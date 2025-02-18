package utils

import (
	libPath "path"
	"strings"
)

// CleanPath clean path
// example:
//
//	/a////b => /a/b
//	/a/b////c/ => /a/b/c/
func CleanPath(path string) string {
	suffix := ""
	if strings.HasSuffix(path, "/") {
		suffix = "/"
	}
	return libPath.Clean(path) + suffix
}
