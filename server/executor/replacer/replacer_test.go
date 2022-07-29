package replacer_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/executor/replacer"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPReplacer(t *testing.T) {
	test := model.Test{
		Name: "A test",
		ServiceUnderTest: model.Trigger{
			Type: model.TriggerTypeHTTP,
			HTTP: &model.HTTPRequest{
				Method: model.HTTPMethodPOST,
				URL:    "http://my-api.com/api/users",
				Headers: []model.HTTPHeader{
					{Key: "X-api-key", Value: "{{ uuid() }}"},
					{Key: "my-key", Value: "my-value"},
				},
				Body: `{ "id": "{{ uuid() }}", "name": "{{ fullName() }}", "age": {{ randomInt(18, 99) }} }`,
			},
		},
	}

	newTest, err := replacer.ReplaceTestPlaceholders(test)
	require.NoError(t, err)
	for _, header := range newTest.ServiceUnderTest.HTTP.Headers {
		assert.NotContains(t, header.Value, "{{")
	}

	assert.NotContains(t, newTest.ServiceUnderTest.HTTP.Body, "{{")
}

func TestGRPCReplacer(t *testing.T) {
	test := model.Test{
		Name: "A test",
		ServiceUnderTest: model.Trigger{
			Type: model.TriggerTypeGRPC,
			GRPC: &model.GRPCRequest{
				ProtobufFile: `syntax = \"proto3\";\r\n\r\noption java_multiple_files = true;\r\noption
				java_outer_classname = \"PokeshopProto\";\r\noption objc_class_prefix = \"PKS\";\r\n\r\npackage
				pokeshop;\r\n\r\nservice Pokeshop {\r\n  rpc getPokemonList (GetPokemonRequest)
				returns (GetPokemonListResponse) {}\r\n  rpc createPokemon (Pokemon) returns
				(Pokemon) {}\r\n  rpc importPokemon (ImportPokemonRequest) returns (ImportPokemonRequest)
				{}\r\n}\r\n\r\nmessage ImportPokemonRequest {\r\n  int32 id = 1;\r\n}\r\n\r\nmessage
				GetPokemonRequest {\r\n  optional int32 skip = 1;\r\n  optional int32 take =
				2;\r\n}\r\n\r\nmessage GetPokemonListResponse {\r\n  repeated Pokemon items
				= 1;\r\n  int32 totalCount = 2;\r\n}\r\n\r\nmessage Pokemon {\r\n  optional
				int32 id = 1;\r\n  string name = 2;\r\n  string type = 3;\r\n  bool isFeatured
				= 4;\r\n  optional string imageUrl = 5;\r\n}`,
				Request: `{"id": {{ randomInt(1, 151) }}}`,
				Method:  "pokeshop.Pokeshop.importPokemon",
				Address: "localhost:8082",
			},
		},
	}

	newTest, err := replacer.ReplaceTestPlaceholders(test)
	require.NoError(t, err)

	assert.NotContains(t, newTest.ServiceUnderTest.GRPC.Request, "{{")
}
