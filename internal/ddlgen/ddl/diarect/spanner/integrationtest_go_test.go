//nolint:testpackage
package spanner

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/kunitsucom/util.go/testing/assert"
	"github.com/kunitsucom/util.go/testing/require"

	"github.com/kunitsucom/ddlgen/internal/config"
	"github.com/kunitsucom/ddlgen/internal/contexts"
	ddlgengo "github.com/kunitsucom/ddlgen/internal/ddlgen/lang/go"
)

func Test_integrationtest_go_spanner(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := contexts.WithArgs(context.Background(), []string{
			"ddlgen",
			"--timestamp=2021-01-01T09:00:00+09:00",
			"--lang=go",
			"--dialect=spanner",
			"--column-key-go=dbtest",
			"--ddl-key-go=spanddl",
			"--src=integrationtest_go_spanner_001.source",
			"--dst=dummy",
		})

		_, err := config.Load(ctx)
		require.NoError(t, err)

		ctx = contexts.WithNowString(ctx, time.RFC3339, config.Timestamp())

		ddl, err := ddlgengo.Parse(ctx, config.Source())
		require.NoError(t, err)

		buf := bytes.NewBuffer(nil)

		require.NoError(t, Fprint(buf, ddl))

		golden, err := os.ReadFile("integrationtest_go_spanner_001.golden")
		require.NoError(t, err)

		if !assert.Equal(t, string(golden), buf.String()) {
			fmt.Println(buf.String()) //nolint:forbidigo
		}
	})
}
