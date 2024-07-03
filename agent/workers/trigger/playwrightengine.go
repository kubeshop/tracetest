package trigger

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	node       = "node"
	app        = "npx"
	libName    = "@tracetest/playwright-engine"
	scriptPath = "script.js"
)

func PLAYWRIGHTENGINE() Triggerer {
	return &playwrightTriggerer{}
}

type playwrightTriggerer struct{}

func (te *playwrightTriggerer) Trigger(ctx context.Context, triggerConfig Trigger, opts *Options) (Response, error) {
	response := Response{
		Result: TriggerResult{
			Type: te.Type(),
			PlaywrightEngine: &PlaywrightEngineResponse{
				Success: false,
			},
		},
	}

	err := validate()
	if err != nil {
		return response, err
	}

	err = os.WriteFile(scriptPath, []byte(triggerConfig.PlaywrightEngine.Script), 0644)
	if err != nil {
		return response, err
	}

	err = start(opts.TraceID.String(), opts.SpanID.String(), triggerConfig.PlaywrightEngine.Target, triggerConfig.PlaywrightEngine.Method)
	if err != nil {
		os.Remove(scriptPath)
		return response, err
	}

	os.Remove(scriptPath)
	response.Result.PlaywrightEngine.Success = true
	return response, err
}

func (t *playwrightTriggerer) Type() TriggerType {
	return TriggerTypePlaywrightEngine
}

const TriggerTypePlaywrightEngine TriggerType = "playwrightengine"

func validate() error {
	_, err := exec.LookPath(node)
	if err != nil {
		return fmt.Errorf("node not found in PATH")
	}

	_, err = exec.LookPath(app)
	if err != nil {
		return fmt.Errorf("npm not found in PATH")
	}

	return nil
}

func start(traceId, spanId, url, method string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	path, err := filepath.Abs(fmt.Sprintf("%s/%s", wd, scriptPath))
	if err != nil {
		return err
	}

	res, err := execCommand(
		app,
		libName,
		"--scriptPath",
		path,
		"--traceId",
		traceId,
		"--spanId",
		spanId,
		"--url",
		url,
		"--method",
		method)

	if err != nil {
		return fmt.Errorf("error executing playwright engine: %s, %w", res, err)
	}

	return nil
}

func execCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		return fmt.Sprint(err) + ": " + stderr.String(), err
	}

	return out.String(), nil
}

type PlaywrightEngineRequest struct {
	Target string `json:"target,omitempty"`
	Script string `json:"script,omitempty"`
	Method string `json:"method,omitempty"`
}

type PlaywrightEngineResponse struct {
	Success bool `json:"success"`
}
