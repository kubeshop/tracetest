package utils

import (
	"fmt"
	"io"
	URL "net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func StringToIOReader(s string) io.Reader {
	return strings.NewReader(s)
}

func IOReadCloserToString(r io.ReadCloser) string {
	b, _ := io.ReadAll(r)
	return string(b)
}

func StringReferencesFile(path string) bool {
	// for the current working dir, check if the file exists
	// by finding its absolute path and executing a stat command

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	info, err := os.Stat(absolutePath)
	if err != nil {
		return false
	}

	// if the string is empty the absolute path will the entire dir
	// otherwise the user also could send a directory by mistake
	return info != nil && !info.IsDir()
}

func OpenBrowser(url string) error {
	_, err := URL.Parse(url)
	if err != nil {
		return err
	}

	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}
