package config

import (
	"context"

	"github.com/kunitsucom/util.go/flagenv"
)

const (
	_EnvKeyDebug  = "DEBUG"
	_DefaultDebug = false
)

func loadDebug(_ context.Context, fes *flagenv.FlagEnvSet) *bool {
	v := fes.Bool("debug", _EnvKeyDebug, _DefaultDebug, "debug mode (default: false)")
	return v
}

func Debug() bool {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return *globalConfig.Debug
}
