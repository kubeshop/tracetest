package rules

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/model"
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

func (r requiredAttributesRule) validateSpan(span *model.Span) model.Result {
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

	return model.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}
}

func (r requiredAttributesRule) validateHttpSpan(span *model.Span) model.Result {
	missingAttrs := r.getMissingAttrs(span, httpAttr, "http")
	result := model.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}

	if span.Kind == model.SpanKindClient {
		missingAttrs = append(missingAttrs, r.getMissingAttrs(span, httpAttrClient, "http")...)
	} else if span.Kind == model.SpanKindServer {
		missingAttrs = append(missingAttrs, r.getMissingAttrs(span, httpAttrServer, "http")...)
	}

	if len(missingAttrs) > 0 {
		result.Passed = false
		result.Errors = missingAttrs
	}

	return result
}

func (r requiredAttributesRule) getMissingAttrs(span *model.Span, matchingAttrList []string, spanType string) []model.Error {
	missingAttributes := make([]model.Error, 0)
	for _, requiredAttribute := range matchingAttrList {
		if _, attributeExists := span.Attributes[requiredAttribute]; !attributeExists {
			missingAttributes = append(missingAttributes, model.Error{
				Error: "missing_attribute_error",
				Value: requiredAttribute,
				// Expected:    requiredAttribute,
				Description: fmt.Sprintf(`Attribute "%s" is missing from span of type "%s"`, requiredAttribute, spanType),
			})
		}
	}

	return missingAttributes
}

func (r requiredAttributesRule) validateDatabaseSpan(span *model.Span) model.Result {
	missingAttrs := r.getMissingAttrs(span, databaseAttr, "database")
	result := model.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}

	if len(missingAttrs) > 0 {
		result.Passed = false
		result.Errors = missingAttrs
	}

	return result
}

func (r requiredAttributesRule) validateRPCSpan(span *model.Span) model.Result {
	missingAttrs := r.getMissingAttrs(span, rpcAttr, "rpc")
	result := model.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}

	if len(missingAttrs) > 0 {
		result.Passed = false
		result.Errors = missingAttrs
	}

	return result
}

func (r requiredAttributesRule) validateMessagingSpan(span *model.Span) model.Result {
	missingAttrs := r.getMissingAttrs(span, messagingAttr, "messaging")
	result := model.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}

	if len(missingAttrs) > 0 {
		result.Passed = false
		result.Errors = missingAttrs
	}

	return result
}

func (r requiredAttributesRule) validateFaasSpan(span *model.Span) model.Result {
	missingAttrs := make([]model.Error, 0)
	result := model.Result{
		Passed: true,
		SpanID: span.ID.String(),
	}

	if span.Kind == model.SpanKindClient {
		missingAttrs = r.getMissingAttrs(span, faasAttrClient, "faas")
	} else if span.Kind == model.SpanKindServer {
		missingAttrs = r.getMissingAttrs(span, faasAttrServer, "faas")
	}

	if len(missingAttrs) > 0 {
		result.Passed = false
		result.Errors = missingAttrs
	}

	return result
}
