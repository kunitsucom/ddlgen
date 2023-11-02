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
	_OptionTrace = "trace"
	_EnvKeyTrace = "DDLGEN_TRACE"

	_OptionDebug = "debug"
	_EnvKeyDebug = "DDLGEN_DEBUG"

	_OptionTimestamp = "timestamp"
	_EnvKeyTimestamp = "DDLGEN_TIMESTAMP"

	_OptionLanguage = "lang"
	_EnvKeyLanguage = "DDLGEN_LANGUAGE"

	_OptionDialect = "dialect"
	_EnvKeyDialect = "DDLGEN_DIALECT"

	_OptionSource = "src"
	_EnvKeySource = "DDLGEN_SOURCE"

	_OptionDestination = "dst"
	_EnvKeyDestination = "DDLGEN_DESTINATION"

	// Golang

	_OptionColumnKeyGo = "column-key-go"
	_EnvKeyColumnKeyGo = "DDLGEN_COLUMN_KEY_GO"

	_OptionDDLKeyGo = "ddl-key-go"
	_EnvKeyDDLKeyGo = "DDLGEN_DDL_KEY_GO"
)

// MEMO: Since there is a possibility of returning some kind of error in the future, the signature is made to return an error.
//
//nolint:funlen
func load(ctx context.Context) (cfg *config, err error) { //nolint:unparam
	cmd := &cliz.Command{
		Name:        "ddlgen",
		Description: "Generate DDL from annotated source code.",
		Options: []cliz.Option{
			&cliz.BoolOption{
				Name:        _OptionTrace,
				Environment: _EnvKeyTrace,
				Description: "trace mode enabled",
				Default:     cliz.Default(false),
			},
			&cliz.BoolOption{
				Name:        _OptionDebug,
				Environment: _EnvKeyDebug,
				Description: "debug mode",
				Default:     cliz.Default(false),
			},
			&cliz.StringOption{
				Name:        _OptionTimestamp,
				Environment: _EnvKeyTimestamp,
				Description: "code generation timestamp",
				Default:     cliz.Default(time.Now().Format(time.RFC3339)),
			},
			&cliz.StringOption{
				Name:        _OptionLanguage,
				Environment: _EnvKeyLanguage,
				Description: "programming language to generate DDL",
				Default:     cliz.Default("go"),
			},
			&cliz.StringOption{
				Name:        _OptionDialect,
				Environment: _EnvKeyDialect,
				Description: "SQL dialect to generate DDL",
			},
			&cliz.StringOption{
				Name:        _OptionSource,
				Environment: _EnvKeySource,
				Description: "source file or directory",
				Default:     cliz.Default("/dev/stdin"),
			},
			&cliz.StringOption{
				Name:        _OptionDestination,
				Environment: _EnvKeyDestination,
				Description: "destination file or directory",
				Default:     cliz.Default("/dev/stdout"),
			},
			// Golang
			&cliz.StringOption{
				Name:        _OptionColumnKeyGo,
				Environment: _EnvKeyColumnKeyGo,
				Description: "column annotation key for Go struct tag",
				Default:     cliz.Default("db"),
			},
			&cliz.StringOption{
				Name:        _OptionDDLKeyGo,
				Environment: _EnvKeyDDLKeyGo,
				Description: "DDL annotation key for Go struct tag",
				Default:     cliz.Default("ddlgen"),
			},
		},
	}

	if _, _, err := cmd.Parse(contexts.Args(ctx)); err != nil {
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
