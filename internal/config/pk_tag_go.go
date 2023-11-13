package config

import (
	"context"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

func loadPKTagGo(_ context.Context, cmd *cliz.Command) string {
	v, _ := cmd.GetOptionString(_OptionPKTagGo)
	return v
}

func PKTagGo() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig.PKTagGo
}
