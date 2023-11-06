package config

import (
	"context"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

func loadVersion(_ context.Context, cmd *cliz.Command) bool {
	v, _ := cmd.GetOptionBool(_OptionVersion)
	return v
}

func Version() bool {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig.Version
}

// nolint: deadcode,gochecknoglobals,unused,varcheck
var (
	version   string
	revision  string
	branch    string
	timestamp string
)

func BuildVersion() string   { return version }
func BuildRevision() string  { return revision }
func BuildBranch() string    { return branch }
func BuildTimestamp() string { return timestamp }
