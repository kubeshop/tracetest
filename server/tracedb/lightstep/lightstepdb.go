package lightstep

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/traces"
)

type LightstepDB struct {
	Config config.LightstepConfig
}

func (db *LightstepDB) Close() error {
	// As it doesn't keep an open connection, we don't need to close anything here
	return nil
}

func (db *LightstepDB) GetTraceByIdentification(ctx context.Context, identification traces.TraceIdentification) (traces.Trace, error) {
	url := fmt.Sprintf(
		"https://api.lightstep.com/public/v0.2/%s/projects/%s/stored-traces?span-id=%s",
		db.Config.Organization,
		db.Config.Project,
		identification.RootSpanID.String(),
	)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return traces.Trace{}, fmt.Errorf("could not create request: %s", err)
	}

	request.Header.Add("Authorization", db.Config.Token)

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
