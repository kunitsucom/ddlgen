package config

import (
	"context"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

func loadColumnKeyGo(_ context.Context, cmd *cliz.Command) string {
	v, _ := cmd.GetOptionString(_OptionColumnKeyGo)
	return v
}

func ColumnKeyGo() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig.ColumnKeyGo
}
