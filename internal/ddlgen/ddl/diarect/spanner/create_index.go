package spanner

import (
	"fmt"
	"io"

	errorz "github.com/kunitsucom/util.go/errors"
	filepathz "github.com/kunitsucom/util.go/path/filepath"

	ddlast "github.com/kunitsucom/ddlgen/internal/ddlgen/ddl"
)

//nolint:cyclop,funlen
func fprintCreateIndex(w io.Writer, _ string, stmt *ddlast.CreateIndexStmt) error {
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

	// CREATE INDEX
	if _, err := io.WriteString(w, stmt.CreateIndex); err != nil {
		return errorz.Errorf("io.WriteString: %w", err)
	}

	if _, err := io.WriteString(w, ";\n"); err != nil {
		return errorz.Errorf("io.WriteString: %w", err)
	}

	return nil
}
