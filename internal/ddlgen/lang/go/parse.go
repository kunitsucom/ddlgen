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

	errorz "github.com/kunitsucom/util.go/errors"

	"github.com/kunitsucom/ddlgen/internal/config"
	ddlast "github.com/kunitsucom/ddlgen/internal/ddlgen/ddl"
	"github.com/kunitsucom/ddlgen/internal/ddlgen/lang/util"
	"github.com/kunitsucom/ddlgen/internal/logs"
	apperr "github.com/kunitsucom/ddlgen/pkg/errors"
)

//nolint:cyclop
func Parse(ctx context.Context, src string) (*ddlast.DDL, error) {
	// MEMO: get absolute path for parser.ParseFile()
	sourceAbs, err := filepath.Abs(src)
	if err != nil {
		return nil, errorz.Errorf("filepath.Abs: %w", err)
	}

	info, err := os.Stat(sourceAbs)
	if err != nil {
		return nil, errorz.Errorf("os.Stat: %w", err)
	}

	ddl := ddlast.NewDDL(ctx)

	if info.IsDir() {
		if err := filepath.WalkDir(sourceAbs, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err //nolint:wrapcheck
			}

			if d.IsDir() {
				return nil
			}

			if !strings.HasSuffix(path, ".go") {
				return nil
			}

			stmts, err := parseFile(ctx, path)
			if err != nil {
				if errors.Is(err, apperr.ErrDDLKeyGoNotFoundInSource) {
					logs.Debug.Printf("parseFile: %s: %v", path, err)
					return nil
				}
				return errorz.Errorf("parseFile: %w", err)
			}

			ddl.Stmts = append(ddl.Stmts, stmts...)

			return nil
		}); err != nil {
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

//nolint:cyclop,funlen,gocognit
func parseFile(ctx context.Context, filename string) ([]ddlast.Stmt, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, errorz.Errorf("parser.ParseFile: %w", err)
	}

	ddlSrc, err := extractDDLSourceFromDDLKeyGo(ctx, fset, f)
	if err != nil {
		return nil, errorz.Errorf("extractDDLSource: %w", err)
	}

	dumpDDLSource(fset, ddlSrc)

	stmts := make([]ddlast.Stmt, 0)
	for _, r := range ddlSrc {
		createTableStmt := &ddlast.CreateTableStmt{}

		// source
		source := fset.Position(r.CommentGroup.Pos())
		createTableStmt.SourceFile = source.Filename
		createTableStmt.SourceLine = source.Line

		// CREATE TABLE (or INDEX) / CONSTRAINT / OPTIONS (from comments)
		comments := strings.Split(strings.Trim(r.CommentGroup.Text(), "\n"), "\n")
		for _, comment := range comments {
			logs.Debug.Printf("[COMMENT DETECTED]: %s:%d: %s", createTableStmt.SourceFile, createTableStmt.SourceLine, comment)

			// NOTE: CREATE INDEX may be written in CREATE TABLE annotation, so process it here
			if /* CREATE INDEX */ matches := util.StmtRegexCreateIndex.Regex.FindStringSubmatch(comment); len(matches) > util.StmtRegexCreateIndex.Index {
				commentMatchedCreateIndex := comment
				source := fset.Position(extractContainingCommentFromCommentGroup(r.CommentGroup, commentMatchedCreateIndex).Pos())
				createIndexStmt := &ddlast.CreateIndexStmt{
					Comments:   []string{commentMatchedCreateIndex},
					SourceFile: source.Filename,
					SourceLine: source.Line,
				}
				createIndexStmt.SetCreateIndex(matches[util.StmtRegexCreateIndex.Index])
				stmts = append(stmts, createIndexStmt)
				continue
			}

			if /* CREATE TABLE */ matches := util.StmtRegexCreateTable.Regex.FindStringSubmatch(comment); len(matches) > util.StmtRegexCreateTable.Index {
				createTableStmt.SetCreateTable(matches[util.StmtRegexCreateTable.Index])
			} else if /* CONSTRAINT */ matches := util.StmtRegexCreateTableConstraint.Regex.FindStringSubmatch(comment); len(matches) > util.StmtRegexCreateTableConstraint.Index {
				createTableStmt.Constraints = append(createTableStmt.Constraints, &ddlast.CreateTableConstraint{
					Constraint: matches[util.StmtRegexCreateTableConstraint.Index],
				})
			} else if /* OPTIONS */ matches := util.StmtRegexCreateTableOptions.Regex.FindStringSubmatch(comment); len(matches) > util.StmtRegexCreateTableOptions.Index {
				createTableStmt.Options = append(createTableStmt.Options, &ddlast.CreateTableOption{
					Option: matches[util.StmtRegexCreateTableOptions.Index],
				})
			}
			// comment
			createTableStmt.Comments = append(createTableStmt.Comments, comment)
		}

		// CREATE TABLE (default: struct name)
		if createTableStmt.CreateTable == "" && r.TypeSpec != nil {
			name := r.TypeSpec.Name.String()
			createTableStmt.Comments = append(createTableStmt.Comments, fmt.Sprintf("NOTE: the comment does not have a key for table (%s: table: CREATE TABLE <table>), so the struct name \"%s\" is used as the table name.", config.DDLKeyGo(), name))
			createTableStmt.SetCreateTable(name)
		}

		// CREATE TABLE (error)
		if createTableStmt.CreateTable == "" {
			createTableStmt.Comments = append(createTableStmt.Comments, fmt.Sprintf("ERROR: the comment does not have a key for table (%s: table: CREATE TABLE <table>), or the comment is not associated with struct.", config.DDLKeyGo()))
		}

		// columns
		if r.StructType != nil {
			for _, field := range r.StructType.Fields.List {
				column := &ddlast.CreateTableColumn{}

				tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))

				// column name
				switch columnName := tag.Get(config.ColumnKeyGo()); columnName {
				case "-":
					logs.Info.Printf("[%s]: Ignore columns with column name \"-\": %s", createTableStmt.CreateTable, field.Names[0].Name)
					continue
				case "":
					name := field.Names[0].Name
					column.Comments = append(column.Comments, fmt.Sprintf("NOTE: the struct does not have a tag for column name (`%s:\"<ColumnName>\"`), so the field name \"%s\" is used as the column name.", config.ColumnKeyGo(), name))
					column.Column = name
				default:
					column.Column = columnName
				}

				// column type and constraint
				switch columnTypeConstraint := tag.Get(config.DDLKeyGo()); columnTypeConstraint {
				case "":
					column.Comments = append(column.Comments, fmt.Sprintf("ERROR: the struct does not have a tag for column type and constraint (`%s:\"<TYPE> [CONSTRAINT]\"`)", config.DDLKeyGo()))
					column.TypeConstraint = "ERROR"
				default:
					column.TypeConstraint = columnTypeConstraint
				}

				// comments
				comments := strings.Split(strings.Trim(field.Doc.Text(), "\n"), "\n")
				column.Comments = append(column.Comments, util.TrimTailEmptyCommentElement(util.TrimDDLGenCommentElement(comments))...)

				createTableStmt.Columns = append(createTableStmt.Columns, column)
			}
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
