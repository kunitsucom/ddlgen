package config

import (
	"context"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

const (
	_EnvKeySource  = "SOURCE"
	_DefaultSource = "/dev/stdin"
)

func loadSource(_ context.Context, cmd *cliz.Command) string {
	v, _ := cmd.GetStringOption(optionSource)
	return v
}

func Source() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig.Source
}
