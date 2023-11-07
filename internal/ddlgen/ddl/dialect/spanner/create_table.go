package spanner

import (
	"fmt"
	"strconv"

	filepathz "github.com/kunitsucom/util.go/path/filepath"
	slicez "github.com/kunitsucom/util.go/slices"

	ddlast "github.com/kunitsucom/ddlgen/internal/ddlgen/ddl"
)

//nolint:cyclop,funlen,gocognit
func fprintCreateTable(buf *string, indent string, stmt *ddlast.CreateTableStmt) {
	// source
	if stmt.SourceFile != "" {
		fprintComment(buf, "", fmt.Sprintf("source: %s:%d", filepathz.Short(stmt.SourceFile), stmt.SourceLine))
	}

	// comments
	for _, comment := range stmt.Comments {
		fprintComment(buf, "", comment)
	}

	if stmt.CreateTable != "" { //nolint:nestif
		// CREATE TABLE and Left Parenthesis
		*buf += stmt.CreateTable + " (\n"

		hasTableConstraint := len(stmt.Constraints) > 0

		// COLUMNS
		fprintCreateTableColumn(buf, indent, stmt.Columns, hasTableConstraint)

		// CONSTRAINT
		for i, constraint := range stmt.Constraints {
			fprintCreateTableConstraint(buf, indent, constraint)
			if lastConstraint := len(stmt.Constraints) - 1; i == lastConstraint {
				*buf += "\n"
			} else {
				*buf += ",\n"
			}
		}

		// Right Parenthesis
		*buf += ")"

		// OPTIONS
		for i, option := range stmt.Options {
			*buf += "\n"
			fprintCreateTableOption(buf, "", option)
			if lastOption := len(stmt.Options) - 1; i != lastOption {
				*buf += ","
			}
		}

		*buf += ";\n"
	}

	return //nolint:gosimple
}

func fprintCreateTableColumn(buf *string, indent string, columns []*ddlast.CreateTableColumn, tailComma bool) {
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
			fprintComment(buf, indent, comment)
		}

		*buf += indent + fmt.Sprintf(columnNameFormat, Quotation+column.Column+Quotation) + " " + column.TypeConstraint

		if lastColumn := len(columns) - 1; i == lastColumn && !tailComma {
			*buf += "\n"
		} else {
			*buf += ",\n"
		}
	}

	return //nolint:gosimple
}

func fprintCreateTableConstraint(buf *string, indent string, constraint *ddlast.CreateTableConstraint) {
	for _, comment := range constraint.Comments {
		fprintComment(buf, indent, comment)
	}

	*buf += indent + constraint.Constraint

	return //nolint:gosimple
}

func fprintCreateTableOption(buf *string, indent string, option *ddlast.CreateTableOption) {
	for _, comment := range option.Comments {
		fprintComment(buf, indent, comment)
	}

	*buf += indent + option.Option

	return //nolint:gosimple
}
