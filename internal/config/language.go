package config

import (
	"context"

	"github.com/kunitsucom/util.go/flagenv"
)

const (
	_EnvKeyLanguage  = "LANGUAGE"
	_DefaultLanguage = "go"
)

func loadLanguage(_ context.Context, fes *flagenv.FlagEnvSet) *string {
	v := fes.String("lang", _EnvKeyLanguage, _DefaultLanguage, "language (default: go)")
	return v
}

func Language() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return *globalConfig.Language
}
