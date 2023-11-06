package config

import (
	"context"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

func loadTimestamp(_ context.Context, cmd *cliz.Command) string {
	v, _ := cmd.GetOptionString(_OptionTimestamp)
	return v
}

func Timestamp() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig.Timestamp
}
