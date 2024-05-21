package runner

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

func (s *Runner) getAgentToken(ctx context.Context, endpoint, orgID, envID, token string) (string, error) {
	s.logger.Debug("Getting agent token", zap.String("endpoint", endpoint), zap.String("token", token))
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/organizations/%s/environments/%s/agent", endpoint, orgID, envID), nil)
	if err != nil {
		s.logger.Debug("Failed to create request", zap.Error(err))
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	headers := req.Header
	headers.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Debug("Failed to get agent token", zap.Error(err))
		return "", fmt.Errorf("failed to get agent token: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		s.logger.Debug("Failed to get agent token", zap.String("status", res.Status))
		return "", fmt.Errorf("failed to get agent token: %s", res.Status)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Debug("Failed to read response body", zap.Error(err))
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var agentResponse agentResponse
	err = json.Unmarshal(body, &agentResponse)
	if err != nil {
		s.logger.Debug("Failed to unmarshal response body", zap.Error(err))
		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	s.logger.Debug("Agent token received")

	return agentResponse.Key, nil
}
