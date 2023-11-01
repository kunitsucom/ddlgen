package config

import (
	"context"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

const (
	_EnvKeyTimestamp  = "TIMESTAMP"
	_DefaultTimestamp = ""
)

func loadTimestamp(_ context.Context, cmd *cliz.Command) string {
	v, _ := cmd.GetStringOption(optionTimestamp)
	return v
}

func Timestamp() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig.Timestamp
}
