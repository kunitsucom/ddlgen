package ddl

import (
	"regexp"
)

var _ Stmt = (*CreateIndexStmt)(nil)

type CreateIndexStmt struct {
	SourceFile  string
	SourceLine  int
	Comments    []string // -- <Comment>
	CreateIndex string   // CREATE INDEX [IF NOT EXISTS] <Index>
}

func (stmt *CreateIndexStmt) GetSourceFile() string {
	return stmt.SourceFile
}

func (stmt *CreateIndexStmt) GetSourceLine() int {
	return stmt.SourceLine
}

func (*CreateIndexStmt) private() {}

var createIndexRegex = regexp.MustCompile(`\s*CREATE\s+INDEX\s+(IF\s+NOT\s+EXISTS\s+)?([^\s]+)`)

func (stmt *CreateIndexStmt) SetCreateIndex(createIndex string) {
	if len(createIndexRegex.FindStringSubmatch(createIndex)) > 2 {
		stmt.CreateIndex = createIndex
		return
	}

	stmt.CreateIndex = "CREATE INDEX " + createIndex
}
