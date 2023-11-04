package ddlgengo

import (
	"context"
	"fmt"
	goast "go/ast"
	"go/token"
	"regexp"
	"strings"
	"sync"

	errorz "github.com/kunitsucom/util.go/errors"
	filepathz "github.com/kunitsucom/util.go/path/filepath"

	"github.com/kunitsucom/ddlgen/internal/config"
	"github.com/kunitsucom/ddlgen/internal/logs"
	apperr "github.com/kunitsucom/ddlgen/pkg/errors"
)

type ddlSource struct {
	// TypeSpec is used to guess the table name if the CREATE TABLE annotation is not found.
	TypeSpec *goast.TypeSpec
	// StructType is used to determine the column name. If the tag specified by --column-key-go is not found, the field name is used.
	StructType   *goast.StructType
	CommentGroup *goast.CommentGroup
}

//nolint:gochecknoglobals
var (
	_DDLKeyGoCommentLineRegex     *regexp.Regexp
	_DDLKeyGoCommentLineRegexOnce sync.Once
)

const (
	//	                                       _______________________ <- 1. comment prefix
	//	                                                                __ <- 2. tag name
	//	                                                                        __ <- 3. tag value
	//	                                                                                ___ <- 4. comment suffix
	_DDLKeyGoCommentLineRegexFormat       = `^(\s*//+\s*|\s*/\*\s*|\s*)(%s):\s*(.*)?\s*(\*/)?`
	_DDLKeyGoCommentLineRegexContentIndex = /*                                  ^^ */ 3
)

func ddlKeyGoCommentLineRegex() *regexp.Regexp {
	_DDLKeyGoCommentLineRegexOnce.Do(func() {
		_DDLKeyGoCommentLineRegex = regexp.MustCompile(fmt.Sprintf(_DDLKeyGoCommentLineRegexFormat, config.DDLKeyGo()))
	})
	return _DDLKeyGoCommentLineRegex
}

//
//nolint:cyclop
func extractDDLSourceFromDDLKeyGo(_ context.Context, fset *token.FileSet, f *goast.File) ([]*ddlSource, error) {
	ddlSrc := make([]*ddlSource, 0)

	for commentedNode, commentGroups := range goast.NewCommentMap(fset, f, f.Comments) {
		for _, commentGroup := range commentGroups {
			commentLines := strings.Split(strings.TrimSuffix(commentGroup.Text(), "\n"), "\n")
		CommentGroupLoop:
			for _, commentLine := range commentLines {
				logs.Trace.Printf("commentLine=%s: %s", filepathz.Short(fset.Position(commentGroup.Pos()).String()), commentLine)
				// NOTE: If the comment line matches the DDLKeyGo, it is assumed to be a comment line for the struct.
				if matches := ddlKeyGoCommentLineRegex().FindStringSubmatch(commentLine); len(matches) > _DDLKeyGoCommentLineRegexContentIndex {
					r := &ddlSource{
						CommentGroup: commentGroup,
					}
					goast.Inspect(commentedNode, func(node goast.Node) bool {
						switch n := node.(type) {
						case *goast.TypeSpec:
							r.TypeSpec = n
							// NOTE: Continue searching deeper until StructType appears.
							return true
						case *goast.StructType:
							r.StructType = n
							return false
						default:
							return true
						}
					})
					ddlSrc = append(ddlSrc, r)
					break CommentGroupLoop // NOTE: There may be multiple "DDLKeyGo"s in the same commentGroup, so once you find the first one, break.
				} else if len(matches) > 0 {
					logs.Warn.Printf("commentLine=%s: %v", commentLine, apperr.ErrUnknownError)
				}
			}
		}
	}

	if len(ddlSrc) == 0 {
		return nil, errorz.Errorf("ddl-key-go=%s: %w", config.DDLKeyGo(), apperr.ErrDDLKeyGoNotFoundInSource)
	}

	return ddlSrc, nil
}
