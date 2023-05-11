package global_parameters

type Global struct {
	Verbose    bool
	ConfigFile string
	Output     string

	// overrides
	OverrideEndpoint string
}

func NewGlobal() *Global {
	return &Global{}
}

func (p *Global) Validate() error {
	return nil
}

func (p *Global) Warnings() string {
	return ""
}
