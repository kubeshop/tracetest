package trigger_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kubeshop/tracetest/agent/workers/trigger"
	triggerer "github.com/kubeshop/tracetest/server/executor/trigger"
	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func createContext() context.Context {
	return trace.ContextWithSpanContext(context.TODO(), trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: id.NewRandGenerator().TraceID(),
		SpanID:  id.NewRandGenerator().SpanID(),
	}))
}

func TestTriggerGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		t.Log(req.Header)

		tp, ok := req.Header["Traceparent"]
		if !ok {
			t.Fatalf("missing Traceparent header %#v", req.Header)
		}
		assert.Len(t, tp, 1)

		assert.Equal(t, "GET", req.Method)
		assert.Equal(t, "Value1", req.Header.Get("Key1"))

		b, err := io.ReadAll(req.Body)
		assert.NoError(t, err)
		assert.Equal(t, "body", string(b))

		rw.WriteHeader(200)
		_, err = rw.Write([]byte(`OK`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	test := test.Test{
		Name: "test",
		Trigger: trigger.Trigger{
			Type: trigger.TriggerTypeHTTP,
			HTTP: &trigger.HTTPRequest{
				URL:    server.URL,
				Method: trigger.HTTPMethodGET,
				Headers: []trigger.HTTPHeader{
					{Key: "Key1", Value: "Value1"},
				},
				Body: "body",
			},
		},
	}

	ex := triggerer.HTTP()

	resp, err := ex.Trigger(createContext(), test, nil)
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.Result.HTTP.StatusCode)
	assert.Equal(t, "OK", resp.Result.HTTP.Body)
}

func TestTriggerPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		t.Log(req.Header)

		tp, ok := req.Header["Traceparent"]
		if !ok {
			t.Fatalf("missing Traceparent header %#v", req.Header)
		}
		assert.Len(t, tp, 1)

		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "Value1", req.Header.Get("Key1"))

		b, err := io.ReadAll(req.Body)
		assert.NoError(t, err)
		assert.Equal(t, "body", string(b))

		rw.WriteHeader(200)
		_, err = rw.Write([]byte(`OK`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	test := test.Test{
		Name: "test",
		Trigger: trigger.Trigger{
			Type: trigger.TriggerTypeHTTP,
			HTTP: &trigger.HTTPRequest{
				URL:    server.URL,
				Method: trigger.HTTPMethodPOST,
				Headers: []trigger.HTTPHeader{
					{Key: "Key1", Value: "Value1"},
				},
				Body: "body",
			},
		},
	}

	ex := triggerer.HTTP()

	resp, err := ex.Trigger(createContext(), test, nil)
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.Result.HTTP.StatusCode)
	assert.Equal(t, "OK", resp.Result.HTTP.Body)
}

func TestTriggerPostWithApiKeyAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		t.Log(req.Header)

		tp, ok := req.Header["Traceparent"]
		if !ok {
			t.Fatalf("missing Traceparent header %#v", req.Header)
		}
		assert.Len(t, tp, 1)

		tp, ok = req.Header["Key"]
		if !ok {
			t.Fatalf("missing key header %#v", req.Header)
		}
		assert.Len(t, tp, 1)
		assert.Equal(t, tp[0], "value")
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "Value1", req.Header.Get("Key1"))

		b, err := io.ReadAll(req.Body)
		assert.NoError(t, err)
		assert.Equal(t, "body", string(b))

		rw.WriteHeader(200)
		_, err = rw.Write([]byte(`OK`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	test := test.Test{
		Name: "test",
		Trigger: trigger.Trigger{
			Type: trigger.TriggerTypeHTTP,
			HTTP: &trigger.HTTPRequest{
				URL:    server.URL,
				Method: trigger.HTTPMethodPOST,
				Headers: []trigger.HTTPHeader{
					{Key: "Key1", Value: "Value1"},
				},
				Auth: &trigger.HTTPAuthenticator{
					Type: "apiKey",
					APIKey: &trigger.APIKeyAuthenticator{
						Key:   "key",
						Value: "value",
						In:    trigger.APIKeyPositionHeader,
					},
				},
				Body: "body",
			},
		},
	}

	ex := triggerer.HTTP()

	resp, err := ex.Trigger(createContext(), test, nil)
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.Result.HTTP.StatusCode)
	assert.Equal(t, "OK", resp.Result.HTTP.Body)
}

func TestTriggerPostWithBasicAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		t.Log(req.Header)

		tp, ok := req.Header["Traceparent"]
		if !ok {
			t.Fatalf("missing Traceparent header %#v", req.Header)
		}
		assert.Len(t, tp, 1)

		tp, ok = req.Header["Authorization"]
		if !ok {
			t.Fatalf("missing Authorization header %#v", req.Header)
		}
		assert.Len(t, tp, 1)
		assert.Equal(t, tp[0], "Basic dXNlcm5hbWU6cGFzc3dvcmQ=")
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "Value1", req.Header.Get("Key1"))

		b, err := io.ReadAll(req.Body)
		assert.NoError(t, err)
		assert.Equal(t, "body", string(b))

		rw.WriteHeader(200)
		_, err = rw.Write([]byte(`OK`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	test := test.Test{
		Name: "test",
		Trigger: trigger.Trigger{
			Type: trigger.TriggerTypeHTTP,
			HTTP: &trigger.HTTPRequest{
				URL:    server.URL,
				Method: trigger.HTTPMethodPOST,
				Headers: []trigger.HTTPHeader{
					{Key: "Key1", Value: "Value1"},
				},
				Auth: &trigger.HTTPAuthenticator{
					Type: "basic",
					Basic: &trigger.BasicAuthenticator{
						Username: "username",
						Password: "password",
					},
				},
				Body: "body",
			},
		},
	}

	ex := triggerer.HTTP()

	resp, err := ex.Trigger(createContext(), test, nil)
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.Result.HTTP.StatusCode)
	assert.Equal(t, "OK", resp.Result.HTTP.Body)
}

func TestTriggerPostWithBearerAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		t.Log(req.Header)

		tp, ok := req.Header["Traceparent"]
		if !ok {
			t.Fatalf("missing Traceparent header %#v", req.Header)
		}
		assert.Len(t, tp, 1)

		tp, ok = req.Header["Authorization"]
		if !ok {
			t.Fatalf("missing Authorization header %#v", req.Header)
		}
		assert.Len(t, tp, 1)
		assert.Equal(t, tp[0], "token")
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "Value1", req.Header.Get("Key1"))

		b, err := io.ReadAll(req.Body)
		assert.NoError(t, err)
		assert.Equal(t, "body", string(b))

		rw.WriteHeader(200)
		_, err = rw.Write([]byte(`OK`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	test := test.Test{
		Name: "test",
		Trigger: trigger.Trigger{
			Type: trigger.TriggerTypeHTTP,
			HTTP: &trigger.HTTPRequest{
				URL:    server.URL,
				Method: trigger.HTTPMethodPOST,
				Headers: []trigger.HTTPHeader{
					{Key: "Key1", Value: "Value1"},
				},
				Auth: &trigger.HTTPAuthenticator{
					Type: "bearer",
					Bearer: &trigger.BearerAuthenticator{
						Bearer: "token",
					},
				},
				Body: "body",
			},
		},
	}

	ex := triggerer.HTTP()

	resp, err := ex.Trigger(createContext(), test, nil)
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.Result.HTTP.StatusCode)
	assert.Equal(t, "OK", resp.Result.HTTP.Body)
}

func TestTriggerResolveLargePayloadWithSingleQuote(t *testing.T) {
	testObject := test.Test{
		Name: "Run API to trigger Gemini ",
		Trigger: trigger.Trigger{
			Type: trigger.TriggerTypeHTTP,
			HTTP: &trigger.HTTPRequest{
				URL:    "http://llm-api:8800/summarizeText",
				Method: trigger.HTTPMethodPOST,
				Headers: []trigger.HTTPHeader{
					{Key: "Content-Type", Value: "application/json"},
				},
				Body: "{\n          \"provider\": \"Google (Gemini)\",\n          \"text\": \"Born in London, Turing was raised in southern England. He graduated from King's College, Cambridge, and in 1938, earned a doctorate degree from Princeton University. During World War II, Turing worked for the Government Code and Cypher School at Bletchley Park, Britain's codebreaking centre that produced Ultra intelligence. He led Hut 8, the section responsible for German naval cryptanalysis. Turing devised techniques for speeding the breaking of German ciphers, including improvements to the pre-war Polish bomba method, an electromechanical machine that could find settings for the Enigma machine. He played a crucial role in cracking intercepted messages that enabled the Allies to defeat the Axis powers in many crucial engagements, including the Battle of the Atlantic.\\n\\nAfter the war, Turing worked at the National Physical Laboratory, where he designed the Automatic Computing Engine, one of the first designs for a stored-program computer. In 1948, Turing joined Max Newman's Computing Machine Laboratory at the Victoria University of Manchester, where he helped develop the Manchester computers[12] and became interested in mathematical biology. Turing wrote on the chemical basis of morphogenesis and predicted oscillating chemical reactions such as the Belousov–Zhabotinsky reaction, first observed in the 1960s. Despite these accomplishments, he was never fully recognised during his lifetime because much of his work was covered by the Official Secrets Act.\"\n        }",
			},
		},
	}

	triggerOptions := &triggerer.ResolveOptions{
		Executor: expression.NewExecutor(),
	}

	httpTriggerer := triggerer.HTTP()
	resolvedTest, err := httpTriggerer.Resolve(context.Background(), testObject, triggerOptions)

	assert.NoError(t, err)
	assert.NotNil(t, resolvedTest)
}
