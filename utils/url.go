package utils

import (
	"net/url"
	"strings"
)

// CheckURLPath 检查uri是路径还是文件
func CheckURLPath(uri *url.URL) bool {
	if uri.String()[len(uri.String())-1] == '/' {
		return true
	}

	paths := strings.Split(uri.Path, "/")
	if len(paths) < 1 {
		return true
	}

	if strings.Contains(paths[len(paths)-1], ".") {
		return false
	}

	return true
}
