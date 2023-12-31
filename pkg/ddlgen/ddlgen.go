package ddlgen

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	errorz "github.com/kunitsucom/util.go/errors"
	cliz "github.com/kunitsucom/util.go/exp/cli"
	"github.com/kunitsucom/util.go/version"

	"github.com/kunitsucom/ddlgen/internal/config"
	ddlast "github.com/kunitsucom/ddlgen/internal/ddlgen/ddl"
	"github.com/kunitsucom/ddlgen/internal/ddlgen/ddl/dialect/postgres"
	"github.com/kunitsucom/ddlgen/internal/ddlgen/ddl/dialect/spanner"
	ddlgengo "github.com/kunitsucom/ddlgen/internal/ddlgen/lang/go"
	"github.com/kunitsucom/ddlgen/internal/logs"
	apperr "github.com/kunitsucom/ddlgen/pkg/errors"
)

//nolint:cyclop
func DDLGen(ctx context.Context) error {
	if _, err := config.Load(ctx); err != nil {
		if errors.Is(err, cliz.ErrHelp) {
			return nil
		}
		return fmt.Errorf("config.Load: %w", err)
	}

	if config.Version() {
		fmt.Printf("version: %s\n", version.Version())           //nolint:forbidigo
		fmt.Printf("revision: %s\n", version.Revision())         //nolint:forbidigo
		fmt.Printf("build branch: %s\n", version.Branch())       //nolint:forbidigo
		fmt.Printf("build timestamp: %s\n", version.Timestamp()) //nolint:forbidigo
		return nil
	}

	src := config.Source()
	logs.Info.Printf("source: %s", src)

	ddl, err := parse(ctx, src)
	if err != nil {
		return errorz.Errorf("parse: %w", err)
	}

	if info, err := os.Stat(config.Destination()); err == nil && info.IsDir() {
		for _, stmt := range ddl.Stmts {
			dst := filepath.Join(config.Destination(), filepath.Base(stmt.GetSourceFile())) + ".gen.sql"
			logs.Info.Printf("destination: %s", dst)

			f, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
			if err != nil {
				return errorz.Errorf("os.OpenFile: %w", err)
			}

			if err := fprint(f, &ddlast.DDL{
				Header: ddl.Header,
				Indent: ddl.Indent,
				Stmts:  []ddlast.Stmt{stmt},
			}); err != nil {
				return errorz.Errorf("fprint: %w", err)
			}
		}
		return nil
	}

	dst := config.Destination()
	logs.Info.Printf("destination: %s", dst)

	f, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return errorz.Errorf("os.OpenFile: %w", err)
	}

	if err := fprint(f, ddl); err != nil {
		return errorz.Errorf("fprint: %w", err)
	}
	return nil
}

func parse(ctx context.Context, src string) (*ddlast.DDL, error) {
	switch language := config.Language(); language {
	case ddlgengo.Language:
		ddl, err := ddlgengo.Parse(ctx, src)
		if err != nil {
			return nil, errorz.Errorf("ddlgengo.Parse: %w", err)
		}
		return ddl, nil
	default:
		return nil, errorz.Errorf("language=%s: %w", language, apperr.ErrNotSupported)
	}
}

func fprint(w io.Writer, ddl *ddlast.DDL) error {
	switch dialect := config.Dialect(); dialect {
	case spanner.Dialect:
		if err := spanner.Fprint(w, ddl); err != nil {
			return errorz.Errorf("spanner.Fprint: %w", err)
		}
		return nil
	case postgres.Dialect:
		if err := postgres.Fprint(w, ddl); err != nil {
			return errorz.Errorf("postgres.Fprint: %w", err)
		}
		return nil
	default:
		return errorz.Errorf("dialect=%s: %w", dialect, apperr.ErrNotSupported)
	}
}
