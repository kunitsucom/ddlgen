//nolint:testpackage
package ddlgengo

import (
	"bytes"
	"context"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kunitsucom/ddlgen/internal/config"
	"github.com/kunitsucom/ddlgen/internal/contexts"
	"github.com/kunitsucom/ddlgen/internal/logs"
)

func TestExtractDDLSource(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := contexts.WithArgs(context.Background(), []string{
			"ddlgen",
			"--lang=go",
			"--dialect=spanner",
			"--column-key-go=db",
			"--ddl-key-go=spanddl",
			"--src=dummy",
			"--dst=dummy",
		})

		if _, err := config.Load(ctx); err != nil {
			t.Fatalf("❌: config.Load: %+v", err)
		}

		fset := token.NewFileSet()
		src := bytes.NewBufferString(`
package main

type (
	// User is a user.
	//
	// spanddl:table:` + "`users`" + `
	// spanddl:options:PRIMARY KEY (` + "`id`" + `)
	User struct {
		// ID is a user ID.
		ID string   ` + "`db:\"Id\"   spanddl:\"STRING(36)  NOT NULL\"`" + `
		// Name is a user name.
		Name string ` + "`db:\"Name\" spanddl:\"STRING(255) NOT NULL\"`" + `
	}

	// Book is a book.
	//
	// spanddl:table:` + "`books`" + `
	Book struct {
		// ID is a book ID.
		ID string   ` + "`db:\"Id\"   spanddl:\"STRING(36)  NOT NULL\"`" + `
	}
)
`)

		f, err := parser.ParseFile(fset, "", bytes.NewBuffer(src.Bytes()), parser.ParseComments)
		if err != nil {
			t.Fatalf("❌: parser.ParseFile: %+v", err)
		}

		ddlSrc, err := extractDDLSource(ctx, fset, f)
		if err != nil {
			t.Fatalf("❌: extractDDLSource: %+v", err)
		}

		logs.Trace = logs.NewTrace()
		dumpDDLSource(fset, ddlSrc)
		for _, s := range ddlSrc {
			t.Logf("✅: ddlSrc: %#v", s)
		}
	})
}
