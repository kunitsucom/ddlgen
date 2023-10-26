package spanner

import (
	"fmt"
	"io"
	"strconv"

	ddlast "github.com/kunitsucom/ddlgen/internal/ddlgen/ast"
	"github.com/kunitsucom/ddlgen/internal/logs"
	"github.com/kunitsucom/ddlgen/pkg/errors"
	errorz "github.com/kunitsucom/util.go/errors"
	filepathz "github.com/kunitsucom/util.go/path/filepath"
	slicez "github.com/kunitsucom/util.go/slices"
)

const (
	Dialect       = "spanner"
	commentPrefix = "--"
	quotation     = "`"
)

func Fprint(w io.Writer, ddl *ddlast.DDL) error {
	for _, header := range ddl.Header {
		if err := fprintComment(w, "", header); err != nil {
			return errorz.Errorf("fprintComment: %w", err)
		}
	}

	for _, statement := range ddl.Stmts {
		if _, err := io.WriteString(w, "\n"); err != nil {
			return errorz.Errorf("io.WriteString: %w", err)
		}
		switch stmt := statement.(type) {
		case *ddlast.CreateTableStmt:
			if err := fprintCreateTable(w, ddl.Indent, stmt); err != nil {
				return errorz.Errorf("fprintCreateTable: %w", err)
			}
		default:
			logs.Warn.Printf("unknown statement type: %T: %v", stmt, errors.ErrNotSupported)
			continue
		}
	}
	return nil
}

func fprintComment(w io.Writer, indent string, comment string) error {
	if comment == "" {
		if _, err := io.WriteString(w, indent+commentPrefix+"\n"); err != nil {
			return errorz.Errorf("io.WriteString: %w", err)
		}
		return nil
	}

	if _, err := io.WriteString(w, indent+commentPrefix+" "+comment+"\n"); err != nil {
		return errorz.Errorf("io.WriteString: %w", err)
	}
	return nil
}

//nolint:cyclop,funlen
func fprintCreateTable(w io.Writer, indent string, stmt *ddlast.CreateTableStmt) error {
	// source
	if stmt.SourceFile != "" {
		if err := fprintComment(w, "", fmt.Sprintf("source: %s:%d", filepathz.Short(stmt.SourceFile), stmt.SourceLine)); err != nil {
			return errorz.Errorf("fprintComment: %w", err)
		}
	}

	// comments
	for _, comment := range stmt.Comments {
		if err := fprintComment(w, "", comment); err != nil {
			return errorz.Errorf("fprintComment: %w", err)
		}
	}

	// CREATE TABLE and Left Parenthesis
	if _, err := io.WriteString(w, stmt.CreateTable+" (\n"); err != nil {
		return errorz.Errorf("io.WriteString: %w", err)
	}

	hasTableConstraint := len(stmt.Constraints) > 0

	// COLUMNS
	if err := fprintCreateTableColumn(w, indent, stmt.Columns, hasTableConstraint); err != nil {
		return errorz.Errorf("fprintCreateTableColumn: %w", err)
	}

	// CONSTRAINT
	for i, constraint := range stmt.Constraints {
		if err := fprintCreateTableConstraint(w, indent, constraint); err != nil {
			return errorz.Errorf("fprintCreateTableConstraint: %w", err)
		}
		if lastConstraint := len(stmt.Constraints) - 1; i == lastConstraint {
			if _, err := io.WriteString(w, "\n"); err != nil {
				return errorz.Errorf("io.WriteString: %w", err)
			}
		} else {
			if _, err := io.WriteString(w, ",\n"); err != nil {
				return errorz.Errorf("io.WriteString: %w", err)
			}
		}
	}

	// Right Parenthesis
	if _, err := io.WriteString(w, ")"); err != nil {
		return errorz.Errorf("io.WriteString: %w", err)
	}

	// OPTIONS
	for _, option := range stmt.Options {
		if _, err := io.WriteString(w, " "); err != nil {
			return errorz.Errorf("io.WriteString: %w", err)
		}
		if err := fprintCreateTableOption(w, "", option); err != nil {
			return errorz.Errorf("fprintCreateTableOption: %w", err)
		}
	}

	if _, err := io.WriteString(w, ";\n"); err != nil {
		return errorz.Errorf("io.WriteString: %w", err)
	}

	return nil
}

func fprintCreateTableColumn(w io.Writer, indent string, columns []*ddlast.CreateTableColumn, tailComma bool) error {
	columnNameMaxLength := 0
	slicez.Each(columns, func(index int, elem *ddlast.CreateTableColumn) {
		if columnLength := len(elem.Column); columnLength > columnNameMaxLength {
			columnNameMaxLength = columnLength
		}
	})
	const quotationCharsLength = 2
	columnNameFormat := "%-" + strconv.Itoa(quotationCharsLength+columnNameMaxLength) + "s"

	for i, column := range columns {
		for _, comment := range column.Comments {
			if err := fprintComment(w, indent, comment); err != nil {
				return errorz.Errorf("fprintComment: %w", err)
			}
		}

		columnLine := indent + fmt.Sprintf(columnNameFormat, quotation+column.Column+quotation) + " " + column.TypeConstraint
		if _, err := io.WriteString(w, columnLine); err != nil {
			return errorz.Errorf("io.WriteString: %w", err)
		}

		if lastColumn := len(columns) - 1; i == lastColumn && !tailComma {
			if _, err := io.WriteString(w, "\n"); err != nil {
				return errorz.Errorf("io.WriteString: %w", err)
			}
		} else {
			if _, err := io.WriteString(w, ",\n"); err != nil {
				return errorz.Errorf("io.WriteString: %w", err)
			}
		}
	}
	return nil
}

func fprintCreateTableConstraint(w io.Writer, indent string, constraint *ddlast.CreateTableConstraint) error {
	for _, comment := range constraint.Comments {
		if err := fprintComment(w, indent, comment); err != nil {
			return errorz.Errorf("fprintComment: %w", err)
		}
	}
	if _, err := io.WriteString(w, indent+constraint.Constraint); err != nil {
		return errorz.Errorf("io.WriteString: %w", err)
	}
	return nil
}

func fprintCreateTableOption(w io.Writer, indent string, option *ddlast.CreateTableOption) error {
	for _, comment := range option.Comments {
		if err := fprintComment(w, indent, comment); err != nil {
			return errorz.Errorf("fprintComment: %w", err)
		}
	}
	if _, err := io.WriteString(w, indent+option.Option); err != nil {
		return errorz.Errorf("io.WriteString: %w", err)
	}
	return nil
}
