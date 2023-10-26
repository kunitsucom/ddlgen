package config

import (
	"context"

	"github.com/kunitsucom/util.go/flagenv"
)

const (
	_EnvKeySource  = "SOURCE"
	_DefaultSource = "/dev/stdin"
)

func loadSource(_ context.Context, fes *flagenv.FlagEnvSet) *string {
	v := fes.String("src", _EnvKeySource, _DefaultSource, "source file or directory (default: /dev/stdin)")
	return v
}

func Source() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return *globalConfig.Source
}
