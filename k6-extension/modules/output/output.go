package output

import (
	"fmt"
	"math/rand"

	"github.com/sirupsen/logrus"
	"github.com/xoscar/xk6-tracetest-tracing/modules/tracetest"

	"go.k6.io/k6/metrics"
	"go.k6.io/k6/output"
)

type Output struct {
	config    Config
	testRunID int64
	logger    logrus.FieldLogger
	tracetest *tracetest.Tracetest
}

var _ output.Output = new(Output)

func New(params output.Params, tracetest *tracetest.Tracetest) (*Output, error) {
	config, err := NewConfig(params)
	if err != nil {
		return nil, err
	}

	return &Output{
		config:    config,
		tracetest: tracetest,
		logger:    params.Logger.WithField("component", "xk6-tracetest-output"),
	}, nil
}

func (o *Output) Description() string {
	return fmt.Sprintf("xk6-crocospans (TestRunID: %d)", o.testRunID)
}

func (o *Output) AddMetricSamples(samples []metrics.SampleContainer) {
	if len(samples) == 0 {
		return
	}

	for _, s := range samples {
		o.handleSample(s)
	}
}

func (o *Output) Stop() error {
	o.logger.Debug("Stopping...")
	defer o.logger.Debug("Stopped!")

	return nil
}

func (o *Output) Start() error {
	o.logger.Debug("Starting...")
	o.testRunID = 10000 + rand.Int63n(99999-10000)
	o.logger.Debug("Started!")

	return nil
}
