package config

import (
	"context"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

func loadDestination(_ context.Context, cmd *cliz.Command) string {
	v, _ := cmd.GetOptionString(_OptionDestination)
	return v
}

func Destination() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig.Destination
}
