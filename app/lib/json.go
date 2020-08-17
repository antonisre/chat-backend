package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/**
@desc Success response JSON
*/
func Success(response http.ResponseWriter, statusCode int, message string, data interface{}) {
	response.WriteHeader(statusCode)
	responseMap := map[string]interface{}{
		"Status":  "Success",
		"Message": message,
		"Data":    data,
	}
	err := json.NewEncoder(response).Encode(responseMap)

	if err != nil {
		fmt.Fprintf(response, "%s", err.Error())
	}
}

/**
@desc Error response JSON
*/
func Error(response http.ResponseWriter, statusCode int, message string) {
	response.WriteHeader(statusCode)
	if message == "" {
		message = "Something went wrong"
	}
	responseMap := map[string]interface{}{
		"Status":  "Error",
		"Message": message,
	}

	err := json.NewEncoder(response).Encode(responseMap)
	if err != nil {
		fmt.Fprintf(response, "%s", err.Error())
	}
}
