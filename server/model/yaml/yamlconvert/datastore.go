package yamlconvert

import (
	"github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/yaml"
)

func DataStore(in model.DataStore) yaml.File {
	out := yaml.DataStore{}
	deepcopy.DeepCopy(in, &out)
	deepcopy.DeepCopy(in.Values.Jaeger, &out.Jaeger)
	deepcopy.DeepCopy(in.Values.Jaeger.TLSSetting, &out.Jaeger.Tls)
	deepcopy.DeepCopy(in.Values.Tempo, &out.Tempo)
	deepcopy.DeepCopy(in.Values.Tempo.Grpc.TLSSetting, &out.Tempo.Grpc.Tls)
	deepcopy.DeepCopy(in.Values.Tempo.Http.TLSSetting, &out.Tempo.Http.Tls)
	deepcopy.DeepCopy(in.Values.OpenSearch, &out.OpenSearch)
	deepcopy.DeepCopy(in.Values.SignalFx, &out.SignalFx)

	return yaml.File{
		Type: yaml.FileTypeDataStore,
		Spec: out,
	}
}
