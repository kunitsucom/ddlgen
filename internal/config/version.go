package config

import (
	"context"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

func loadVersion(_ context.Context, cmd *cliz.Command) bool {
	v, _ := cmd.GetOptionBool(_OptionVersion)
	return v
}

func Version() bool {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig.Version
}
