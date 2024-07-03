package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"

	"github.com/kubeshop/tracetest/cli/ui"
	"go.uber.org/zap"
)

type OnAuthSuccess func(token string, jwt string)
type OnAuthFailure func(err error)

type OAuthServer struct {
	endpoint         string
	frontendEndpoint string
	onSuccess        OnAuthSuccess
	onFailure        OnAuthFailure
	port             int
	server           *http.Server
	mutex            sync.Mutex
	ui               ui.UI
}

type Option func(*OAuthServer)

func NewOAuthServer(endpoint, frontendEndpoint string) *OAuthServer {
	return &OAuthServer{
		endpoint:         endpoint,
		frontendEndpoint: frontendEndpoint,
		ui:               ui.DefaultUI,
	}
}

func (s *OAuthServer) WithOnSuccess(onSuccess OnAuthSuccess) *OAuthServer {
	s.onSuccess = onSuccess
	return s
}

func (s *OAuthServer) WithOnFailure(onFailure OnAuthFailure) *OAuthServer {
	s.onFailure = onFailure
	return s
}

func (s *OAuthServer) GetAuthJWT() error {
	confirmed := s.ui.Enter("Lets get to it! Press enter to launch a browser and authenticate:")
	if !confirmed {
		s.ui.Finish()
		return nil
	}

	url, err := s.getUrl()
	if err != nil {
		return fmt.Errorf("failed to start oauth server: %w", err)
	}

	loginUrl := fmt.Sprintf("%soauth?callback=%s", s.frontendEndpoint, url)

	ui := ui.DefaultUI
	err = ui.OpenBrowser(loginUrl)
	if err != nil {
		return fmt.Errorf("failed to open the oauth url: %s", loginUrl)
	}

	return s.start()
}

type JWTResponse struct {
	Jwt string `json:"jwt"`
}

var logger = zap.NewNop()

func SetLogger(l *zap.Logger) {
	logger = l
}

func ExchangeToken(endpoint string, token string) (string, error) {
	logger.Debug("Exchanging token", zap.String("endpoint", endpoint), zap.String("token", token))
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tokens/%s/exchange", endpoint, token), nil)
	if err != nil {
		logger.Debug("Failed to create request", zap.Error(err))
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Debug("Failed to exchange token", zap.Error(err))
		return "", fmt.Errorf("failed to exchange token: %w", err)
	}

	if res.StatusCode != http.StatusCreated {
		logger.Debug("Failed to exchange token", zap.String("status", res.Status))
		return "", fmt.Errorf("failed to exchange token: %s", res.Status)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Debug("Failed to read response body", zap.Error(err))
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var jwtResponse JWTResponse
	err = json.Unmarshal(body, &jwtResponse)
	if err != nil {
		logger.Debug("Failed to unmarshal response body", zap.Error(err))
		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	logger.Debug("Token exchanged")

	return jwtResponse.Jwt, nil
}

func (s *OAuthServer) getUrl() (string, error) {
	port, err := getFreePort()
	if err != nil {
		return "", fmt.Errorf("failed to get free port: %w", err)
	}

	s.port = port
	return fmt.Sprintf("http://localhost:%d", port), nil
}

func (s *OAuthServer) start() error {
	srv := &http.Server{Addr: fmt.Sprintf(":%d", s.port)}
	s.server = srv

	http.HandleFunc("/", s.callback)
	return srv.ListenAndServe()
}

func (s *OAuthServer) callback(w http.ResponseWriter, r *http.Request) {
	tokenId, jwt, err := s.handleResult(r)
	if err != nil {
		redirect(w, r, false)
		go s.onFailure(err)
		return
	}

	redirect(w, r, true)
	go s.onSuccess(tokenId, jwt)
}

func (s *OAuthServer) handleResult(r *http.Request) (string, string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	tokenId := r.URL.Query().Get("tokenId")
	if tokenId == "" {
		return "", "", fmt.Errorf("tokenId not found")
	}

	jwt, err := ExchangeToken(s.endpoint, tokenId)
	if err != nil {
		return "", "", err
	}

	return tokenId, jwt, nil
}

func redirect(w http.ResponseWriter, r *http.Request, success bool) {
	returnUrl := r.URL.Query().Get("returnUrl")
	if returnUrl != "" {
		returnUrl = fmt.Sprintf("%s?success=%t", returnUrl, success)
	}

	http.Redirect(w, r, returnUrl, http.StatusMovedPermanently)
}

func getFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}
