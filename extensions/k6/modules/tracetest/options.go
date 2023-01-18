package tracetest

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/kubeshop/tracetest/extensions/k6/models"
	"github.com/kubeshop/tracetest/extensions/k6/utils"
	"go.k6.io/k6/js/modules"
)

var defaultPropagatorList = []models.PropagatorName{
	models.PropagatorW3C,
}

type Options struct {
	ServerUrl  string
	ServerPath string
}

const (
	DefaultServerUrl = "http://localhost:3000"
	ServerURL        = "serverUrl"
	ServerPath       = "serverPath"
)

func getOptions(vu modules.VU, val goja.Value) (Options, error) {
	rawOptions := utils.ParseOptions(vu, val)
	options := Options{
		ServerUrl:  DefaultServerUrl,
		ServerPath: "",
	}

	if len(rawOptions) == 0 {
		return options, nil
	}

	for key, value := range rawOptions {
		switch key {
		case ServerURL:
			options.ServerUrl = value.ToString().String()
		case ServerPath:
			options.ServerPath = value.ToString().String()
		default:
			return options, fmt.Errorf("unknown Tracetest option '%s'", key)
		}
	}

	return options, nil
}
