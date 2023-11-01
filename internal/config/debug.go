package config

import (
	"context"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

const (
	_EnvKeyDebug  = "DEBUG"
	_DefaultDebug = false
)

func loadDebug(_ context.Context, cmd *cliz.Command) bool {
	v, _ := cmd.GetBoolOption(optionDebug)
	return v
}

func Debug() bool {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig.Debug
}
