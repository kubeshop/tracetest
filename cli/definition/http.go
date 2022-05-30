package definition

type HttpRequest struct {
	URL            string             `yaml:"url"`
	Method         string             `yaml:"method"`
	Headers        []HTTPHeader       `yaml:"headers"`
	Authentication HTTPAuthentication `yaml:"authentication"`
	Body           HTTPBody           `yaml:"body"`
}

type HTTPHeader struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type HTTPAuthentication struct {
	Type     string `yaml:"type"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	APIKey   string `yaml:"apikey"`
	Token    string `yaml:"token"`
}

type HTTPBody struct {
	Type string `yaml:"type"`
	Raw  string `yaml:"raw"`
}
