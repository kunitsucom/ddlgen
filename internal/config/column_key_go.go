package config

import (
	"context"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

const (
	_EnvKeyColumnKeyGo  = "COLUMN_KEY_GO"
	_DefaultColumnKeyGo = "db"
)

func loadColumnKeyGo(_ context.Context, cmd *cliz.Command) string {
	v, _ := cmd.GetStringOption(optionColumnKeyGo)
	return v
}

func ColumnKeyGo() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig.ColumnKeyGo
}
