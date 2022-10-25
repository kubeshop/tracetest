package variable_test

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/cli/variable"
	"github.com/kubeshop/tracetest/server/model/yaml"
	"github.com/stretchr/testify/assert"
)

type testingVarProvider struct {
	variables map[string]string
}

func (p testingVarProvider) GetVariable(name string) (string, error) {
	if _, ok := p.variables[name]; !ok {
		return "", fmt.Errorf("could not resolve variable \"%s\"", name)
	}

	return p.variables[name], nil
}

type structWithUnexportedFields struct {
	name    string
	content []byte
	File    yaml.File
}

func TestInjectorWithStruct(t *testing.T) {
	provider := testingVarProvider{
		variables: map[string]string{
			"TRACETEST_URL":       "http://localhost:11633",
			"POKEMON_API_URL":     "http://pokemon.api:11633",
			"EXPECTED_POKEMON_ID": "521",
		},
	}
	injector := variable.NewInjector(variable.WithVariableProvider(provider))

	input := yaml.File{
		Type: "Test",
		Spec: yaml.Test{
			Name: "Test ${TRACETEST_URL}",
			Trigger: yaml.TestTrigger{
				Type: "http",
				HTTPRequest: yaml.HTTPRequest{
					URL:    "${POKEMON_API_URL}",
					Method: "GET",
				},
			},
			Specs: []yaml.TestSpec{
				{
					Selector: "http.url = \"${POKEMON_API_URL}\"",
					Assertions: []string{
						"tracetest.span.duration < 100",
						`tracetest.response.body contains '"id": ${EXPECTED_POKEMON_ID}'`,
					},
				},
			},
		},
	}

	expectedDefinition := yaml.File{
		Type: "Test",
		Spec: yaml.Test{
			Name: "Test http://localhost:11633",
			Trigger: yaml.TestTrigger{
				Type: "http",
				HTTPRequest: yaml.HTTPRequest{
					URL:    "http://pokemon.api:11633",
					Method: "GET",
				},
			},
			Specs: []yaml.TestSpec{
				{
					Selector: "http.url = \"http://pokemon.api:11633\"",
					Assertions: []string{
						"tracetest.span.duration < 100",
						`tracetest.response.body contains '"id": 521'`,
					},
				},
			},
		},
	}

	file := structWithUnexportedFields{
		name:    "haha",
		content: []byte{},
		File:    input,
	}

	expectedFile := structWithUnexportedFields{
		name:    "haha",
		content: []byte{},
		File:    expectedDefinition,
	}

	err := injector.Inject(&file)
	assert.NoError(t, err)

	assert.Equal(t, expectedFile, file)
}
