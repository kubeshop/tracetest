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
								"name": "Span Name Convention",
								"errorDescription": "",
								"description": "Ensure all spans follow the naming convention",
								"tips": []
							},
							{
								"id": "required-attributes",
								"weight": 25,
								"errorLevel": "error",
								"name": "Required Attributes By Span Type",
								"errorDescription": "This span is missing the following required attributes:",
								"description": "Ensure all required attributes are present",
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
								"description": "Ensure all attributes follow the naming convention",
								"tips": [
									"You should always add namespaces to your span names to ensure they will not be overwritten",
									"Use snake_case to separate multi-words. Ex: http.status_code instead of http.statusCode"
								]
							},
							{
								"id": "no-empty-attributes",
								"weight": 25,
								"errorLevel": "error",
								"name": "Not Empty Attributes",
								"errorDescription": "The following attributes are empty:",
								"description": "Does not allow empty attribute values in any span",
								"tips": [
									"Empty attributes don't provide any information about the operation and should be removed"
								]
							}
						],
						"name": "OTel Semantic Conventions",
						"description": "Enforce standards for spans and attributes"
					},
					{
						"id": "common",
						"enabled": true,
						"rules": [
							{
								"id": "prefer-dns",
								"weight": 100,
								"errorLevel": "error",
								"name": "Enforce DNS Over IP usage",
								"errorDescription": "The following attributes are using IP addresses instead of DNS:",
								"description": "Enforce DNS usage over IP addresses",
								"tips": []
							}
						],
						"name": "Common problems",
						"description": "Helps you find common problems with your application"
					},
					{
						"id": "security",
						"enabled": true,
						"rules": [
							{
								"id": "secure-https-protocol",
								"weight": 30,
								"errorLevel": "error",
								"name": "Enforce HTTPS protocol",
								"errorDescription": "The following spans are using http protocol:",
								"description": "Ensure all request use https",
								"tips": []
							},
							{
								"id": "no-api-key-leak",
								"weight": 70,
								"errorLevel": "error",
								"name": "No API Key Leak",
								"errorDescription": "The following attributes are exposing API keys:",
								"description": "Ensure no API keys are leaked in http headers",
								"tips": []
							}
						],
						"name": "Security",
						"description": "Enforce security for spans and attributes"
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
								"name": "Span Name Convention",
								"errorDescription": "",
								"description": "Ensure all spans follow the naming convention",
								"tips": []
							},
							{
								"id": "required-attributes",
								"weight": 25,
								"errorLevel": "error",
								"name": "Required Attributes By Span Type",
								"errorDescription": "This span is missing the following required attributes:",
								"description": "Ensure all required attributes are present",
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
								"description": "Ensure all attributes follow the naming convention",
								"tips": [
									"You should always add namespaces to your span names to ensure they will not be overwritten",
									"Use snake_case to separate multi-words. Ex: http.status_code instead of http.statusCode"
								]
							},
							{
								"id": "no-empty-attributes",
								"weight": 25,
								"errorLevel": "error",
								"name": "Not Empty Attributes",
								"errorDescription": "The following attributes are empty:",
								"description": "Does not allow empty attribute values in any span",
								"tips": [
									"Empty attributes don't provide any information about the operation and should be removed"
								]
							}
						],
						"name": "OTel Semantic Conventions",
						"description": "Enforce standards for spans and attributes"
					},
					{
						"id": "common",
						"enabled": true,
						"rules": [
							{
								"id": "prefer-dns",
								"weight": 100,
								"errorLevel": "error",
								"name": "Enforce DNS Over IP usage",
								"errorDescription": "The following attributes are using IP addresses instead of DNS:",
								"description": "Enforce DNS usage over IP addresses",
								"tips": []
							}
						],
						"name": "Common problems",
						"description": "Helps you find common problems with your application"
					},
					{
						"id": "security",
						"enabled": true,
						"rules": [
							{
								"id": "secure-https-protocol",
								"weight": 30,
								"errorLevel": "error",
								"name": "Enforce HTTPS protocol",
								"errorDescription": "The following spans are using http protocol:",
								"description": "Ensure all request use https",
								"tips": []
							},
							{
								"id": "no-api-key-leak",
								"weight": 70,
								"errorLevel": "error",
								"name": "No API Key Leak",
								"errorDescription": "The following attributes are exposing API keys:",
								"description": "Ensure no API keys are leaked in http headers",
								"tips": []
							}
						],
						"name": "Security",
						"description": "Enforce security for spans and attributes"
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
