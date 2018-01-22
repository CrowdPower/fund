package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

/*
  Follows the Google JSON style guide
*/

type httpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type httpResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error *httpError  `json:"error,omitempty"`
}

func SendSuccess(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")

	resp := httpResponse{Data: data}

	respBody, err := json.Marshal(resp)
	if err == nil {
		w.WriteHeader(status)
		w.Write(respBody)
	} else {
		log.Printf("error marshalling response body\n %v\n error\n %v", resp, err)
		SendError(w, "Could not convert response to JSON", http.StatusInternalServerError)
	}
}

func SendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	resp := httpResponse{Error: &httpError{Code: status, Message: message}}

	respBody, err := json.Marshal(resp)
	if err == nil {
		w.WriteHeader(status)
		w.Write(respBody)
	} else {
		log.Printf("error marshalling response body\n %v\n error\n %v", resp, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
