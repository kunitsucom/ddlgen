package util

import "regexp"

type StmtRegex struct {
	Regex *regexp.Regexp
	Index int
}

//nolint:gochecknoglobals
var (
	RegexStmtCreateTable = StmtRegex{
		Regex: regexp.MustCompile(`\s*table\s*:\s*((CREATE\s+TABLE\s+)?.*)`),
		Index: 1,
	}
	RegexStmtCreateTableConstraint = StmtRegex{
		Regex: regexp.MustCompile(`\s*constraint\s*:\s*(.*)`),
		Index: 1,
	}
	RegexStmtCreateTableOptions = StmtRegex{
		Regex: regexp.MustCompile(`\s*options\s*:\s*(.*)`),
		Index: 1,
	}
)
