package e2e

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func IsJsonWithFormat(t *testing.T, input string, format interface{}) {
	err := json.Unmarshal([]byte(input), &format)
	assert.NoError(t, err, "string is not a json or isn't in the expected format")
}
