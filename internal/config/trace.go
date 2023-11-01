package config

import (
	"context"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

const (
	_EnvKeyTrace  = "TRACE"
	_DefaultTrace = false
)

func loadTrace(_ context.Context, cmd *cliz.Command) bool {
	v, _ := cmd.GetBoolOption(optionTrace)
	return v
}

func Trace() bool {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig.Trace
}
