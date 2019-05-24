package utils

import (
	"net/url"
	"testing"
)

func TestCheckURLPath(t *testing.T) {
	var (
		exists bool
		uri    *url.URL
	)

	uri, _ = url.Parse("http://www.google.com")
	exists = CheckURLPath(uri)
	if !exists {
		t.Errorf("%s 是目录", uri.String())
	}

	uri, _ = url.Parse("http://www.google.com/a")
	exists = CheckURLPath(uri)
	if !exists {
		t.Errorf("%s 是目录", uri.String())
	}

	uri, _ = url.Parse("http://www.google.com/a/b")
	exists = CheckURLPath(uri)
	if !exists {
		t.Errorf("%s 是目录", uri.String())
	}

	uri, _ = url.Parse("http://www.google.com/a/b/c.c/")
	exists = CheckURLPath(uri)
	if !exists {
		t.Errorf("%s 是目录", uri.String())
	}

	uri, _ = url.Parse("http://www.google.com/a/b/c.c/abc.html")
	exists = CheckURLPath(uri)
	if exists {
		t.Errorf("%s 是文件", uri.String())
	}

	uri, _ = url.Parse("http://www.google.com/a/b/c.c/r.")
	exists = CheckURLPath(uri)
	if exists {
		t.Errorf("%s 是文件", uri.String())
	}

	uri, _ = url.Parse("http://www.google.com/a/b/c.c/v.mp3")
	exists = CheckURLPath(uri)
	if exists {
		t.Errorf("%s 是文件", uri.String())
	}

	uri, _ = url.Parse("http://www.google.com/a/b/c.c/20190516.log")
	exists = CheckURLPath(uri)
	if exists {
		t.Errorf("%s 是文件", uri.String())
	}
}
