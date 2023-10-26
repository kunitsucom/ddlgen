package ast

import (
	"regexp"
)

var _ Stmt = (*CreateTableStmt)(nil)

type CreateTableStmt struct {
	SourceFile  string
	SourceLine  int
	Comments    []string                 // -- <Comment>
	CreateTable string                   // CREATE TABLE [IF NOT EXISTS] <Table>
	Columns     []*CreateTableColumn     // ( <Column>, ...
	Constraints []*CreateTableConstraint // <Constraint> )
	Options     []*CreateTableOption     // <Options>;
}

func (stmt *CreateTableStmt) GetSourceFile() string {
	return stmt.SourceFile
}

func (stmt *CreateTableStmt) GetSourceLine() int {
	return stmt.SourceLine
}

func (*CreateTableStmt) stmt() {}

var createTableRegex = regexp.MustCompile(`\s*CREATE\s+TABLE\s+(IF\s+NOT\s+EXISTS\s+)?([^\s]+)`)

func (stmt *CreateTableStmt) SetCreateTable(createTable string) {
	if len(createTableRegex.FindStringSubmatch(createTable)) > 2 {
		stmt.CreateTable = createTable
		return
	}

	stmt.CreateTable = "CREATE TABLE " + createTable
}

type CreateTableColumn struct {
	Comments       []string
	Column         string
	TypeConstraint string
}

type CreateTableConstraint struct {
	Comments   []string
	Constraint string
}

type CreateTableOption struct {
	Comments []string
	Option   string
}
