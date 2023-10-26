package ddlgengo

import (
	"bufio"
	"bytes"
	goast "go/ast"
	"go/token"
	"strings"

	"github.com/kunitsucom/ddlgen/internal/logs"
)

func dumpDDLSource(fset *token.FileSet, ddlSrc []*ddlSource) {
	for _, r := range ddlSrc {
		logs.Debug.Print("== ddlSource ================================================================================================================================")
		for _, comment := range strings.Split(strings.Trim(r.CommentGroup.Text(), "\n"), "\n") {
			logs.Debug.Printf("r.CommentGroup.Text: %s", comment)
		}
		logs.Debug.Print("-- CommentGroup --------------------------------------------------------------------------------------------------------------------------------")
		commentGroupAST := bytes.NewBuffer(nil)
		goast.Fprint(commentGroupAST, fset, r.CommentGroup, goast.NotNilFilter)
		for s := bufio.NewScanner(bytes.NewBuffer(commentGroupAST.Bytes())); s.Scan(); {
			logs.Debug.Print(s.Text())
		}
		logs.Debug.Print("-- TypeSpec --------------------------------------------------------------------------------------------------------------------------------")
		typeSpecAST := bytes.NewBuffer(nil)
		goast.Fprint(typeSpecAST, fset, r.TypeSpec, goast.NotNilFilter)
		for s := bufio.NewScanner(bytes.NewBuffer(typeSpecAST.Bytes())); s.Scan(); {
			logs.Debug.Print(s.Text())
		}
		logs.Debug.Print("-- StructType --------------------------------------------------------------------------------------------------------------------------------")
		structTypeAST := bytes.NewBuffer(nil)
		goast.Fprint(structTypeAST, fset, r.StructType, goast.NotNilFilter)
		for s := bufio.NewScanner(bytes.NewBuffer(structTypeAST.Bytes())); s.Scan(); {
			logs.Debug.Print(s.Text())
		}
	}
}
