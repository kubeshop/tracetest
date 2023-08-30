package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"runtime"
)

type OnAuthSuccess func(token string, jwt string)
type OnAuthFailure func(err error)

type OAuthServer struct {
	endpoint         string
	frontendEndpoint string
	onSuccess        OnAuthSuccess
	onFailure        OnAuthFailure
	port             int
}

type Option func(*OAuthServer)

func NewOAuthServer(endpoint, frontendEndpoint string) *OAuthServer {
	return &OAuthServer{
		endpoint:         endpoint,
		frontendEndpoint: frontendEndpoint,
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
	url, err := s.getUrl()
	if err != nil {
		return fmt.Errorf("failed to start oauth server: %w", err)
	}

	loginUrl := fmt.Sprintf("%s/oauth?callback=%s", s.frontendEndpoint, url)

	err = openBrowser(loginUrl)
	if err != nil {
		return fmt.Errorf("failed to open the oauth url: %s", loginUrl)
	}

	return s.start()
}

type JWTResponse struct {
	Jwt string `json:"jwt"`
}

func (s *OAuthServer) ExchangeToken(token string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tokens/%s/exchange", s.endpoint, token), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to exchange token: %w", err)
	}

	if res.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to exchange token: %s", res.Status)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var jwtResponse JWTResponse
	err = json.Unmarshal(body, &jwtResponse)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}

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
	http.HandleFunc("/", s.callback)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

func (s *OAuthServer) callback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", s.frontendEndpoint)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))

	go s.handleResult(r)
}

func (s *OAuthServer) handleResult(r *http.Request) {
	tokenId := r.URL.Query().Get("tokenId")
	if tokenId == "" {
		s.onFailure(fmt.Errorf("tokenId not found"))
		return
	}

	jwt, err := s.ExchangeToken(tokenId)
	if err != nil {
		s.onFailure(err)
		return
	}

	s.onSuccess(tokenId, jwt)
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

func openBrowser(u string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", u).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", u).Start()
	case "darwin":
		return exec.Command("open", u).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}
