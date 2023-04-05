package yamlconvert

import (
	"github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/yaml"
)

func DataStore(in model.DataStore) yaml.File {
	out := yaml.DataStore{}
	deepcopy.DeepCopy(in, &out)
	if in.Values.Jaeger != nil {
		deepcopy.DeepCopy(in.Values.Jaeger, &out.Jaeger)
		deepcopy.DeepCopy(in.Values.Jaeger.TLSSetting, &out.Jaeger.Tls)
	}
	if in.Values.Tempo != nil {
		deepcopy.DeepCopy(in.Values.Tempo, &out.Tempo)
		deepcopy.DeepCopy(in.Values.Tempo.Grpc.TLSSetting, &out.Tempo.Grpc.Tls)
		deepcopy.DeepCopy(in.Values.Tempo.Http.TLSSetting, &out.Tempo.Http.Tls)
	}

	if in.Values.OpenSearch != nil {
		deepcopy.DeepCopy(in.Values.OpenSearch, &out.OpenSearch)
	}
	if in.Values.SignalFx != nil {
		deepcopy.DeepCopy(in.Values.SignalFx, &out.SignalFx)
	}
	if in.Values.AwsXRay != nil {
		deepcopy.DeepCopy(in.Values.AwsXRay, &out.AwsXRay)
	}

	return yaml.File{
		Type: yaml.FileTypeDataStore,
		Spec: out,
	}
}
