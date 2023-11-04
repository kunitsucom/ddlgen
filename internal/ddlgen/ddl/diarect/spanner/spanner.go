package spanner

import (
	"io"

	errorz "github.com/kunitsucom/util.go/errors"

	ddlast "github.com/kunitsucom/ddlgen/internal/ddlgen/ddl"
	"github.com/kunitsucom/ddlgen/internal/logs"
	"github.com/kunitsucom/ddlgen/pkg/errors"
)

const (
	Dialect       = "spanner"
	CommentPrefix = "--"
	Quotation     = "`"
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
		case *ddlast.CreateIndexStmt:
			if err := fprintCreateIndex(w, ddl.Indent, stmt); err != nil {
				return errorz.Errorf("fprintCreateIndex: %w", err)
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
		if _, err := io.WriteString(w, indent+CommentPrefix+"\n"); err != nil {
			return errorz.Errorf("io.WriteString: %w", err)
		}
		return nil
	}

	if _, err := io.WriteString(w, indent+CommentPrefix+" "+comment+"\n"); err != nil {
		return errorz.Errorf("io.WriteString: %w", err)
	}
	return nil
}
