package config

import (
	"context"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

func loadTrace(_ context.Context, cmd *cliz.Command) bool {
	v, _ := cmd.GetOptionBool(_OptionTrace)
	return v
}

func Trace() bool {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig.Trace
}
