package main

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/agent/ui/dashboard"
)

func main() {
	err := dashboard.StartDashboard(context.Background())
	if err != nil {
		fmt.Println(err.Error())
	}
}
