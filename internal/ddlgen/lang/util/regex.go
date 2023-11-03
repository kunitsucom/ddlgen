package util

import "regexp"

type StmtRegex struct {
	Regex *regexp.Regexp
	Index int
}

//nolint:gochecknoglobals
var (
	StmtRegexCreateTable = StmtRegex{
		Regex: regexp.MustCompile(`\s*tables?\s*:\s*((CREATE\s+TABLE\s+)?.*)`),
		Index: 1,
	}
	StmtRegexCreateTableConstraint = StmtRegex{
		Regex: regexp.MustCompile(`\s*constraints?\s*:\s*(.*)`),
		Index: 1,
	}
	StmtRegexCreateTableOptions = StmtRegex{
		Regex: regexp.MustCompile(`\s*options?\s*:\s*(.*)`),
		Index: 1,
	}
)
