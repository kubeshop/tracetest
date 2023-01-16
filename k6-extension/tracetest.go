package tracetest

import (
	"github.com/kubeshop/tracetest/k6-extension/modules/instance"
	tracetestOutput "github.com/kubeshop/tracetest/k6-extension/modules/output"
	"github.com/kubeshop/tracetest/k6-extension/modules/tracetest"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/output"
)

func init() {
	tracetest := tracetest.New()
	modules.Register("k6/x/tracetest", New(tracetest))

	output.RegisterExtension("xk6-tracetest", func(params output.Params) (output.Output, error) {
		return tracetestOutput.New(params, tracetest)
	})
}

type RootModule struct {
	tracetest *tracetest.Tracetest
}

var _ modules.Module = &RootModule{}

func New(tracetest *tracetest.Tracetest) *RootModule {
	return &RootModule{
		tracetest: tracetest,
	}
}

func (r *RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	r.tracetest.Vu = vu
	return instance.New(vu, r.tracetest)
}
