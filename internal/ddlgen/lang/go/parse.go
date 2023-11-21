package ddlgengo

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"unicode"

	errorz "github.com/kunitsucom/util.go/errors"
	filepathz "github.com/kunitsucom/util.go/path/filepath"
	slicez "github.com/kunitsucom/util.go/slices"

	"github.com/kunitsucom/ddlgen/internal/config"
	ddlast "github.com/kunitsucom/ddlgen/internal/ddlgen/ddl"
	langutil "github.com/kunitsucom/ddlgen/internal/ddlgen/lang/util"
	"github.com/kunitsucom/ddlgen/internal/logs"
	"github.com/kunitsucom/ddlgen/internal/util"
	apperr "github.com/kunitsucom/ddlgen/pkg/errors"
)

const (
	Language                                = "go"
	DDLGEN_ERROR_STRUCT_FIELD_TAG_NOT_FOUND = "DDLGEN_ERROR_STRUCT_FIELD_TAG_NOT_FOUND" //nolint:revive,stylecheck
)

//nolint:cyclop
func Parse(ctx context.Context, src string) (*ddlast.DDL, error) {
	// MEMO: get absolute path for parser.ParseFile()
	sourceAbs := util.Abs(src)

	info, err := os.Stat(sourceAbs)
	if err != nil {
		return nil, errorz.Errorf("os.Stat: %w", err)
	}

	ddl := ddlast.NewDDL(ctx)

	if info.IsDir() {
		if err := filepath.WalkDir(sourceAbs, walkDirFn(ctx, ddl)); err != nil {
			return nil, errorz.Errorf("filepath.WalkDir: %w", err)
		}

		return ddl, nil
	}

	stmts, err := parseFile(ctx, sourceAbs)
	if err != nil {
		return nil, errorz.Errorf("Parse: %w", err)
	}
	ddl.Stmts = append(ddl.Stmts, stmts...)

	return ddl, nil
}

//nolint:gochecknoglobals
var fileSuffix = ".go"

func walkDirFn(ctx context.Context, ddl *ddlast.DDL) func(path string, d os.DirEntry, err error) error {
	return func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err //nolint:wrapcheck
		}

		if d.IsDir() || !strings.HasSuffix(path, fileSuffix) || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		stmts, err := parseFile(ctx, path)
		if err != nil {
			if errors.Is(err, apperr.ErrDDLTagGoAnnotationNotFoundInSource) {
				logs.Debug.Printf("parseFile: %s: %v", path, err)
				return nil
			}
			return errorz.Errorf("parseFile: %w", err)
		}

		ddl.Stmts = append(ddl.Stmts, stmts...)

		return nil
	}
}

