package config

import (
	"context"

	"github.com/kunitsucom/util.go/flagenv"
)

const (
	_EnvKeyDDLKeyGo  = "DDL_KEY_GO"
	_DefaultDDLKeyGo = "ddlgen"
)

func loadDDLKeyGo(_ context.Context, fes *flagenv.FlagEnvSet) *string {
	v := fes.String("ddl-key-go", _EnvKeyDDLKeyGo, _DefaultDDLKeyGo, "DDL Go struct tag key (default: ddl)")
	return v
}

func DDLKeyGo() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return *globalConfig.DDLKeyGo
}
