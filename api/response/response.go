package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

func ERROR(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JSON(w, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	JSON(w, http.StatusBadRequest, nil)
}

func ToJSON(data interface{}) (result string, err error) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err == nil {
		result = string(b)
	}

	return
}

func ToJSONQuite(data interface{}) (result string) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err == nil {
		result = string(b)
	} else {
		result = ""
	}

	return
}
