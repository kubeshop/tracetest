package rules

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/traces"
)

var (
	httpAttr       = []string{"http.method"}
	httpAttrClient = []string{"http.url", "net.peer.name"}
	httpAttrServer = []string{"http.target", "http.scheme", "net.host.name"}

	databaseAttr   = []string{"db.system"}
	rpcAttr        = []string{"rpc.system", "net.peer.name"}
	messagingAttr  = []string{"messaging.system", "messaging.operation"}
	faasAttrServer = []string{"faas.trigger"}
	faasAttrClient = []string{"faas.invoked_name", "faas.invoked_provider"}
)

func (r requiredAttributesRule) validateSpan(span *traces.Span) analyzer.Result {
	switch span.Attributes.Get("tracetest.span.type") {
	case "http":
		return r.validateHttpSpan(span)
	case "database":
		return r.validateDatabaseSpan(span)
	case "rpc":
		return r.validateRPCSpan(span)
	case "messaging":
		return r.validateMessagingSpan(span)
	case "faas":
		return r.validateFaasSpan(span)
	}

	return analyzer.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}
}

func (r requiredAttributesRule) validateHttpSpan(span *traces.Span) analyzer.Result {
	missingAttrs := r.getMissingAttrs(span, httpAttr, "http")
	result := analyzer.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}

	if span.Kind == traces.SpanKindClient {
		missingAttrs = append(missingAttrs, r.getMissingAttrs(span, httpAttrClient, "http")...)
	} else if span.Kind == traces.SpanKindServer {
		missingAttrs = append(missingAttrs, r.getMissingAttrs(span, httpAttrServer, "http")...)
	}

	if len(missingAttrs) > 0 {
		result.Passed = false
		result.Errors = missingAttrs
	}

	return result
}

func (r requiredAttributesRule) getMissingAttrs(span *traces.Span, matchingAttrList []string, spanType string) []analyzer.Error {
	missingAttributes := make([]analyzer.Error, 0)
	for _, requiredAttribute := range matchingAttrList {
		if _, attributeExists := span.Attributes.GetExists(requiredAttribute); !attributeExists {
			missingAttributes = append(missingAttributes, analyzer.Error{
				Value:       requiredAttribute,
				Description: fmt.Sprintf(`Attribute "%s" is missing from span of type "%s"`, requiredAttribute, spanType),
			})
		}
	}

	return missingAttributes
}

func (r requiredAttributesRule) validateDatabaseSpan(span *traces.Span) analyzer.Result {
	missingAttrs := r.getMissingAttrs(span, databaseAttr, "database")
	result := analyzer.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}

	if len(missingAttrs) > 0 {
		result.Passed = false
		result.Errors = missingAttrs
	}

	return result
}

func (r requiredAttributesRule) validateRPCSpan(span *traces.Span) analyzer.Result {
	missingAttrs := r.getMissingAttrs(span, rpcAttr, "rpc")
	result := analyzer.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}

	if len(missingAttrs) > 0 {
		result.Passed = false
		result.Errors = missingAttrs
	}

	return result
}

func (r requiredAttributesRule) validateMessagingSpan(span *traces.Span) analyzer.Result {
	missingAttrs := r.getMissingAttrs(span, messagingAttr, "messaging")
	result := analyzer.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}

	if len(missingAttrs) > 0 {
		result.Passed = false
		result.Errors = missingAttrs
	}

	return result
}

func (r requiredAttributesRule) validateFaasSpan(span *traces.Span) analyzer.Result {
	missingAttrs := make([]analyzer.Error, 0)
	result := analyzer.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}

	if span.Kind == traces.SpanKindClient {
		missingAttrs = r.getMissingAttrs(span, faasAttrClient, "faas")
	} else if span.Kind == traces.SpanKindServer {
		missingAttrs = r.getMissingAttrs(span, faasAttrServer, "faas")
	}

	if len(missingAttrs) > 0 {
		result.Passed = false
		result.Errors = missingAttrs
	}

	return result
}
