package server

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func RenderJSON(w http.ResponseWriter, statusCode int, message string, data any) {
	var body interface{}
	if statusCode >= 400 {
		body = ErrorResponse{
			Error: message,
			Code:  statusCode,
		}
	} else if data != nil {
		body = data
	} else {
		body = Response{Message: message}
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)
	resp, _ := json.Marshal(body)
	_, _ = w.Write(resp)
}
