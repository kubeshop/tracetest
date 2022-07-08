package lightstep

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/traces"
)

type LightstepAPIConfig struct {
	Organization string
	Project      string
	Token        string
}

type lightstepDB struct {
	config LightstepAPIConfig
}

func New(config LightstepAPIConfig) tracedb.TraceDB {
	return &lightstepDB{config}
}

func (db *lightstepDB) Close() error {
	// As it doesn't keep an open connection, we don't need to close anything here
	return nil
}

func (db *lightstepDB) GetTraceByIdentification(ctx context.Context, identification traces.TraceIdentification) (traces.Trace, error) {
	url := fmt.Sprintf(
		"https://api.lightstep.com/public/v0.2/%s/projects/%s/stored-traces?span-id=%s",
		db.config.Organization,
		db.config.Project,
		identification.RootSpanID.String(),
	)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return traces.Trace{}, fmt.Errorf("could not create request: %s", err)
	}

	request.Header.Add("Authorization", db.config.Token)

	// TODO: inject a configured client here
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return traces.Trace{}, fmt.Errorf("could not execute request: %w", err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return traces.Trace{}, fmt.Errorf("could not read response body: %w", err)
	}

	var traceResponse GetTraceResponse
	err = json.Unmarshal(body, &traceResponse)
	if err != nil {
		return traces.Trace{}, fmt.Errorf("could not unmarshall json: %w", err)
	}

	return ConvertResponseToOtelFormat(traceResponse), nil
}

var _ tracedb.TraceDB = &lightstepDB{}
