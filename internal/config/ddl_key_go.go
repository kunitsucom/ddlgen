package config

import (
	"context"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

const (
	_EnvKeyDDLKeyGo  = "DDL_KEY_GO"
	_DefaultDDLKeyGo = "ddlgen"
)

func loadDDLKeyGo(_ context.Context, cmd *cliz.Command) string {
	v, _ := cmd.GetStringOption(optionDDLKeyGo)
	return v
}

func DDLKeyGo() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig.DDLKeyGo
}
