package analyzer

var (
	// plugins
	StandardsSlug = "standards"
	CommonSlug    = "common"
	SecuritySlug  = "security"

	// rules
	EnsureSpanNamingRuleSlug      string = "span_naming"
	RequiredAttributesRuleSlug    string = "required_attributes"
	EnsureAttributeNamingRuleSlug string = "attribute_naming"
	NotEmptyAttributesRuleSlug    string = "not_empty_attributes"
	EnforceDnsRuleSlug            string = "enforce_dns"
	EnforceHttpsProtocolRuleSlug  string = "enforce_https_protocol"
	EnsuresNoApiKeyLeakRuleSlug   string = "ensures_no_api_key_leak"

	ErrorLevelWarning  string = "warning"
	ErrorLevelError    string = "error"
	ErrorLevelDisabled string = "disabled"

	DefaultPlugins = []LinterPlugin{
		StandardsPlugin,
		CommonPlugin,
		SecurityPlugin,
	}

	AvailablePlugins = []string{StandardsPlugin.Slug, CommonPlugin.Slug, SecurityPlugin.Slug}

	// standards
	StandardsPlugin = LinterPlugin{
		Slug:        StandardsSlug,
		Name:        "OTel Semantic Conventions",
		Description: "Enforce standards for spans and attributes",
		Enabled:     true,
		Rules: []LinterRule{
			EnsureSpanNamingRule,
			RequiredAttributesRule,
			EnsureAttributeNamingRule,
			NotEmptyAttributesRule,
		},
	}

	EnsureSpanNamingRule = LinterRule{
		Slug:             EnsureSpanNamingRuleSlug,
		Name:             "Span Name Convention",
		Description:      "Ensure all spans follow the naming convention",
		ErrorDescription: "",
		Tips:             []string{},
		Weight:           25,
		ErrorLevel:       "error",
	}

	RequiredAttributesRule = LinterRule{
		Slug:             RequiredAttributesRuleSlug,
		Name:             "Required Attributes By Span Type",
		Description:      "Ensure all required attributes are present",
		ErrorDescription: "This span is missing the following required attributes:",
		Tips:             []string{"This rule checks if all required attributes are present in spans of given type"},
		Weight:           25,
		ErrorLevel:       "error",
	}

	EnsureAttributeNamingRule = LinterRule{
		Slug:             EnsureAttributeNamingRuleSlug,
		Name:             "Attribute Naming",
		Description:      "Ensure all attributes follow the naming convention",
		ErrorDescription: "The following attributes do not follow the naming convention:",
		Tips: []string{
			"You should always add namespaces to your span names to ensure they will not be overwritten",
			"Use snake_case to separate multi-words. Ex: http.status_code instead of http.statusCode",
		},
		Weight:     25,
		ErrorLevel: "error",
	}

	NotEmptyAttributesRule = LinterRule{
		Slug:             NotEmptyAttributesRuleSlug,
		Name:             "Not Empty Attributes",
		Description:      "Does not allow empty attribute values in any span",
		ErrorDescription: "The following attributes are empty:",
		Tips:             []string{"Empty attributes don't provide any information about the operation and should be removed"},
		Weight:           25,
		ErrorLevel:       "error",
	}

	// common
	CommonPlugin = LinterPlugin{
		Slug:        CommonSlug,
		Name:        "Common problems",
		Description: "Helps you find common problems with your application",
		Enabled:     true,
		Rules: []LinterRule{
			EnforceDnsRule,
		},
	}

	EnforceDnsRule = LinterRule{
		Slug:             EnforceDnsRuleSlug,
		Name:             "Enforce DNS Over IP usage",
		Description:      "Enforce DNS usage over IP addresses",
		ErrorDescription: "The following attributes are using IP addresses instead of DNS:",
		Tips:             []string{},
		Weight:           100,
		ErrorLevel:       "error",
	}

	// security
	SecurityPlugin = LinterPlugin{
		Slug:        SecuritySlug,
		Name:        "Security",
		Description: "Enforce security for spans and attributes",
		Enabled:     true,
		Rules: []LinterRule{
			EnforceHttpsProtocolRule,
			EnsuresNoApiKeyLeakRule,
		},
	}

	EnforceHttpsProtocolRule = LinterRule{
		Slug:             EnforceHttpsProtocolRuleSlug,
		Name:             "Enforce HTTPS protocol",
		Description:      "Ensure all request use https",
		ErrorDescription: "The following spans are using http protocol:",
		Tips:             []string{},
		Weight:           30,
		ErrorLevel:       "error",
	}

	EnsuresNoApiKeyLeakRule = LinterRule{
		Slug:             EnsuresNoApiKeyLeakRuleSlug,
		Name:             "No API Key Leak",
		Description:      "Ensure no API keys are leaked in http headers",
		ErrorDescription: "The following attributes are exposing API keys:",
		Tips:             []string{},
		Weight:           70,
		ErrorLevel:       "error",
	}
)
