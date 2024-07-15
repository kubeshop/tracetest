package trigger

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var (
	node = "node"
	app  = "npx"
	// libName = "../../tracetest-js/packages/tracetest-playwright-engine"
	libName        = "@tracetest/playwright-engine"
	baseScriptPath = "script.js"
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

	scriptPath := fmt.Sprintf("%s-%s", opts.TraceID, baseScriptPath)

	err = os.WriteFile(scriptPath, []byte(triggerConfig.PlaywrightEngine.Script), 0644)
	if err != nil {
		return response, err
	}

	out, err := start(ctx, opts.TraceID.String(), opts.SpanID.String(), triggerConfig.PlaywrightEngine.Target, triggerConfig.PlaywrightEngine.Method, scriptPath)
	os.Remove(scriptPath)
	if err != nil {
		return response, err
	}

	response.Result.PlaywrightEngine.Success = true
	response.Result.PlaywrightEngine.Out = out
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

func start(ctx context.Context, traceId, spanId, url, method, scriptPath string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if os.Getenv("PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD") != "1" {
		res, err := execCommand(
			ctx,
			app,
			"playwright",
			"install",
		)

		if err != nil {
			return "", fmt.Errorf("error installing playwright: %s, %w", res, err)
		}
	}

	path, err := filepath.Abs(fmt.Sprintf("%s/%s", wd, scriptPath))
	if err != nil {
		return "", err
	}

	res, err := execCommand(
		ctx,
		app,
		"--yes",
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
		return "", fmt.Errorf("error executing playwright engine: %s, %w", res, err)
	}

	return res, nil
}

func execCommand(ctx context.Context, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*3) //3 minutes
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		res := fmt.Sprint(err) + ": " + out.String() + stderr.String()
		return res, err
	}

	return out.String(), nil
}

type PlaywrightEngineRequest struct {
	Target string `json:"target,omitempty"  expr_enabled:"true"`
	Script string `json:"script,omitempty"`
	Method string `json:"method,omitempty"  expr_enabled:"true"`
}

type PlaywrightEngineResponse struct {
	Success bool   `json:"success"`
	Out     string `json:"out"`
}
