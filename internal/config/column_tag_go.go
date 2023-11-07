package config

import (
	"context"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

func loadColumnTagGo(_ context.Context, cmd *cliz.Command) string {
	v, _ := cmd.GetOptionString(_OptionColumnTagGo)
	return v
}

func ColumnTagGo() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig.ColumnTagGo
}
