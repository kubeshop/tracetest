package analyzer

var (
	// plugins
	StandardsID = "standards"
	CommonID    = "common"
	SecurityID  = "security"

	// rules
	EnsureSpanNamingRuleID      string = "span-naming"
	RequiredAttributesRuleID    string = "required-attributes"
	EnsureAttributeNamingRuleID string = "attribute-naming"
	NotEmptyAttributesRuleID    string = "no-empty-attributes"
	EnforceDnsRuleID            string = "prefer-dns"
	EnforceHttpsProtocolRuleID  string = "secure-https-protocol"
	EnsuresNoApiKeyLeakRuleID   string = "no-api-key-leak"

	ErrorLevelWarning  string = "warning"
	ErrorLevelError    string = "error"
	ErrorLevelDisabled string = "disabled"

	SortWeight = map[string]int{
		ErrorLevelError:    2,
		ErrorLevelWarning:  1,
		ErrorLevelDisabled: 0,
	}

	DefaultPlugins = []LinterPlugin{
		StandardsPlugin,
		CommonPlugin,
		SecurityPlugin,
	}

	AvailablePlugins = []string{StandardsPlugin.ID, CommonPlugin.ID, SecurityPlugin.ID}

	// standards
	StandardsPlugin = LinterPlugin{
		ID:          StandardsID,
		Name:        "OTel Semantic Conventions",
		Description: "Enforce trace standards following OTel Semantic Conventions",
		Enabled:     true,
		Rules: []LinterRule{
			EnsureSpanNamingRule,
			RequiredAttributesRule,
			EnsureAttributeNamingRule,
			NotEmptyAttributesRule,
		},
	}

	EnsureSpanNamingRule = LinterRule{
		ID:               EnsureSpanNamingRuleID,
		Name:             "Span Naming",
		Description:      "Enforce span names that identify a class of Spans",
		ErrorDescription: "",
		Tips:             []string{},
		Weight:           25,
		ErrorLevel:       "error",
	}

	RequiredAttributesRule = LinterRule{
		ID:               RequiredAttributesRuleID,
		Name:             "Required Attributes",
		Description:      "Enforce required attributes by span type",
		ErrorDescription: "This span is missing the following required attributes:",
		Tips:             []string{"This rule checks if all required attributes are present in spans of given type"},
		Weight:           25,
		ErrorLevel:       "error",
	}

	EnsureAttributeNamingRule = LinterRule{
		ID:               EnsureAttributeNamingRuleID,
		Name:             "Attribute Naming",
		Description:      "Enforce attribute keys to follow common specifications",
		ErrorDescription: "The following attributes do not follow the naming convention:",
		Tips: []string{
			"You should always add namespaces to your span names to ensure they will not be overwritten",
			"Use snake_case to separate multi-words. Ex: http.status_code instead of http.statusCode",
		},
		Weight:     25,
		ErrorLevel: "error",
	}

	NotEmptyAttributesRule = LinterRule{
		ID:               NotEmptyAttributesRuleID,
		Name:             "No Empty Attributes",
		Description:      "Disallow empty attribute values",
		ErrorDescription: "The following attributes are empty:",
		Tips:             []string{"Empty attributes don't provide any information about the operation and should be removed"},
		Weight:           25,
		ErrorLevel:       "error",
	}

	// common
	CommonPlugin = LinterPlugin{
		ID:          CommonID,
		Name:        "Common Problems",
		Description: "Help you find common mistakes with your application",
		Enabled:     true,
		Rules: []LinterRule{
			EnforceDnsRule,
		},
	}

	EnforceDnsRule = LinterRule{
		ID:               EnforceDnsRuleID,
		Name:             "Prefer DNS",
		Description:      "Enforce usage of DNS instead of IP addresses",
		ErrorDescription: "The following attributes are using IP addresses instead of DNS:",
		Tips:             []string{},
		Weight:           100,
		ErrorLevel:       "error",
	}

	// security
	SecurityPlugin = LinterPlugin{
		ID:          SecurityID,
		Name:        "Security",
		Description: "Help you find security problems with your application",
		Enabled:     true,
		Rules: []LinterRule{
			EnforceHttpsProtocolRule,
			EnsuresNoApiKeyLeakRule,
		},
	}

	EnforceHttpsProtocolRule = LinterRule{
		ID:               EnforceHttpsProtocolRuleID,
		Name:             "Secure HTTPS Protocol",
		Description:      "Enforce usage of secure protocol for HTTP server spans",
		ErrorDescription: "The following attributes are using insecure http protocol:",
		Tips:             []string{},
		Weight:           30,
		ErrorLevel:       "error",
	}

	EnsuresNoApiKeyLeakRule = LinterRule{
		ID:               EnsuresNoApiKeyLeakRuleID,
		Name:             "No API Key Leak",
		Description:      "Disallow leaked API keys for HTTP spans",
		ErrorDescription: "The following attributes are exposing API keys:",
		Tips:             []string{},
		Weight:           70,
		ErrorLevel:       "error",
	}
)
