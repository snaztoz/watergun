package watergun

import (
	"encoding/json"
	"net/http"
)

func RespondWithError(
	w http.ResponseWriter,
	err error,
	message string,
	statusCode int,
) {
	logger.Error(message, "err", err)

	if err := json.NewEncoder(w).Encode(map[string]any{
		"err": err.Error(),
	}); err != nil {
		panic("Failed to write error JSON response")
	}
}
