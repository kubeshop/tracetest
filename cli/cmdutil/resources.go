package cmdutil

import (
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
)

func GetResourceTypeFromFile(filePath string) (string, error) {
	f, err := fileutil.Read(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot read file %s: %w", filePath, err)
	}

	return strings.ToLower(f.Type()), nil
}
