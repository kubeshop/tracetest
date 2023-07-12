package analyzer_test

import (
	"database/sql"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
)

func TestLinterResource(t *testing.T) {
	rmtests.TestResourceType(t, rmtests.ResourceTypeTest{
		ResourceTypeSingular: analyzer.ResourceName,
		ResourceTypePlural:   analyzer.ResourceNamePlural,
		RegisterManagerFn: func(router *mux.Router, db *sql.DB) resourcemanager.Manager {
			repo := analyzer.NewRepository(db)

			manager := resourcemanager.New[analyzer.Linter](
				analyzer.ResourceName,
				analyzer.ResourceNamePlural,
				repo,
				resourcemanager.DisableDelete(),
			)
			manager.RegisterRoutes(router)

			return manager
		},
		SampleJSON: `{
			"type": "Analyzer",
			"spec": {
				"id": "current",
				"name": "analyzer",
				"enabled": true,
				"minimumScore": 0,
				"plugins": [
					{
						"id": "standards",
						"enabled": true,
						"rules": [
							{
								"id": "span-naming",
								"weight": 25,
								"errorLevel": "error",
								"name": "Span Naming",
								"errorDescription": "",
								"description": "Enforce span names that identify a class of Spans",
								"tips": []
							},
							{
								"id": "required-attributes",
								"weight": 25,
								"errorLevel": "error",
								"name": "Required Attributes",
								"errorDescription": "This span is missing the following required attributes:",
								"description": "Enforce required attributes by span type",
								"tips": [
									"This rule checks if all required attributes are present in spans of given type"
								]
							},
							{
								"id": "attribute-naming",
								"weight": 25,
								"errorLevel": "error",
								"name": "Attribute Naming",
								"errorDescription": "The following attributes do not follow the naming convention:",
								"description": "Enforce attribute keys to follow common specifications",
								"tips": [
									"You should always add namespaces to your span names to ensure they will not be overwritten",
									"Use snake_case to separate multi-words. Ex: http.status_code instead of http.statusCode"
								]
							},
							{
								"id": "no-empty-attributes",
								"weight": 25,
								"errorLevel": "error",
								"name": "No Empty Attributes",
								"errorDescription": "The following attributes are empty:",
								"description": "Disallow empty attribute values",
								"tips": [
									"Empty attributes don't provide any information about the operation and should be removed"
								]
							}
						],
						"name": "OTel Semantic Conventions",
						"description": "Enforce trace standards following OTel Semantic Conventions"
					},
					{
						"id": "common",
						"enabled": true,
						"rules": [
							{
								"id": "prefer-dns",
								"weight": 100,
								"errorLevel": "error",
								"name": "Prefer DNS",
								"errorDescription": "The following attributes are using IP addresses instead of DNS:",
								"description": "Enforce usage of DNS instead of IP addresses",
								"tips": []
							}
						],
						"name": "Common Problems",
						"description": "Help you find common mistakes with your application"
					},
					{
						"id": "security",
						"enabled": true,
						"rules": [
							{
								"id": "secure-https-protocol",
								"weight": 30,
								"errorLevel": "error",
								"name": "Secure HTTPS Protocol",
								"errorDescription": "The following attributes are using insecure http protocol:",
								"description": "Enforce usage of secure protocol for HTTP server spans",
								"tips": []
							},
							{
								"id": "no-api-key-leak",
								"weight": 70,
								"errorLevel": "error",
								"name": "No API Key Leak",
								"errorDescription": "The following attributes are exposing API keys:",
								"description": "Disallow leaked API keys for HTTP spans",
								"tips": []
							}
						],
						"name": "Security",
						"description": "Help you find security problems with your application"
					}
				]
			}
		}`,
		SampleJSONUpdated: `{
			"type": "Analyzer",
			"spec": {
				"id": "current",
				"name": "analyzer",
				"enabled": true,
				"minimumScore": 0,
				"plugins": [
					{
						"id": "standards",
						"enabled": true,
						"rules": [
							{
								"id": "span-naming",
								"weight": 25,
								"errorLevel": "error",
								"name": "Span Naming",
								"errorDescription": "",
								"description": "Enforce span names that identify a class of Spans",
								"tips": []
							},
							{
								"id": "required-attributes",
								"weight": 25,
								"errorLevel": "error",
								"name": "Required Attributes",
								"errorDescription": "This span is missing the following required attributes:",
								"description": "Enforce required attributes by span type",
								"tips": [
									"This rule checks if all required attributes are present in spans of given type"
								]
							},
							{
								"id": "attribute-naming",
								"weight": 25,
								"errorLevel": "error",
								"name": "Attribute Naming",
								"errorDescription": "The following attributes do not follow the naming convention:",
								"description": "Enforce attribute keys to follow common specifications",
								"tips": [
									"You should always add namespaces to your span names to ensure they will not be overwritten",
									"Use snake_case to separate multi-words. Ex: http.status_code instead of http.statusCode"
								]
							},
							{
								"id": "no-empty-attributes",
								"weight": 25,
								"errorLevel": "error",
								"name": "No Empty Attributes",
								"errorDescription": "The following attributes are empty:",
								"description": "Disallow empty attribute values",
								"tips": [
									"Empty attributes don't provide any information about the operation and should be removed"
								]
							}
						],
						"name": "OTel Semantic Conventions",
						"description": "Enforce trace standards following OTel Semantic Conventions"
					},
					{
						"id": "common",
						"enabled": true,
						"rules": [
							{
								"id": "prefer-dns",
								"weight": 100,
								"errorLevel": "error",
								"name": "Prefer DNS",
								"errorDescription": "The following attributes are using IP addresses instead of DNS:",
								"description": "Enforce usage of DNS instead of IP addresses",
								"tips": []
							}
						],
						"name": "Common Problems",
						"description": "Help you find common mistakes with your application"
					},
					{
						"id": "security",
						"enabled": true,
						"rules": [
							{
								"id": "secure-https-protocol",
								"weight": 30,
								"errorLevel": "error",
								"name": "Secure HTTPS Protocol",
								"errorDescription": "The following spans are using http protocol:",
								"description": "Enforce usage of secure protocol for HTTP server spans",
								"tips": []
							},
							{
								"id": "no-api-key-leak",
								"weight": 70,
								"errorLevel": "error",
								"name": "No API Key Leak",
								"errorDescription": "The following attributes are exposing API keys:",
								"description": "Disallow leaked API keys for HTTP spans",
								"tips": []
							}
						],
						"name": "Security",
						"description": "Help you find security problems with your application"
					}
				]
			}
		}`,
	},
		rmtests.ExcludeOperations(
			rmtests.OperationGetNotFound,
			rmtests.OperationUpdateNotFound,
			rmtests.OperationListSortSuccess,
			rmtests.OperationListNoResults,
		),
	)
}
