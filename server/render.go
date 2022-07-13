package server

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func RenderJSON(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	var body interface{}
	if statusCode >= 400 {
		body = ErrorResponse{
			Error: message,
			Code: statusCode,
		}
	} else if data != nil {
		body = Response{Data: data}
	} else {
		body = Response{Message: message}
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)
	resp, _ := json.Marshal(body)
	_, _ = w.Write(resp)
}
