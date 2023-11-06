package config

import (
	"context"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

func loadDDLKeyGo(_ context.Context, cmd *cliz.Command) string {
	v, _ := cmd.GetOptionString(_OptionDDLKeyGo)
	return v
}

func DDLKeyGo() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig.DDLKeyGo
}
