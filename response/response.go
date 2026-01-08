package response

import (
	"encoding/json"
	"net/http"
)

func SendJSON(w http.ResponseWriter, body any) {
	if err := json.NewEncoder(w).Encode(body); err != nil {
		panic(err)
	}
}

func SendErrorJSON(w http.ResponseWriter, msg string, statusCode int) {
	w.WriteHeader(statusCode)
	SendJSON(w, map[string]any{"err": msg})
}
