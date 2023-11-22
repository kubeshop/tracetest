package runner

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	agentConfig "github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/agent/event"
	"github.com/kubeshop/tracetest/agent/proto"
	consoleUI "github.com/kubeshop/tracetest/agent/ui"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

func (s *Runner) RunVerboseStrategy(ctx context.Context, cfg agentConfig.Config) error {
	s.ui.Infof("%s Starting Agent with name %s...", consoleUI.Emoji_Truck, cfg.Name)

	session, err := StartSession(ctx, cfg, &verboseObserver{ui: s.ui}, s.logger)

	if err != nil && errors.Is(err, ErrOtlpServerStart) {
		s.ui.Errorf("%s Tracetest Agent binds to the OpenTelemetry ports 4317 and 4318 which are used to receive trace information from your system. The agent tried to bind to these ports, but failed.", consoleUI.Emoji_RedCircle)

		s.ui.Finish()
		return err
	}

	s.ui.Infof("%s Agent is started!", consoleUI.Emoji_Truck)
	s.ui.Println("")

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-signalChannel
	s.ui.Infof("%s Received stop signal %s. Stopping agent...", consoleUI.Emoji_YellowCircle, sig)

	session.Close()
	s.ui.Finish()

	return nil
}

type verboseObserver struct {
	ui consoleUI.ConsoleUI
}

var _ event.Observer = &verboseObserver{}

func (o *verboseObserver) EndDataStoreConnection(request *proto.DataStoreConnectionTestRequest, err error) {
	if err != nil {
		o.ui.Warningf("%s Testing connection to DataStore %s failed!", consoleUI.Emoji_YellowCircle, request.Datastore.Type)
		o.ui.Warningf("Error: %s", err.Error())
		o.ui.Println("")
		return
	}

	o.ui.Infof("%s Testing connection to %s data store succeed", consoleUI.Emoji_WhiteCheckMark, request.Datastore.Type)
	o.ui.Println("")
}

func (o *verboseObserver) EndSpanReceive(spans []*v1.Span, err error) {
	if err != nil {
		o.ui.Warningf("%s Trace spans reception failed!", consoleUI.Emoji_YellowCircle)
		o.ui.Warningf("Error: %s", err.Error())
		o.ui.Println("")
		return
	}

	o.ui.Infof("%s Trace spans received with success. %d spans received", consoleUI.Emoji_WhiteCheckMark, len(spans))
	o.ui.Println("")
}

func (o *verboseObserver) EndTracePoll(request *proto.PollingRequest, err error) {
	if err != nil {
		o.ui.Warningf("%s Test run %d, test %s trace spans fetch failed!", consoleUI.Emoji_YellowCircle, request.RunID, request.TestID)
		o.ui.Warningf("Error: %s", err.Error())
		o.ui.Println("")
		return
	}

	o.ui.Infof("%s Test run %d, test %s trace spans fetch with success", consoleUI.Emoji_WhiteCheckMark, request.RunID, request.TestID)
	o.ui.Println("")
}

func (o *verboseObserver) EndTriggerExecution(request *proto.TriggerRequest, err error) {
	if err != nil {
		o.ui.Warningf("%s Test run %d, test %s trigger failed!", consoleUI.Emoji_YellowCircle, request.RunID, request.TestID)
		o.ui.Warningf("Error: %s", err.Error())
		o.ui.Println("")
		return
	}

	o.ui.Infof("%s Test run %d, test %s trigger executed with success", consoleUI.Emoji_WhiteCheckMark, request.RunID, request.TestID)
	o.ui.Println("")
}

func (o *verboseObserver) Error(err error) {
	o.ui.Errorf("%s An unknown error happened on Tracetest agent.", consoleUI.Emoji_RedCircle)
	o.ui.Errorf("Error: %s", err.Error())
	o.ui.Println("")
}

func (o *verboseObserver) StartDataStoreConnection(request *proto.DataStoreConnectionTestRequest) {
	o.ui.Infof("%s Testing connection to %s data store ...", consoleUI.Emoji_Magnifier, request.Datastore.Type)
}

func (o *verboseObserver) StartSpanReceive(spans []*v1.Span) {
	o.ui.Infof("%s Receiving trace spans...", consoleUI.Emoji_Sparkles)
}

func (o *verboseObserver) StartTracePoll(request *proto.PollingRequest) {
	o.ui.Infof("%s Polling traces and spans for test run %d, test %s ...", consoleUI.Emoji_Magnifier, request.RunID, request.TestID)
}

func (o *verboseObserver) StartTriggerExecution(request *proto.TriggerRequest) {
	o.ui.Infof("%s Executing trigger for test run %d, test %s ...", consoleUI.Emoji_Sparkles, request.RunID, request.TestID)
}
