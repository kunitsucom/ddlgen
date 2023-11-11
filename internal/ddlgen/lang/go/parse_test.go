//nolint:testpackage
package ddlgengo

import (
	"context"
	"go/ast"
	"os"
	"path/filepath"
	"testing"

	"github.com/kunitsucom/util.go/testing/assert"
	"github.com/kunitsucom/util.go/testing/require"

	"github.com/kunitsucom/ddlgen/internal/config"
	"github.com/kunitsucom/ddlgen/internal/contexts"
	ddlast "github.com/kunitsucom/ddlgen/internal/ddlgen/ddl"
	apperr "github.com/kunitsucom/ddlgen/pkg/errors"
)

//nolint:paralleltest
func TestParse(t *testing.T) {
	t.Run("success,common.source", func(t *testing.T) {
		ctx := contexts.WithArgs(context.Background(), []string{
			"ddlgen",
			"--lang=go",
			"--dialect=spanner",
			"--column-tag-go=dbtest",
			"--ddl-tag-go=spanddl",
			"--src=tests/common.source",
			"--dst=dummy",
		})

		_, err := config.Load(ctx)
		require.NoError(t, err)

		ddl, err := Parse(ctx, config.Source())
		require.NoError(t, err)
		if !assert.Equal(t, 9, len(ddl.Stmts)) {
			for _, stmt := range ddl.Stmts {
				t.Logf("🚧: ddl.Stmts: %#v", stmt)
			}
		}
	})

	t.Run("success,info.IsDir", func(t *testing.T) {
		ctx := contexts.WithArgs(context.Background(), []string{
			"ddlgen",
			"--lang=go",
			"--dialect=spanner",
			"--column-tag-go=dbtest",
			"--ddl-tag-go=spanddl",
			"--src=tests",
			"--dst=dummy",
		})

		_, err := config.Load(ctx)
		require.NoError(t, err)

		backup := fileSuffix
		t.Cleanup(func() { fileSuffix = backup })
		fileSuffix = ".source"

		{
			ddl, err := Parse(ctx, config.Source())
			require.NoError(t, err)
			if !assert.Equal(t, 10, len(ddl.Stmts)) {
				for _, stmt := range ddl.Stmts {
					t.Logf("🚧: ddl.Stmts: %#v", stmt)
				}
			}
		}
	})

	t.Run("failure,info.IsDir", func(t *testing.T) {
		tempDir := t.TempDir()
		{
			f, err := os.Create(filepath.Join(tempDir, "error.go"))
			require.NoError(t, err)
			_ = f.Close()
		}

		ctx := contexts.WithArgs(context.Background(), []string{
			"ddlgen",
			"--lang=go",
			"--dialect=spanner",
			"--column-tag-go=dbtest",
			"--ddl-tag-go=spanddl",
			"--src=" + tempDir,
			"--dst=dummy",
		})

		_, err := config.Load(ctx)
		require.NoError(t, err)

		{
			_, err := Parse(ctx, config.Source())
			require.ErrorContains(t, err, "expected 'package', found 'EOF'")
		}
	})

	t.Run("failure,os.ErrNotExist", func(t *testing.T) {
		ctx := contexts.WithArgs(context.Background(), []string{
			"ddlgen",
			"--lang=go",
			"--dialect=spanner",
			"--column-tag-go=dbtest",
			"--ddl-tag-go=spanddl",
			"--src=tests/no-such-file.source",
			"--dst=dummy",
		})

		_, err := config.Load(ctx)
		require.NoError(t, err)

		{
			t.Setenv("PWD", "\\")
			_, err := Parse(ctx, config.Source())
			require.Error(t, err)
			assert.ErrorIs(t, err, os.ErrNotExist)
		}
	})

	t.Run("failure,parser.ParseFile", func(t *testing.T) {
		ctx := contexts.WithArgs(context.Background(), []string{
			"ddlgen",
			"--lang=go",
			"--dialect=spanner",
			"--column-tag-go=dbtest",
			"--ddl-tag-go=spanddl",
			"--src=tests/no.errsource",
			"--dst=dummy",
		})

		_, err := config.Load(ctx)
		require.NoError(t, err)

		{
			_, err := Parse(ctx, config.Source())
			require.Error(t, err)
			assert.ErrorContains(t, err, "expected 'package', found 'EOF'")
		}
	})

	t.Run("failure,extractDDLSource", func(t *testing.T) {
		ctx := contexts.WithArgs(context.Background(), []string{
			"ddlgen",
			"--lang=go",
			"--dialect=spanner",
			"--column-tag-go=dbtest",
			"--ddl-tag-go=spanddl",
			"--src=tests/no-ddl-tag-go.source",
			"--dst=dummy",
		})

		_, err := config.Load(ctx)
		require.NoError(t, err)

		{
			_, err := Parse(ctx, config.Source())
			require.Error(t, err)
			assert.ErrorIs(t, err, apperr.ErrDDLTagGoAnnotationNotFoundInSource)
		}
	})
}

func Test_walkDirFn(t *testing.T) {
	t.Parallel()

	t.Run("failure,err", func(t *testing.T) {
		t.Parallel()

		ctx := contexts.WithArgs(context.Background(), []string{
			"ddlgen",
			"--lang=go",
			"--dialect=spanner",
			"--column-tag-go=dbtest",
			"--ddl-tag-go=spanddl",
			"--src=tests",
			"--dst=dummy",
		})

		_, err := config.Load(ctx)
		require.NoError(t, err)

		ddl := ddlast.NewDDL(ctx)
		fn := walkDirFn(ctx, ddl)
		{
			err := fn("", nil, os.ErrPermission)
			require.Error(t, err)
		}
	})
}

func Test_extractContainingCommentFromCommentGroup(t *testing.T) {
	t.Parallel()

	t.Run("failure,no-such-string", func(t *testing.T) {
		t.Parallel()

		actual := extractContainingCommentFromCommentGroup(&ast.CommentGroup{
			List: []*ast.Comment{
				{
					Text: "// spanddl: index: CREATE INDEX `idx_users_name` ON `users` (`name`)",
				},
			},
		}, "no-such-string")
		assert.Nil(t, actual)
	})
}
