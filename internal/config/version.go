package config

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
