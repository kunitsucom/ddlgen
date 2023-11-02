package config

import (
	"context"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

func loadDialect(_ context.Context, cmd *cliz.Command) string {
	v, _ := cmd.GetStringOption(_OptionDialect)
	return v
}

func Dialect() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig.Dialect
}
