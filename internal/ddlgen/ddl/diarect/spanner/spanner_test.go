//nolint:testpackage
package spanner

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/kunitsucom/ddlgen/internal/contexts"
	ddlast "github.com/kunitsucom/ddlgen/internal/ddlgen/ddl"
)

func TestFprint(t *testing.T) {
	t.Parallel()

	t.Run("success,None", func(t *testing.T) {
		t.Parallel()

		ddl := ddlast.NewDDL(contexts.WithNow(context.Background(), time.Date(2023, 10, 27, 10, 27, 30, 0, time.UTC)))
		ddl.Stmts = []ddlast.Stmt{
			&ddlast.CreateTableStmt{
				Comments:    []string{"Spans is Spanner test table."},
				CreateTable: "CREATE TABLE Spans",
				Columns: []*ddlast.CreateTableColumn{
					{
						Column:         "Id",
						TypeConstraint: "STRING(64) NOT NULL",
						Comments:       []string{"Id is Spans's Id."},
					},
					{
						Column:         "Name",
						TypeConstraint: "STRING(100) NOT NULL",
						Comments:       []string{"Name is Spans's Name."},
					},
					{
						Column:         "Description",
						TypeConstraint: "STRING(1024) NOT NULL",
						Comments:       []string{"Description is Spans's Description."},
					},
				},
				Constraints: []*ddlast.CreateTableConstraint{
					{
						Comments:   []string{"primary key is Id"},
						Constraint: "PRIMARY KEY (Id)",
					},
				},
			},
		}

		buf := bytes.NewBuffer(nil)
		if err := Fprint(buf, ddl); err != nil {
			t.Fatalf("failed to Fprint: %+v", err)
		}

		t.Logf("📝: buf:\n%s", buf.String())
	})
}
