package config

import (
	"context"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

func loadDebug(_ context.Context, cmd *cliz.Command) bool {
	v, _ := cmd.GetBoolOption(_OptionDebug)
	return v
}

func Debug() bool {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig.Debug
}
