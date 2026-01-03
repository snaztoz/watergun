package server

import (
	"encoding/json"
	"net/http"

	"github.com/snaztoz/watergun"
)

func respondWithJSON(w http.ResponseWriter, body any) {
	if err := json.NewEncoder(w).Encode(body); err != nil {
		panic(err)
	}
}

func respondWithError(
	w http.ResponseWriter,
	err error,
	message string,
	statusCode int,
) {
	watergun.Logger().Error(message, "err", err)

	w.WriteHeader(statusCode)

	respondWithJSON(w, map[string]any{"err": err.Error()})
}
