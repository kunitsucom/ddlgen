package config

import (
	"context"

	"github.com/kunitsucom/util.go/flagenv"
)

const (
	_EnvKeyDialect  = "DIALECT"
	_DefaultDialect = "postgres"
)

func loadDialect(_ context.Context, fes *flagenv.FlagEnvSet) *string {
	v := fes.String("dialect", _EnvKeyDialect, _DefaultDialect, "dialect (default: postgres)")
	return v
}

func Dialect() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return *globalConfig.Dialect
}
