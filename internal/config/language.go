package config

import (
	"context"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

func loadLanguage(_ context.Context, cmd *cliz.Command) string {
	v, _ := cmd.GetStringOption(_OptionLanguage)
	return v
}

func Language() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig.Language
}
