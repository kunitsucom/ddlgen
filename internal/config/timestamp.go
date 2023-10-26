package config

import (
	"context"

	"github.com/kunitsucom/util.go/flagenv"
)

const (
	_EnvKeyTimestamp  = "TIMESTAMP"
	_DefaultTimestamp = ""
)

func loadTimestamp(_ context.Context, fes *flagenv.FlagEnvSet) *string {
	v := fes.String("timestamp", _EnvKeyTimestamp, _DefaultTimestamp, "source file or directory (default: /dev/stdin)")
	return v
}

func Timestamp() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return *globalConfig.Timestamp
}
