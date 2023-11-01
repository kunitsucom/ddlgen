package ddlgengo

import (
	"context"
	goast "go/ast"
	"go/token"
	"strings"

	errorz "github.com/kunitsucom/util.go/errors"
	filepathz "github.com/kunitsucom/util.go/path/filepath"

	"github.com/kunitsucom/ddlgen/internal/config"
	"github.com/kunitsucom/ddlgen/internal/logs"
	apperr "github.com/kunitsucom/ddlgen/pkg/errors"
)

//nolint:cyclop
func extractDDLSource(_ context.Context, fset *token.FileSet, f *goast.File) ([]*ddlSource, error) {
	var ddlSrc []*ddlSource
	for commentedNode, commentGroups := range goast.NewCommentMap(fset, f, f.Comments) {
		for _, commentGroup := range commentGroups {
		CommentGroupLoop:
			for _, commentLine := range strings.Split(strings.TrimSuffix(commentGroup.Text(), "\n"), "\n") {
				logs.Trace.Printf("commentLine=%s: %s", filepathz.Short(fset.Position(commentGroup.Pos()).String()), commentLine)
				// NOTE: (en) If the comment line matches the DDLKeyGo, it is assumed to be a comment line for the struct.
				if matches := ddlKeyGoCommentLineRegex().FindStringSubmatch(commentLine); len(matches) > 3 {
					r := &ddlSource{
						CommentGroup: commentGroup,
					}
					goast.Inspect(commentedNode, func(node goast.Node) bool {
						switch n := node.(type) {
						case *goast.TypeSpec:
							r.TypeSpec = n
							// NOTE: (en) Continue searching deeper until StructType appears.
							return true
						case *goast.StructType:
							r.StructType = n
							return false
						default:
							return true
						}
					})
					ddlSrc = append(ddlSrc, r)
					break CommentGroupLoop // MEMO: (en) There may be multiple "DDLKeyGo"s in the same commentGroup, so once you find the first one, break.
				} else if len(matches) > 0 {
					logs.Warn.Printf("commentLine=%s: %v", commentLine, apperr.ErrUnknownError)
				}
			}
		}
	}

	if len(ddlSrc) == 0 {
		return nil, errorz.Errorf("searchStructTypeWithDDLKeyGoComment: ddl-key-go=%s: %w", config.DDLKeyGo(), apperr.ErrDDLKeyGoNotFoundInSource)
	}

	return ddlSrc, nil
}
