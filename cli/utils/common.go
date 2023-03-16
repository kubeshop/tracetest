package utils

import (
	"io"
	"strings"
)

func StringToIOReader(s string) io.Reader {
	return strings.NewReader(s)
}

func IOReadCloserToString(r io.ReadCloser) string {
	b, _ := io.ReadAll(r)
	return string(b)
}
