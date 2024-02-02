package main

import (
	"context"

	"github.com/kubeshop/tracetest/agent/ui/dashboard"
)

func main() {
	dashboard.StartDashboard(context.Background())
}
