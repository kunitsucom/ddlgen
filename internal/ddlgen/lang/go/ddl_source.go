package ddlgengo

import goast "go/ast"

type ddlSource struct {
	TypeSpec     *goast.TypeSpec
	StructType   *goast.StructType
	CommentGroup *goast.CommentGroup
}
