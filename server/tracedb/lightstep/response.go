package lightstep

// Used https://mholt.github.io/json-to-go/ to generate this struct, then broke it into more structs

type GetTraceResponse struct {
	Data []traceData `json:"data"`
}

type traceData struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes struct {
		StartTimeMicros int64  `json:"start-time-micros"`
		EndTimeMicros   int64  `json:"end-time-micros"`
		Spans           []span `json:"spans"`
	} `json:"attributes"`
	Relationships struct {
		Reporters []reporter `json:"reporters"`
	} `json:"relationships"`
}

type span struct {
	SpanName        string                 `json:"span-name"`
	SpanID          string                 `json:"span-id"`
	IsError         bool                   `json:"is-error"`
	StartTimeMicros int64                  `json:"start-time-micros"`
	EndTimeMicros   int64                  `json:"end-time-micros"`
	TraceID         string                 `json:"trace-id"`
	ReporterID      string                 `json:"reporter-id"`
	Tags            map[string]interface{} `json:"tags,omitempty"`
}

type reporter struct {
	ReporterID string                 `json:"reporter-id"`
	Attributes map[string]interface{} `json:"attributes"`
}
