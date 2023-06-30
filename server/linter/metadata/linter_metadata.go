package metadata

type PluginMetadata struct {
	Slug           string
	Name           string
	Description    string
	DefaultEnabled bool
	Rules          []RuleMetadata
}

type RuleMetadata struct {
	Slug              string
	Name              string
	ErrorDescription  string
	Description       string
	Tips              []string
	DefaultWeight     int
	DefaultErrorLevel string
}

var (
	ErrorLevelWarning  string = "warning"
	ErrorLevelError    string = "error"
	ErrorLevelDisabled string = "disabled"

	Plugins = []PluginMetadata{
		StandardsPlugin,
		CommonPlugin,
		SecurityPlugin,
	}

	AvailablePlugins = []string{StandardsPlugin.Slug, CommonPlugin.Slug, SecurityPlugin.Slug}

	StandardsPlugin = PluginMetadata{
		Slug:           "standards",
		Name:           "OTel Semantic Conventions",
		Description:    "Enforce standards for spans and attributes",
		DefaultEnabled: true,
		Rules: []RuleMetadata{
			EnsureSpanNamingRule,
			RequiredAttributesRule,
			EnsureAttributeNamingRule,
			NotEmptyAttributesRule,
		},
	}

	CommonPlugin = PluginMetadata{
		Slug:           "common",
		Name:           "Common problems",
		Description:    "Helps you find common problems with your application",
		DefaultEnabled: true,
		Rules: []RuleMetadata{
			EnforceDnsRule,
		},
	}

	SecurityPlugin = PluginMetadata{
		Slug:           "security",
		Name:           "Security",
		Description:    "Enforce security for spans and attributes",
		DefaultEnabled: true,
		Rules: []RuleMetadata{
			EnforceHttpsProtocolRule,
			EnsuresNoApiKeyLeakRule,
		},
	}

	EnforceDnsRule = RuleMetadata{
		Slug:              "enforce_dns",
		Name:              "Enforce DNS Over IP usage",
		Description:       "Enforce DNS usage over IP addresses",
		ErrorDescription:  "The following attributes are using IP addresses instead of DNS:",
		Tips:              []string{},
		DefaultWeight:     100,
		DefaultErrorLevel: "error",
	}

	EnforceHttpsProtocolRule = RuleMetadata{
		Slug:              "enforce_https_protocol",
		Name:              "Enforce HTTPS protocol",
		Description:       "Ensure all request use https",
		ErrorDescription:  "The following spans are using http protocol:",
		Tips:              []string{},
		DefaultWeight:     30,
		DefaultErrorLevel: "error",
	}

	EnsureAttributeNamingRule = RuleMetadata{
		Slug:             "span_attribute_naming",
		Name:             "Attribute Naming",
		Description:      "Ensure all attributes follow the naming convention",
		ErrorDescription: "The following attributes do not follow the naming convention:",
		Tips: []string{
			"You should always add namespaces to your span names to ensure they will not be overwritten",
			"Use snake_case to separate multi-words. Ex: http.status_code instead of http.statusCode",
		},
		DefaultWeight:     25,
		DefaultErrorLevel: "error",
	}

	EnsureSpanNamingRule = RuleMetadata{
		Slug:              "span_naming",
		Name:              "Span Name Convention",
		Description:       "Ensure all spans follow the naming convention",
		ErrorDescription:  "",
		Tips:              []string{},
		DefaultWeight:     25,
		DefaultErrorLevel: "error",
	}

	EnsuresNoApiKeyLeakRule = RuleMetadata{
		Slug:              "no_api_key_leak",
		Name:              "No API Key Leak",
		Description:       "Ensure no API keys are leaked in http headers",
		ErrorDescription:  "The following attributes are exposing API keys:",
		Tips:              []string{},
		DefaultWeight:     70,
		DefaultErrorLevel: "error",
	}

	NotEmptyAttributesRule = RuleMetadata{
		Slug:              "not_empty_attributes",
		Name:              "Not Empty Attributes",
		Description:       "Does not allow empty attribute values in any span",
		ErrorDescription:  "The following attributes are empty:",
		Tips:              []string{"Empty attributes don't provide any information about the operation and should be removed"},
		DefaultWeight:     25,
		DefaultErrorLevel: "error",
	}

	RequiredAttributesRule = RuleMetadata{
		Slug:              "required_attributes",
		Name:              "Required Attributes By Span Type",
		Description:       "Ensure all required attributes are present",
		ErrorDescription:  "This span is missing the following required attributes:",
		Tips:              []string{"This rule checks if all required attributes are present in spans of given type"},
		DefaultWeight:     25,
		DefaultErrorLevel: "error",
	}
)
