package config

import (
	"context"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

func loadDDLTagGo(_ context.Context, cmd *cliz.Command) string {
	v, _ := cmd.GetOptionString(_OptionDDLTagGo)
	return v
}

func DDLTagGo() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig.DDLTagGo
}
