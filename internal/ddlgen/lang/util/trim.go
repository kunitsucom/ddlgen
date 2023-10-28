package util

import (
	"strings"

	slicez "github.com/kunitsucom/util.go/slices"

	"github.com/kunitsucom/ddlgen/internal/config"
)

func TrimDDLGenCommentElement(stringSlice []string) []string {
	return slicez.Filter(stringSlice, func(_ int, s string) bool {
		return !strings.HasPrefix(s, config.DDLKeyGo())
	})
}

func TrimTailEmptyCommentElement(stringSlice []string) []string {
	if len(stringSlice) == 0 {
		return stringSlice
	}

	if stringSlice[len(stringSlice)-1] == "" {
		return stringSlice[:len(stringSlice)-1]
	}

	return stringSlice
}