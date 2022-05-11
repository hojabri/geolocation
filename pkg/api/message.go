package api

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func SendHTTPError(w http.ResponseWriter, errMessage string, statusCode int) error {
	e := Error{
		Message: errMessage,
		Code:    statusCode,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		return err
	}
	return nil
}

func SendHTTPMessage(writer http.ResponseWriter, object interface{}, statusCode int) error {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)

	if err := json.NewEncoder(writer).Encode(object); err != nil {
		return err
	}
	return nil
}
