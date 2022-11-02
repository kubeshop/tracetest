package mappings

import (
	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/traces"
)

type Mappings struct {
	In  Model
	Out OpenAPI
}

func New(tcc traces.ConversionConfig, cr comparator.Registry, tr model.TestRepository) Mappings {
	return Mappings{
		In: Model{
			comparators:           cr,
			traceConversionConfig: tcc,
			testRepository:        tr,
		},
		Out: OpenAPI{
			traceConversionConfig: tcc,
		},
	}
}
