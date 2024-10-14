package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type agentResponse struct {
	Key string `json:"key"`
}

func GetAgentToken(ctx context.Context, logger *zap.Logger, endpoint, orgID, envID, token string) (string, error) {
	reqUrl := fmt.Sprintf("%s/organizations/%s/environments/%s/agent", endpoint, orgID, envID)
	logger.Debug("Getting agent token", zap.String("endpoint", reqUrl))
	req, err := http.NewRequestWithContext(ctx, "GET", reqUrl, nil)
	if err != nil {
		logger.Debug("Failed to create request", zap.Error(err))
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	headers := req.Header
	headers.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Debug("Failed to get agent token", zap.Error(err))
		return "", fmt.Errorf("failed to get agent token: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		logger.Debug("Failed to get agent token", zap.String("status", res.Status))
		return "", fmt.Errorf("failed to get agent token: %s", res.Status)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Debug("Failed to read response body", zap.Error(err))
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var agentResponse agentResponse
	err = json.Unmarshal(body, &agentResponse)
	if err != nil {
		logger.Debug("Failed to unmarshal response body", zap.Error(err))
		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	logger.Debug("Agent token received")

	return agentResponse.Key, nil
}
