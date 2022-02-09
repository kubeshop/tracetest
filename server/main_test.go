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

	ApiApiService := openapi.NewApiApiService(nil, db)
	ApiApiController := openapi.NewApiApiController(ApiApiService)

	ApiApiController.CreateTest(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, w.Code, http.StatusOK)

	t.Logf("response: %s\n", string(data))
}
