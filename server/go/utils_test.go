package openapi_test

import (
	"os"
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/stretchr/testify/assert"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

func TestTransformTrace(t *testing.T) {
	t.Skip("TODO")
	file, err := os.Open("./testdata/out.json")
	assert.NoError(t, err)
	m := jsonpb.Unmarshaler{}

	var td v1.TracesData
	err = m.Unmarshal(file, &td)
	assert.NoError(t, err)
	//TODO: unknown field "stringValue" in v1.AnyValue
}
