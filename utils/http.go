package utils

import (
	"net/http"
	"encoding/json"
	"log"
	"fmt"
)

/*
  Follows the Google JSON style guide
*/

type httpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type httpResponse struct {
	NextLink string      `json:"nextLink,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	Error    *httpError  `json:"error,omitempty"`
}

func SendSuccess(w http.ResponseWriter, data interface{}, status int) {
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

func SendPage(w http.ResponseWriter, r *http.Request, data interface{}, offset int, count int, more bool) {
	nextLink := ""
	if more {
		query := r.URL.RawQuery
		if query != "" {
			query += "&"
		} else {
			query += "?"
		}
		query += fmt.Sprintf("offset=%v&count=%v", offset, count)
		nextLink = fmt.Sprintf("https://%v%v%v", r.Host, r.RequestURI, query)
		log.Printf(nextLink)
	}

	resp := httpResponse{NextLink: nextLink, Data: data}

	w.Header().Set("Content-Type", "application/json")
	respBody, err := json.Marshal(resp)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(respBody)
	} else {
		log.Printf("error marshalling response body\n %v\n error\n %v", resp, err)
		SendError(w, "Could not convert response to JSON", http.StatusInternalServerError)
	}
}

func SendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")

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
