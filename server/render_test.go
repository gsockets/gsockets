package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderJsonForSuccess(t *testing.T) {
	message := "hello world"
	statusCode := 200

	expectedJson := Response{
		Message: message,
	}

	jsonParsed, err := json.Marshal(expectedJson)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	RenderJSON(w, statusCode, message, nil)

	res := w.Result()
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)

	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"), "renderJson must return json content type")
	assert.Equal(t, res.StatusCode, statusCode, "response status code must match the code passed in the argument")
	assert.Equal(t, string(jsonParsed), string(resBody), "the response json must match the expected json")
}

func TestRenderJsonForSuccessWithData(t *testing.T) {
	data := struct {
		Data string `json:"data"`
	}{Data: "Hello world"}
	statusCode := 200

	expectedJson := Response{Data: data}

	jsonParsed, err := json.Marshal(expectedJson)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	RenderJSON(w, statusCode, "", data)

	res := w.Result()
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)

	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"), "renderJson must return json content type")
	assert.Equal(t, res.StatusCode, statusCode, "response status code must match the code passed in the argument")
	assert.Equal(t, string(jsonParsed), string(resBody), "the response json must match the expected json")
}

func TestRenderJsonForError(t *testing.T) {
	errorMessage := "invalid request"
	statusCode := 400

	expectedJson := ErrorResponse{Error: errorMessage, Code: statusCode}

	jsonParsed, err := json.Marshal(expectedJson)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	RenderJSON(w, statusCode, errorMessage, nil)

	res := w.Result()
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)

	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"), "renderJson must return json content type")
	assert.Equal(t, res.StatusCode, statusCode, "response status code must match the code passed in the argument")
	assert.Equal(t, string(jsonParsed), string(resBody), "the response json must match the expected json")
}
