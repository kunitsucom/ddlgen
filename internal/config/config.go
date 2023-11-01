package config

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	errorz "github.com/kunitsucom/util.go/errors"
	cliz "github.com/kunitsucom/util.go/exp/cli"

	"github.com/kunitsucom/ddlgen/internal/contexts"
	"github.com/kunitsucom/ddlgen/internal/logs"
)

// Use a structure so that settings can be backed up.
//
//nolint:tagliatelle
type config struct {
	Trace       bool   `json:"trace"`
	Debug       bool   `json:"debug"`
	Timestamp   string `json:"timestamp"`
	Language    string `json:"language"`
	Dialect     string `json:"dialect"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	// Golang
	ColumnKeyGo string `json:"column_key_go"`
	DDLKeyGo    string `json:"ddl_key_go"`
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

const (
	optionTrace       = "trace"
	optionDebug       = "debug"
	optionTimestamp   = "timestamp"
	optionLanguage    = "lang"
	optionDialect     = "dialect"
	optionSource      = "src"
	optionDestination = "dst"
	// Golang
	optionColumnKeyGo = "column-key-go"
	optionDDLKeyGo    = "ddl-key-go"
)

// MEMO: Since there is a possibility of returning some kind of error in the future, the signature is made to return an error.
//
//nolint:funlen
func load(ctx context.Context) (cfg *config, err error) { //nolint:unparam
	cmd := &cliz.Command{
		Name:        "ddlgen",
		Description: "Generate DDL from Go source code",
		Options: []cliz.Option{
			&cliz.BoolOption{
				Name:        optionTrace,
				Environment: "DDLGEN_TRACE",
				Description: "trace mode enabled",
				Default:     cliz.Default(false),
			},
			&cliz.BoolOption{
				Name:        optionDebug,
				Environment: "DDLGEN_DEBUG",
				Description: "debug mode enabled",
				Default:     cliz.Default(false),
			},
			&cliz.StringOption{
				Name:        optionTimestamp,
				Environment: "DDLGEN_TIMESTAMP",
				Description: "timestamp format",
				Default:     cliz.Default(time.Now().Format(time.RFC3339)),
			},
			&cliz.StringOption{
				Name:        optionLanguage,
				Environment: "DDLGEN_LANG",
				Description: "programming language",
			},
			&cliz.StringOption{
				Name:        optionDialect,
				Environment: "DDLGEN_DIALECT",
				Description: "SQL dialect",
			},
			&cliz.StringOption{
				Name:        optionSource,
				Environment: "DDLGEN_SOURCE",
				Description: "source file or directory",
			},
			&cliz.StringOption{
				Name:        optionDestination,
				Environment: "DDLGEN_DESTINATION",
				Description: "destination file or directory",
			},
			// Golang
			&cliz.StringOption{
				Name:        optionColumnKeyGo,
				Environment: "DDLGEN_COLUMN_KEY_GO",
				Description: "column annotation key for Go struct tag",
				Default:     cliz.Default("db"),
			},
			&cliz.StringOption{
				Name:        optionDDLKeyGo,
				Environment: "DDLGEN_DDL_KEY_GO",
				Description: "DDL annotation key for Go struct tag",
				Default:     cliz.Default("ddlgen"),
			},
		},
	}

	if _, err := cmd.Parse(contexts.Args(ctx)); err != nil {
		return nil, errorz.Errorf("cmd.Parse: %w", err)
	}

	c := &config{
		Trace:       loadTrace(ctx, cmd),
		Debug:       loadDebug(ctx, cmd),
		Timestamp:   loadTimestamp(ctx, cmd),
		Language:    loadLanguage(ctx, cmd),
		Dialect:     loadDialect(ctx, cmd),
		Source:      loadSource(ctx, cmd),
		Destination: loadDestination(ctx, cmd),
		ColumnKeyGo: loadColumnKeyGo(ctx, cmd),
		DDLKeyGo:    loadDDLKeyGo(ctx, cmd),
	}

	if c.Debug {
		logs.Debug = logs.NewDebug()
		logs.Trace.Print("debug mode enabled")
	}
	if c.Trace {
		logs.Trace = logs.NewTrace()
		logs.Debug = logs.NewDebug()
		logs.Debug.Print("trace mode enabled")
	}

	if err := json.NewEncoder(logs.Debug).Encode(c); err != nil {
		logs.Debug.Printf("config: %#v", c)
	}

	return c, nil
}
