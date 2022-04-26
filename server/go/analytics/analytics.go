package analytics

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/denisbrodbeck/machineid"
)

var (
	appName          = "tracetest"
	analyticsEnabled = false
	gaURL            = "https://www.google-analytics.com/mp/collect?measurement_id=%s&api_secret=%s"
	gaValidationURL  = "https://www.google-analytics.com/debug/mp/collect?measurement_id=%s&api_secret=%s"
	gaMeasurementId  = ""
	gaSecretKey      = ""
)

func init() {
	gaMeasurementId = os.Getenv("GOOGLE_ANALYTICS_MEASUREMENT_ID")
	gaSecretKey = os.Getenv("GOOGLE_ANALYTICS_SECRET_KEY")
	analyticsIsEnabled, err := strconv.ParseBool(os.Getenv("ANALYTICS_ENABLED"))
	if err != nil {
		analyticsEnabled = false
	} else {
		analyticsEnabled = analyticsIsEnabled
	}
}

type Params struct {
	EventCount       int64  `json:"event_count,omitempty"`
	EventCategory    string `json:"event_category,omitempty"`
	AppVersion       string `json:"app_version,omitempty"`
	AppName          string `json:"app_name,omitempty"`
	CustomDimensions string `json:"custom_dimensions,omitempty"`
	DataSource       string `json:"data_source,omitempty"`
	Host             string `json:"host,omitempty"`
	MachineID        string `json:"machine_id,omitempty"`
	OperatingSystem  string `json:"operating_system,omitempty"`
	Architecture     string `json:"architecture,omitempty"`
}

type Event struct {
	Name   string `json:"name"`
	Params Params `json:"params,omitempty"`
}

type Payload struct {
	UserID   string  `json:"user_id,omitempty"`
	ClientID string  `json:"client_id,omitempty"`
	Events   []Event `json:"events,omitempty"`
}

type validationResponse struct {
	ValidationMessages []validationMessage `json:"validationMessages"`
}

type validationMessage struct {
	FieldPath      string `json:"fieldPath"`
	Description    string `json:"description"`
	ValidationCode string `json:"validationCode"`
}

func NewEvent(name string, category string) (Event, error) {
	appVersion := os.Getenv("VERSION")
	host, err := os.Hostname()
	if err != nil {
		return Event{}, fmt.Errorf("could not get hostname: %w", err)
	}

	machineID, err := machineid.ProtectedID(appName)
	if err != nil {
		return Event{}, fmt.Errorf("could not get machineID: %w", err)
	}

	return Event{
		Name: name,
		Params: Params{
			EventCount:      1,
			AppName:         "tracetest",
			EventCategory:   category,
			Host:            host,
			MachineID:       machineID,
			AppVersion:      appVersion,
			Architecture:    runtime.GOARCH,
			OperatingSystem: runtime.GOOS,
		},
	}, nil
}

// SendEvent sends an event to Google Analytics.
func SendEvent(event Event) error {
	if !analyticsEnabled {
		return nil
	}

	return sendEvent(event)
}

// CreateAndSendEvent is a syntax-sugar to create and send the event in a single command
func CreateAndSendEvent(name string, category string) error {
	event, err := NewEvent(name, category)
	if err != nil {
		return fmt.Errorf("could not create event: %w", err)
	}

	return SendEvent(event)
}

func sendEvent(event Event) error {
	machineID, err := machineid.ProtectedID(appName)
	if err != nil {
		return fmt.Errorf("could not get machine id: %w", err)
	}
	payload := Payload{
		UserID:   machineID,
		ClientID: machineID,
		Events: []Event{
			event,
		},
	}

	err = sendValidationRequest(payload)
	if err != nil {
		return err
	}

	err = sendDataToGA(payload)
	if err != nil {
		return fmt.Errorf("could not send request to google analytics: %w", err)
	}

	return nil
}

func sendValidationRequest(payload Payload) error {
	response, body, err := sendPayloadToURL(payload, gaValidationURL)

	if err != nil {
		return err
	}

	if response.StatusCode >= 300 {
		return fmt.Errorf("validation response got unexpected status. Got: %d", response.StatusCode)
	}

	validationResponse := validationResponse{}
	err = json.Unmarshal(body, &validationResponse)
	if err != nil {
		return fmt.Errorf("could not unmarshal response body: %w", err)
	}

	if len(validationResponse.ValidationMessages) > 0 {
		fmt.Println(validationResponse)
		return fmt.Errorf("google analytics request validation failed")
	}

	return nil
}

func sendDataToGA(payload Payload) error {
	response, _, err := sendPayloadToURL(payload, gaURL)
	if err != nil {
		return fmt.Errorf("could not send event to google analytics: %w", err)
	}

	if response.StatusCode >= 300 {
		return fmt.Errorf("google analytics returned unexpected status. Got: %d", response.StatusCode)
	}

	return nil
}

func sendPayloadToURL(payload Payload, url string) (*http.Response, []byte, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, []byte{}, fmt.Errorf("could not marshal json payload: %w", err)
	}

	request, err := http.NewRequest("POST", fmt.Sprintf(url, gaMeasurementId, gaSecretKey), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, []byte{}, fmt.Errorf("could not create request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	client := getClient()

	resp, err := client.Do(request)
	if err != nil {
		return nil, []byte{}, fmt.Errorf("could not execute request: %w", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, []byte{}, fmt.Errorf("could not read response body: %w", err)
	}

	return resp, body, err
}

func getClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				// TODO: point to valid certificate
				InsecureSkipVerify: true,
			},
		},
		Timeout: 10 * time.Second,
	}
}
