package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	neturl "net/url"
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

	url, listener, err := s.getUrl()
	if err != nil {
		return fmt.Errorf("failed to start oauth server: %w", err)
	}

	loginUrl, err := neturl.JoinPath(s.frontendEndpoint, "oauth")
	if err != nil {
		return fmt.Errorf("could not build path: %w", err)
	}

	loginUrl += fmt.Sprintf("?callback=%s", url)

	ui := ui.DefaultUI
	err = ui.OpenBrowser(loginUrl)
	if err != nil {
		return fmt.Errorf("failed to open the oauth url: %s", loginUrl)
	}

	return s.start(listener)
}

type JWTResponse struct {
	Jwt string `json:"jwt"`
}

var logger = zap.NewNop()

func SetLogger(l *zap.Logger) {
	logger = l
}

type oauthError struct {
	err error
	msg string
}

func (e oauthError) Error() string {
	return e.err.Error()
}

func (e oauthError) Message() string {
	return e.msg
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
		logger.Debug("Cannot create exchange token request", zap.Error(err))
		return "", fmt.Errorf("cannot create exchange token request: %w", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusNotFound:
		return "", oauthError{err: fmt.Errorf("token not found"), msg: "Token not found"}
	case http.StatusUnauthorized:
		return "", oauthError{err: fmt.Errorf("token expired"), msg: "Token has expired"}
	case http.StatusCreated:
		logger.Debug("Token exchanged")
	default:
		b, _ := io.ReadAll(res.Body)
		logger.Debug("Failed to exchange token", zap.String("status", res.Status), zap.String("response", string(b)))
		return "", oauthError{err: fmt.Errorf("failed to exchange token: %s", res.Status), msg: "Unexpected error exchanging token"}
	}

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

func (s *OAuthServer) getUrl() (string, *net.TCPListener, error) {
	port, listener, err := getFreePort()
	if err != nil {
		return "", nil, fmt.Errorf("failed to get free port: %w", err)
	}

	s.port = port
	return fmt.Sprintf("http://localhost:%d", port), listener, nil
}

func (s *OAuthServer) start(listener *net.TCPListener) error {
	srv := &http.Server{Addr: fmt.Sprintf(":%d", s.port)}
	s.server = srv

	http.HandleFunc("/", s.callback)
	return srv.Serve(listener)
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

func getFreePort() (port int, listener *net.TCPListener, err error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		err = fmt.Errorf("failed to resolve tcp addr: %w", err)
		return
	}

	listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		err = fmt.Errorf("failed to listen: %w", err)
		return
	}
	port = listener.Addr().(*net.TCPAddr).Port

	return
}
