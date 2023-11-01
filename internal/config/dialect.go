package config

import (
	"context"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

const (
	_EnvKeyDialect  = "DIALECT"
	_DefaultDialect = "postgres"
)

func loadDialect(_ context.Context, cmd *cliz.Command) string {
	v, _ := cmd.GetStringOption(optionDialect)
	return v
}

func Dialect() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig.Dialect
}