//nolint:cyclop,funlen,gocognit
func parseFile(ctx context.Context, filename string) ([]ddlast.Stmt, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, errorz.Errorf("parser.ParseFile: %w", err)
	}

	ddlSrc, err := extractDDLSourceFromDDLTagGo(ctx, fset, f)
	if err != nil {
		return nil, errorz.Errorf("extractDDLSourceFromDDLTagGo: %w", err)
	}

	dumpDDLSource(fset, ddlSrc)

	stmts := make([]ddlast.Stmt, 0)
	for _, r := range ddlSrc {
		createTableStmt := &ddlast.CreateTableStmt{}

		// source
		createTableStmt.SourceFile = r.Position.Filename
		createTableStmt.SourceLine = r.Position.Line

		// CREATE TABLE (or INDEX) / CONSTRAINT / OPTIONS (from comments)
		comments := slicez.Select(r.CommentGroup.List, func(_ int, comment *ast.Comment) string {
			return strings.TrimLeftFunc(strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(comment.Text, "//"), "/*"), "*/"), unicode.IsSpace)
		})
		for _, comment := range comments {
			logs.Debug.Printf("[COMMENT DETECTED]: %s:%d: %s", createTableStmt.SourceFile, createTableStmt.SourceLine, comment)

			// NOTE: CREATE INDEX may be written in CREATE TABLE annotation, so process it here
			if /* CREATE INDEX */ matches := langutil.StmtRegexCreateIndex.Regex.FindStringSubmatch(comment); len(matches) > langutil.StmtRegexCreateIndex.Index {
				commentMatchedCreateIndex := comment
				source := fset.Position(extractContainingCommentFromCommentGroup(r.CommentGroup, commentMatchedCreateIndex).Pos())
				createIndexStmt := &ddlast.CreateIndexStmt{
					Comments:   []string{commentMatchedCreateIndex},
					SourceFile: source.Filename,
					SourceLine: source.Line,
				}
				createIndexStmt.SetCreateIndex(matches[langutil.StmtRegexCreateIndex.Index])
				stmts = append(stmts, createIndexStmt)
				continue
			}

			if /* CREATE TABLE */ matches := langutil.StmtRegexCreateTable.Regex.FindStringSubmatch(comment); len(matches) > langutil.StmtRegexCreateTable.Index {
				createTableStmt.SetCreateTable(matches[langutil.StmtRegexCreateTable.Index])
			} else if /* CONSTRAINT */ matches := langutil.StmtRegexCreateTableConstraint.Regex.FindStringSubmatch(comment); len(matches) > langutil.StmtRegexCreateTableConstraint.Index {
				createTableStmt.Constraints = append(createTableStmt.Constraints, &ddlast.CreateTableConstraint{
					Constraint: matches[langutil.StmtRegexCreateTableConstraint.Index],
				})
			} else if /* OPTIONS */ matches := langutil.StmtRegexCreateTableOptions.Regex.FindStringSubmatch(comment); len(matches) > langutil.StmtRegexCreateTableOptions.Index {
				createTableStmt.Options = append(createTableStmt.Options, &ddlast.CreateTableOption{
					Option: matches[langutil.StmtRegexCreateTableOptions.Index],
				})
			}
			// comment
			createTableStmt.Comments = append(createTableStmt.Comments, comment)
		}

		// CREATE TABLE (default: struct name)
		if r.TypeSpec != nil && createTableStmt.CreateTable == "" {
			name := r.TypeSpec.Name.String()
			source := fset.Position(r.CommentGroup.Pos())
			createTableStmt.Comments = append(createTableStmt.Comments, fmt.Sprintf("WARN: the comment (%s:%d) does not have a key for table (%s: table: CREATE TABLE <table>), so the struct name \"%s\" is used as the table name.", filepathz.Short(source.Filename), source.Line, config.DDLTagGo(), name))
			createTableStmt.SetCreateTable(name)
		}

		// columns
		if r.StructType != nil {
			for _, field := range r.StructType.Fields.List {
				column := &ddlast.CreateTableColumn{}

				tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))

				// column name
				switch columnName := tag.Get(config.ColumnTagGo()); columnName {
				case "-":
					createTableStmt.Comments = append(createTableStmt.Comments, fmt.Sprintf("NOTE: the \"%s\" struct's \"%s\" field has a tag for column name (`%s:\"-\"`), so the field is ignored.", r.TypeSpec.Name, field.Names[0], config.ColumnTagGo()))
					continue
				case "":
					name := field.Names[0].Name
					column.Comments = append(column.Comments, fmt.Sprintf("WARN: the \"%s\" struct's \"%s\" field does not have a tag for column name (`%s:\"<ColumnName>\"`), so the field name \"%s\" is used as the column name.", r.TypeSpec.Name, field.Names[0], config.ColumnTagGo(), name))
					column.ColumnName = name
				default:
					column.ColumnName = columnName
				}

				// column type and constraint
				switch columnTypeConstraint := tag.Get(config.DDLTagGo()); columnTypeConstraint {
				case "", "-":
					column.Comments = append(column.Comments, fmt.Sprintf("ERROR: the \"%s\" struct's \"%s\" field does not have a tag for column type and constraint (`%s:\"<TYPE> [CONSTRAINT]\"`)", r.TypeSpec.Name, field.Names[0], config.DDLTagGo()))
					column.TypeConstraint = DDLGEN_ERROR_STRUCT_FIELD_TAG_NOT_FOUND
				default:
					column.TypeConstraint = columnTypeConstraint
				}

				// primary key
				switch primaryKey := tag.Get(config.PKTagGo()); primaryKey {
				case "true", "1":
					createTableStmt.PrimaryKey = append(createTableStmt.PrimaryKey, column.ColumnName)
				case "", "-":
					// do nothing
				default:
					column.Comments = append(column.Comments, fmt.Sprintf("WARN: the field \"%s\" does not have valid primary key tag (`%s:\"true\"`), so the column is not used as primary key.", field.Names[0], config.PKTagGo()))
				}

				// comments
				comments := strings.Split(strings.Trim(field.Doc.Text(), "\n"), "\n")
				column.Comments = append(column.Comments, langutil.TrimCommentElementTailEmpty(langutil.TrimCommentElementHasPrefix(comments, config.DDLTagGo()))...)

				createTableStmt.Columns = append(createTableStmt.Columns, column)
			}
		}

		if createTableStmt.CreateTable == "" {
			// CREATE TABLE (ERROR)
			source := fset.Position(r.CommentGroup.Pos())
			createTableStmt.Comments = append(createTableStmt.Comments, fmt.Sprintf("WARN: the comment (%s:%d) does not have a key for table (%s: table: CREATE TABLE <table>), or the comment is not associated with struct.", filepathz.Short(source.Filename), source.Line, config.DDLTagGo()))
		} else if len(createTableStmt.Columns) == 0 {
			// columns (ERROR)
			source := fset.Position(r.CommentGroup.Pos())
			createTableStmt.Comments = append(createTableStmt.Comments, fmt.Sprintf("ERROR: the comment (%s:%d) does not have struct fields for column type and constraint (`%s:\"<TYPE> [CONSTRAINT]\"`), or the comment is not associated with struct.", filepathz.Short(source.Filename), source.Line, config.DDLTagGo()))
		}

		stmts = append(stmts, createTableStmt)
	}

	sort.Slice(stmts, func(i, j int) bool {
		return fmt.Sprintf("%s:%09d", stmts[i].GetSourceFile(), stmts[i].GetSourceLine()) < fmt.Sprintf("%s:%09d", stmts[j].GetSourceFile(), stmts[j].GetSourceLine())
	})

	for i := range stmts {
		logs.Trace.Print(fmt.Sprintf("%s:%09d", stmts[i].GetSourceFile(), stmts[i].GetSourceLine()))
	}

	return stmts, nil
}

func extractContainingCommentFromCommentGroup(commentGroup *ast.CommentGroup, sub string) *ast.Comment {
	for _, commentLine := range commentGroup.List {
		if strings.Contains(commentLine.Text, sub) {
			return commentLine
		}
	}
	return nil
}
