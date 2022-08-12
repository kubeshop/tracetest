package replacer

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/model"
)

var defaultInjector Injector = NewInjector()

func ReplaceTestPlaceholders(test model.Test) (model.Test, error) {
	err := defaultInjector.Inject(&test)
	if err != nil {
		return model.Test{}, fmt.Errorf("could not replace expressions in test: %w", err)
	}

	return test, nil
}
