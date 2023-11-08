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
	// StructType is used to determine the column name. If the tag specified by --column-tag-go is not found, the field name is used.
	StructType   *goast.StructType
	CommentGroup *goast.CommentGroup
}

//nolint:gochecknoglobals
var (
	_DDLTagGoCommentLineRegex     *regexp.Regexp
	_DDLTagGoCommentLineRegexOnce sync.Once
)

const (
	//	                                       _______________________ <- 1. comment prefix
	//	                                                                __ <- 2. tag name
	//	                                                                        __ <- 3. tag value
	//	                                                                                ___ <- 4. comment suffix
	_DDLTagGoCommentLineRegexFormat       = `^(\s*//+\s*|\s*/\*\s*|\s*)(%s):\s*(.*)?\s*(\*/)?`
	_DDLTagGoCommentLineRegexContentIndex = /*                                  ^^ */ 3
)

func DDLTagGoCommentLineRegex() *regexp.Regexp {
	_DDLTagGoCommentLineRegexOnce.Do(func() {
		_DDLTagGoCommentLineRegex = regexp.MustCompile(fmt.Sprintf(_DDLTagGoCommentLineRegexFormat, config.DDLTagGo()))
	})
	return _DDLTagGoCommentLineRegex
}

//
//nolint:cyclop
func extractDDLSourceFromDDLTagGo(_ context.Context, fset *token.FileSet, f *goast.File) ([]*ddlSource, error) {
	ddlSrc := make([]*ddlSource, 0)

	for commentedNode, commentGroups := range goast.NewCommentMap(fset, f, f.Comments) {
		for _, commentGroup := range commentGroups {
			commentLines := strings.Split(strings.TrimSuffix(commentGroup.Text(), "\n"), "\n")
		CommentGroupLoop:
			for _, commentLine := range commentLines {
				logs.Trace.Printf("commentLine=%s: %s", filepathz.Short(fset.Position(commentGroup.Pos()).String()), commentLine)
				// NOTE: If the comment line matches the DDLTagGo, it is assumed to be a comment line for the struct.
				if matches := DDLTagGoCommentLineRegex().FindStringSubmatch(commentLine); len(matches) > _DDLTagGoCommentLineRegexContentIndex {
					r := &ddlSource{
						CommentGroup: commentGroup,
					}
					goast.Inspect(commentedNode, func(node goast.Node) bool {
						switch n := node.(type) {
						case *goast.TypeSpec:
							r.TypeSpec = n
							switch t := n.Type.(type) {
							case *goast.StructType:
								r.StructType = t
								return false
							default: // noop
							}
						default: // noop
						}
						return true
					})
					ddlSrc = append(ddlSrc, r)
					break CommentGroupLoop // NOTE: There may be multiple "DDLTagGo"s in the same commentGroup, so once you find the first one, break.
				}
			}
		}
	}

	if len(ddlSrc) == 0 {
		return nil, errorz.Errorf("ddl-tag-go=%s: %w", config.DDLTagGo(), apperr.ErrDDLTagGoAnnotationNotFoundInSource)
	}

	return ddlSrc, nil
}
