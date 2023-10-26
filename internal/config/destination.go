package config

import (
	"context"

	"github.com/kunitsucom/util.go/flagenv"
)

const (
	_EnvKeyDestination  = "DESTINATION"
	_DefaultDestination = "/dev/stdout"
)

func loadDestination(_ context.Context, fes *flagenv.FlagEnvSet) *string {
	v := fes.String("dst", _EnvKeyDestination, _DefaultDestination, "destination file or directory (default: /dev/stdout)")
	return v
}

func Destination() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return *globalConfig.Destination
}
