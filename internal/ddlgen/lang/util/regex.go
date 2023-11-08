package util

import "regexp"

type StmtRegex struct {
	Regex *regexp.Regexp
	Index int
}

//nolint:gochecknoglobals
var (
	StmtRegexCreateTable = StmtRegex{
		Regex: regexp.MustCompile(`\s*\S+:\s*tables?\s*:\s*((CREATE\s+TABLE\s+)?\S+.*)`),
		Index: 1,
	}
	StmtRegexCreateTableConstraint = StmtRegex{
		Regex: regexp.MustCompile(`\s*\S+:\s*constraints?\s*:\s*(\S+.*)`),
		Index: 1,
	}
	StmtRegexCreateTableOptions = StmtRegex{
		Regex: regexp.MustCompile(`\s*\S+:\s*options?\s*:\s*(\S+.*)`),
		Index: 1,
	}
	StmtRegexCreateIndex = StmtRegex{
		Regex: regexp.MustCompile(`\s*\S+:\s*index(es)?\s*:\s*(\S+.*)`),
		Index: 2,
	}
)
