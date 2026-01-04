package response

import (
	"encoding/json"
	"net/http"

	"github.com/snaztoz/watergun/log"
)

func SendJSON(w http.ResponseWriter, body any) {
	if err := json.NewEncoder(w).Encode(body); err != nil {
		panic(err)
	}
}

func SendErrorJSON(w http.ResponseWriter, err error, message string, statusCode int) {
	log.Error(message, "err", err)

	w.WriteHeader(statusCode)

	SendJSON(w, map[string]any{"err": err.Error()})
}
