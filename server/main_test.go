package main_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	openapi "github.com/kubeshop/tracetest/server/go"
	"github.com/kubeshop/tracetest/server/go/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewTest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	test := openapi.Test{
		Name: "test",
		ServiceUnderTest: openapi.TestServiceUnderTest{
			Request: openapi.HttpRequest{
				Url:    "http://localhost:3030/test",
				Method: "GET",
			},
		},
	}
	b, err := json.Marshal(test)
	assert.NoError(t, err)
	t.Logf("request: %s\n", b)

	db := mocks.NewMockTestDB(ctrl)
	db.EXPECT().CreateTest(gomock.Any(), &test).Return("id", nil)

	req := httptest.NewRequest(http.MethodPost, "/tests", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	ApiApiService := openapi.NewApiApiService(nil, db, nil)
	controller := openapi.NewApiApiController(ApiApiService)

	ctr := controller.(*openapi.ApiApiController)
	ctr.CreateTest(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, w.Code, http.StatusOK)

	t.Logf("response: %s\n", string(data))
}

func TestRunTest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mocks.NewMockTestDB(ctrl)
	db.EXPECT().GetTest(gomock.Any(), gomock.Any()).Return(&openapi.Test{}, nil)
	db.EXPECT().CreateResult(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	db.EXPECT().UpdateResult(gomock.Any(), gomock.Any()).Return(nil)
	db.EXPECT().UpdateTest(gomock.Any(), gomock.Any()).Return(nil)
	ex := mocks.NewMockTestExecutor(ctrl)

	ex.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any()).Return(&openapi.TestRunResult{ResultId: "2"}, nil)
	req := httptest.NewRequest(http.MethodPost, "/api/tests/1/run", nil)
	w := httptest.NewRecorder()

	ApiApiService := openapi.NewApiApiService(nil, db, ex)
	controller := openapi.NewApiApiController(ApiApiService)

	ctr := controller.(*openapi.ApiApiController)

	ctr.TestsTestIdRunPost(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, w.Code, http.StatusOK)

	t.Logf("response: %s\n", string(data))
	//TODO: test executed in seperate goroutine
	time.Sleep(100 * time.Millisecond)
}

func TestCreateNewAssertion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	assertion := openapi.Assertion{
		Selectors: []openapi.SelectorItem{
			{
				LocationName: "SPAN",
				PropertyName: "operation",
				Value:        "POST /users/verify",
				ValueType:    "stringValue",
			},
		},
		SpanAssertions: []openapi.SpanAssertion{
			{
				LocationName:    "SPAN_ATTRIBUTES",
				PropertyName:    "http.status.code",
				ValueType:       "intValue",
				Operator:        "EQUALS",
				ComparisonValue: "200",
			},
		},
	}
	b, err := json.Marshal(assertion)
	assert.NoError(t, err)
	t.Logf("request: %s\n", b)

	db := mocks.NewMockTestDB(ctrl)
	db.EXPECT().CreateAssertion(gomock.Any(), "", &assertion).Return("id", nil)

	req := httptest.NewRequest(http.MethodPost, "/api/tests/testid/assertions", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	ApiApiService := openapi.NewApiApiService(nil, db, nil)
	controller := openapi.NewApiApiController(ApiApiService)

	ctr := controller.(*openapi.ApiApiController)
	ctr.CreateAssertion(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, w.Code, http.StatusOK)

	t.Logf("response: %s\n", string(data))
}
