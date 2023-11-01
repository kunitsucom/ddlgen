package config

import (
	"context"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

const (
	_EnvKeyLanguage  = "LANGUAGE"
	_DefaultLanguage = "go"
)

func loadLanguage(_ context.Context, cmd *cliz.Command) string {
	v, _ := cmd.GetStringOption(optionLanguage)
	return v
}

func Language() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig.Language
}
