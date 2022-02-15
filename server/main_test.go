package main_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	openapi "github.com/GIT_USER_ID/GIT_REPO_ID/go"
	"github.com/GIT_USER_ID/GIT_REPO_ID/go/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewTest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	test := openapi.Test{
		Name: "test",
		ServiceUnderTest: openapi.TestServiceUnderTest{
			Url: "http://localhost:3030/test",
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
	ex := mocks.NewMockTestExecutor(ctrl)

	ex.EXPECT().Execute(gomock.Any()).Return(&openapi.Result{Id: "2"}, nil)
	req := httptest.NewRequest(http.MethodPost, "/tests/1/run", nil)
	w := httptest.NewRecorder()

	ApiApiService := openapi.NewApiApiService(nil, db, ex)
	controller := openapi.NewApiApiController(ApiApiService)

	ctr := controller.(*openapi.ApiApiController)
	ctr.TestsTestidRunPost(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, w.Code, http.StatusOK)

	t.Logf("response: %s\n", string(data))
}
