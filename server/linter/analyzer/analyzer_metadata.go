package analyzer

var (
	// plugins
	StandardsId = "standards"
	CommonId    = "common"
	SecurityId  = "security"

	// rules
	EnsureSpanNamingRuleId      string = "span_naming"
	RequiredAttributesRuleId    string = "required_attributes"
	EnsureAttributeNamingRuleId string = "attribute_naming"
	NotEmptyAttributesRuleId    string = "not_empty_attributes"
	EnforceDnsRuleId            string = "enforce_dns"
	EnforceHttpsProtocolRuleId  string = "enforce_https_protocol"
	EnsuresNoApiKeyLeakRuleId   string = "ensures_no_api_key_leak"

	ErrorLevelWarning  string = "warning"
	ErrorLevelError    string = "error"
	ErrorLevelDisabled string = "disabled"

	DefaultPlugins = []LinterPlugin{
		StandardsPlugin,
		CommonPlugin,
		SecurityPlugin,
	}

	AvailablePlugins = []string{StandardsPlugin.Id, CommonPlugin.Id, SecurityPlugin.Id}

	// standards
	StandardsPlugin = LinterPlugin{
		Id:          StandardsId,
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
		Id:               EnsureSpanNamingRuleId,
		Name:             "Span Name Convention",
		Description:      "Ensure all spans follow the naming convention",
		ErrorDescription: "",
		Tips:             []string{},
		Weight:           25,
		ErrorLevel:       "error",
	}

	RequiredAttributesRule = LinterRule{
		Id:               RequiredAttributesRuleId,
		Name:             "Required Attributes By Span Type",
		Description:      "Ensure all required attributes are present",
		ErrorDescription: "This span is missing the following required attributes:",
		Tips:             []string{"This rule checks if all required attributes are present in spans of given type"},
		Weight:           25,
		ErrorLevel:       "error",
	}

	EnsureAttributeNamingRule = LinterRule{
		Id:               EnsureAttributeNamingRuleId,
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
		Id:               NotEmptyAttributesRuleId,
		Name:             "Not Empty Attributes",
		Description:      "Does not allow empty attribute values in any span",
		ErrorDescription: "The following attributes are empty:",
		Tips:             []string{"Empty attributes don't provide any information about the operation and should be removed"},
		Weight:           25,
		ErrorLevel:       "error",
	}

	// common
	CommonPlugin = LinterPlugin{
		Id:          CommonId,
		Name:        "Common problems",
		Description: "Helps you find common problems with your application",
		Enabled:     true,
		Rules: []LinterRule{
			EnforceDnsRule,
		},
	}

	EnforceDnsRule = LinterRule{
		Id:               EnforceDnsRuleId,
		Name:             "Enforce DNS Over IP usage",
		Description:      "Enforce DNS usage over IP addresses",
		ErrorDescription: "The following attributes are using IP addresses instead of DNS:",
		Tips:             []string{},
		Weight:           100,
		ErrorLevel:       "error",
	}

	// security
	SecurityPlugin = LinterPlugin{
		Id:          SecurityId,
		Name:        "Security",
		Description: "Enforce security for spans and attributes",
		Enabled:     true,
		Rules: []LinterRule{
			EnforceHttpsProtocolRule,
			EnsuresNoApiKeyLeakRule,
		},
	}

	EnforceHttpsProtocolRule = LinterRule{
		Id:               EnforceHttpsProtocolRuleId,
		Name:             "Enforce HTTPS protocol",
		Description:      "Ensure all request use https",
		ErrorDescription: "The following spans are using http protocol:",
		Tips:             []string{},
		Weight:           30,
		ErrorLevel:       "error",
	}

	EnsuresNoApiKeyLeakRule = LinterRule{
		Id:               EnsuresNoApiKeyLeakRuleId,
		Name:             "No API Key Leak",
		Description:      "Ensure no API keys are leaked in http headers",
		ErrorDescription: "The following attributes are exposing API keys:",
		Tips:             []string{},
		Weight:           70,
		ErrorLevel:       "error",
	}
)
