package server

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func RenderSuccess(w http.ResponseWriter, message string, statusCode int) {
	resp := Response{
		Message: message,
	}

	renderJSON(w, statusCode, resp)
}

func RenderSuccessWithData(w http.ResponseWriter, message string, statusCode int, data interface{}) {
	resp := Response{
		Message: message,
		Data:    data,
	}

	renderJSON(w, statusCode, resp)
}

func RenderError(w http.ResponseWriter, message string, statusCode int) {
	resp := ErrorResponse{
		Error: message,
		Code:  statusCode,
	}

	renderJSON(w, statusCode, resp)
}

func renderJSON(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)
	resp, _ := json.Marshal(v)
	_, _ = w.Write(resp)
}
