package instance

import (
	"github.com/xoscar/xk6-tracetest-tracing/modules/httpClient"
	"github.com/xoscar/xk6-tracetest-tracing/modules/tracetest"
	"go.k6.io/k6/js/modules"
)

var _ modules.Instance = &MainInstance{}

const version = "0.1.0"

type MainInstance struct {
	vu         modules.VU
	httpClient *httpClient.HttpClient
	Tracetest  *tracetest.Tracetest
}

func New(vu modules.VU, tracetest *tracetest.Tracetest) *MainInstance {
	return &MainInstance{
		vu:         vu,
		httpClient: httpClient.New(vu),
		Tracetest:  tracetest,
	}
}

func (i *MainInstance) Exports() modules.Exports {
	return modules.Exports{
		Named: map[string]interface{}{
			"Http":      i.httpClient.Constructor,
			"Tracetest": i.Tracetest.Constructor,
			"version":   version,
		},
	}
}
