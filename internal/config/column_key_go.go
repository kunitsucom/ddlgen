package config

import (
	"context"

	"github.com/kunitsucom/util.go/flagenv"
)

const (
	_EnvKeyColumnKeyGo  = "COLUMN_KEY_GO"
	_DefaultColumnKeyGo = "db"
)

func loadColumnKeyGo(_ context.Context, fes *flagenv.FlagEnvSet) *string {
	v := fes.String("column-key-go", _EnvKeyColumnKeyGo, _DefaultColumnKeyGo, "DB column Go struct tag key (default: db)")
	return v
}

func ColumnKeyGo() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return *globalConfig.ColumnKeyGo
}
