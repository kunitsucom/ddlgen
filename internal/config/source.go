package config

import (
	"context"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

func loadSource(_ context.Context, cmd *cliz.Command) string {
	v, _ := cmd.GetOptionString(_OptionSource)
	return v
}

func Source() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig.Source
}
