package replacer

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/model"
)

var defaultInjector *Injector = nil

func ReplaceTestPlaceholders(test model.Test) (model.Test, error) {
	if defaultInjector == nil {
		injector := NewInjector()
		defaultInjector = &injector
	}

	err := defaultInjector.Inject(&test)
	if err != nil {
		return model.Test{}, fmt.Errorf("could not replace expressions in test: %w", err)
	}

	return test, nil
}
