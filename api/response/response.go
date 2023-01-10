package response

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	// "path/filepath"

	"github.com/xuri/excelize/v2"
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
func JSONDOWNLOAD(w http.ResponseWriter, statusCode int, file *excelize.File) {

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename ="+strconv.Quote(strconv.FormatInt(time.Now().Unix(), 10)+"Download.xlsx"))
	w.Header().Set("Content-Transfer-Encoding", "binary")
	file.Write(w)
	w.WriteHeader(statusCode)
}
