package expression_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/traces"
)

func BenchmarkSimpleExpressions(b *testing.B) {
	statement := `1 = 1`

	executor := expression.NewExecutor()

	for i := 0; i < b.N; i++ {
		executor.Statement(statement)
	}
}

func BenchmarkJSONPathExpressions(b *testing.B) {
	statement := `attr:my_json | json_path '[*].id' | count = 3`

	attributeDataStore := expression.AttributeDataStore{
		Span: traces.Span{
			Attributes: traces.Attributes{
				"my_json": getJSON(),
			},
		},
	}
	executor := expression.NewExecutor(attributeDataStore)

	for i := 0; i < b.N; i++ {
		executor.Statement(statement)
	}
}

func getJSON() string {
	return `[
		{
			"id": "0001",
			"type": "donut",
			"name": "Cake",
			"ppu": 0.55,
			"batters":
				{
					"batter":
						[
							{ "id": "1001", "type": "Regular" },
							{ "id": "1002", "type": "Chocolate" },
							{ "id": "1003", "type": "Blueberry" },
							{ "id": "1004", "type": "Devil's Food" }
						]
				},
			"topping":
				[
					{ "id": "5001", "type": "None" },
					{ "id": "5002", "type": "Glazed" },
					{ "id": "5005", "type": "Sugar" },
					{ "id": "5007", "type": "Powdered Sugar" },
					{ "id": "5006", "type": "Chocolate with Sprinkles" },
					{ "id": "5003", "type": "Chocolate" },
					{ "id": "5004", "type": "Maple" }
				]
		},
		{
			"id": "0002",
			"type": "donut",
			"name": "Raised",
			"ppu": 0.55,
			"batters":
				{
					"batter":
						[
							{ "id": "1001", "type": "Regular" }
						]
				},
			"topping":
				[
					{ "id": "5001", "type": "None" },
					{ "id": "5002", "type": "Glazed" },
					{ "id": "5005", "type": "Sugar" },
					{ "id": "5003", "type": "Chocolate" },
					{ "id": "5004", "type": "Maple" }
				]
		},
		{
			"id": "0003",
			"type": "donut",
			"name": "Old Fashioned",
			"ppu": 0.55,
			"batters":
				{
					"batter":
						[
							{ "id": "1001", "type": "Regular" },
							{ "id": "1002", "type": "Chocolate" }
						]
				},
			"topping":
				[
					{ "id": "5001", "type": "None" },
					{ "id": "5002", "type": "Glazed" },
					{ "id": "5003", "type": "Chocolate" },
					{ "id": "5004", "type": "Maple" }
				]
		}
	]`
}
