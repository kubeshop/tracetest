package testutil

import (
	"sigs.k8s.io/yaml"
)

type ContentTypeConverter struct {
	name        string
	contentType string
	fromJSON    func(input string) string
	toJSON      func(input string) string
}

var contentTypeJSON = ContentTypeConverter{
	name:        "json",
	contentType: "application/json",
	fromJSON:    func(jsonString string) string { return jsonString },
	toJSON:      func(jsonString string) string { return jsonString },
}

var contentTypeYAML = ContentTypeConverter{
	name:        "yaml",
	contentType: "text/yaml",
	fromJSON: func(jsonString string) string {
		y, err := yaml.JSONToYAML([]byte(jsonString))
		if err != nil {
			panic(err)
		}
		return string(y)
	},
	toJSON: func(yamlString string) string {
		j, err := yaml.YAMLToJSON([]byte(yamlString))
		if err != nil {
			panic(err)
		}
		return string(j)
	},
}

var contentTypeConverters = []ContentTypeConverter{contentTypeJSON, contentTypeYAML}
