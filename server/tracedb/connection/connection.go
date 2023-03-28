package connection

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/kubeshop/tracetest/server/model"
)

const reachabilityTimeout = 5 * time.Second

var (
	ErrTraceNotFound        = errors.New("trace not found")
	ErrInvalidConfiguration = errors.New("invalid data store configuration")
	ErrConnectionFailed     = errors.New("could not connect to data store")
)

func CheckReachability(endpoint string, protocol model.Protocol) error {
	if protocol == model.ProtocolHTTP {
		address, err := url.Parse(endpoint)
		if err != nil {
			return err
		}

		endpoint = strings.TrimPrefix(endpoint, "http://")
		endpoint = strings.TrimPrefix(endpoint, "https://")

		if address.Scheme == "https" && address.Port() == "" {
			endpoint = fmt.Sprintf("%s:443", address.Hostname())
		}

		if address.Scheme == "http" && address.Port() == "" {
			endpoint = fmt.Sprintf("%s:80", address.Hostname())
		}
	}

	_, err := net.DialTimeout("tcp", endpoint, reachabilityTimeout)
	if err != nil {
		return err
	}

	return nil
}
