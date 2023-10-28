package config

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"sync"

	errorz "github.com/kunitsucom/util.go/errors"
	"github.com/kunitsucom/util.go/flagenv"

	"github.com/kunitsucom/ddlgen/internal/logs"
)

// Use a structure so that settings can be backed up.
//
//nolint:tagliatelle
type config struct {
	Debug       *bool   `json:"debug"`
	Timestamp   *string `json:"timestamp"`
	Language    *string `json:"language"`
	Dialect     *string `json:"dialect"`
	Source      *string `json:"source"`
	Destination *string `json:"destination"`
	// Golang
	ColumnKeyGo *string `json:"column_key_go"`
	DDLKeyGo    *string `json:"ddl_key_go"`
}

//nolint:gochecknoglobals
var (
	globalConfig   *config
	globalConfigMu sync.RWMutex
)

func MustLoad(ctx context.Context) (rollback func()) {
	rollback, err := Load(ctx)
	if err != nil {
		err = errorz.Errorf("Load: %w", err)
		panic(err)
	}
	return rollback
}

func Load(ctx context.Context) (rollback func(), err error) {
	globalConfigMu.Lock()
	defer globalConfigMu.Unlock()
	backup := globalConfig

	cfg, err := load(ctx)
	if err != nil {
		return nil, errorz.Errorf("load: %w", err)
	}

	globalConfig = cfg

	rollback = func() {
		globalConfigMu.Lock()
		defer globalConfigMu.Unlock()
		globalConfig = backup
	}

	return rollback, nil
}

// MEMO: Since there is a possibility of returning some kind of error in the future, the signature is made to return an error.
func load(ctx context.Context) (cfg *config, err error) { //nolint:unparam
	fe := flagenv.NewFlagEnvSet(filepath.Base(os.Args[0]), flag.ContinueOnError)

	c := &config{
		Debug:       loadDebug(ctx, fe),
		Timestamp:   loadTimestamp(ctx, fe),
		Language:    loadLanguage(ctx, fe),
		Dialect:     loadDialect(ctx, fe),
		Source:      loadSource(ctx, fe),
		Destination: loadDestination(ctx, fe),
		ColumnKeyGo: loadColumnKeyGo(ctx, fe),
		DDLKeyGo:    loadDDLKeyGo(ctx, fe),
	}

	if err := fe.Parse(os.Args[1:]); err != nil {
		return nil, errorz.Errorf("flagenvSet.Parse: %w", err)
	}

	if err := json.NewEncoder(logs.Debug).Encode(c); err != nil {
		logs.Debug.Printf("config: %#v", c)
	}

	return c, nil
}
