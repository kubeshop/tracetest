package traces

type ConversionConfig struct {
	timeFields map[string]bool
}

func NewConversionConfig() ConversionConfig {
	return ConversionConfig{
		timeFields: make(map[string]bool, 0),
	}
}

func (c ConversionConfig) AddTimeFields(fields ...string) {
	for _, field := range fields {
		c.timeFields[field] = true
	}
}

func (c ConversionConfig) IsTimeField(field string) bool {
	_, ok := c.timeFields[field]
	return ok
}
