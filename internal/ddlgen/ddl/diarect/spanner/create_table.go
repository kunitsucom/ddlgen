package spanner

import (
	"fmt"
	"io"
	"strconv"

	errorz "github.com/kunitsucom/util.go/errors"
	filepathz "github.com/kunitsucom/util.go/path/filepath"
	slicez "github.com/kunitsucom/util.go/slices"

	ddlast "github.com/kunitsucom/ddlgen/internal/ddlgen/ddl"
)

//nolint:cyclop,funlen,gocognit
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

	if stmt.CreateTable != "" { //nolint:nestif
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
		for i, option := range stmt.Options {
			if _, err := io.WriteString(w, " "); err != nil {
				return errorz.Errorf("io.WriteString: %w", err)
			}
			if err := fprintCreateTableOption(w, "", option); err != nil {
				return errorz.Errorf("fprintCreateTableOption: %w", err)
			}
			if lastOption := len(stmt.Options) - 1; i != lastOption {
				if _, err := io.WriteString(w, ",\n"); err != nil {
					return errorz.Errorf("io.WriteString: %w", err)
				}
			}
		}

		if _, err := io.WriteString(w, ";\n"); err != nil {
			return errorz.Errorf("io.WriteString: %w", err)
		}
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

		columnLine := indent + fmt.Sprintf(columnNameFormat, Quotation+column.Column+Quotation) + " " + column.TypeConstraint
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
