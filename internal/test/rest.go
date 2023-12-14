package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func VerifyRestResponseWithTestFileAsObject[T any](t *testing.T, resp *httptest.ResponseRecorder, expectedStatusCode int, filepath string) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
	require.NoError(t, err)
	expectedResponseJson, err := io.ReadAll(file)

	assert.Equal(t, expectedStatusCode, resp.Code)
	VerifyRestBodyAsObject[T](t, resp, string(expectedResponseJson))
}

func VerifyRestBodyAsObject[T any](t *testing.T, resp *httptest.ResponseRecorder, expectedBodyJson string) {
	var expected T
	err := json.NewDecoder(bytes.NewBuffer([]byte(expectedBodyJson))).Decode(&expected)
	require.NoError(t, err)

	responseBodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var actual T
	err = json.NewDecoder(bytes.NewBuffer(responseBodyBytes)).Decode(&actual)

	assert.Equal(t, expected, actual)
}

func VerifyRestResponseWithTestFile(t *testing.T, resp *httptest.ResponseRecorder, expectedStatusCode int, filepath string) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
	require.NoError(t, err)
	expectedResponseJson, err := io.ReadAll(file)

	assert.Equal(t, expectedStatusCode, resp.Code)
	VerifyRestBody(t, resp, string(expectedResponseJson))
}

func VerifyRestResponse(t *testing.T, resp *httptest.ResponseRecorder, expectedStatusCode int, expectedBodyJson string) {
	assert.Equal(t, expectedStatusCode, resp.Code)
	VerifyRestBody(t, resp, expectedBodyJson)
}

func VerifyRestBody(t *testing.T, resp *httptest.ResponseRecorder, expectedBodyJson string) {
	responseBodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.JSONEq(t, expectedBodyJson, string(responseBodyBytes))
}
