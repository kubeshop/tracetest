package lintern_plugins_standards_rules

type RequiredAttributesMap map[string][]string

func NewRequiredAttributesMap(data map[string][]string) RequiredAttributesMap {
	return data
}

func (ram RequiredAttributesMap) Types() []string {
	spanTypes := make([]string, len(ram))

	for spanType := range ram {
		spanTypes = append(spanTypes, spanType)
	}

	return spanTypes
}

func (ram RequiredAttributesMap) TypeAttributes(spanType string) []string {
	attributes, found := ram[spanType]

	if !found {
		return []string{}
	}

	return attributes
}
