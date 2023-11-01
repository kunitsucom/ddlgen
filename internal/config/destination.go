package config

import (
	"context"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

const (
	_EnvKeyDestination  = "DESTINATION"
	_DefaultDestination = "/dev/stdout"
)

func loadDestination(_ context.Context, cmd *cliz.Command) string {
	v, _ := cmd.GetStringOption(optionDestination)
	return v
}

func Destination() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig.Destination
}
