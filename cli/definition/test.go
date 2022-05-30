package definition

type Test struct {
	Id      string      `yaml:"id"`
	Name    string      `yaml:"name"`
	Trigger TestTrigger `yaml:"trigger"`
}

type TestTrigger struct {
	Type        string      `yaml:"type"`
	HTTPRequest HttpRequest `yaml:"http_request"`
}

type TestDefinition struct {
	Selector   string   `yaml:"selector"`
	Assertions []string `yaml:"assertions"`
}
