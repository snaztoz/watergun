package server

import (
	"encoding/json"
	"net/http"

	"github.com/snaztoz/watergun"
)

func userCreationHandler(userDomain *watergun.UserDomain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var dto UserCreationDTO
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			watergun.Logger().Error("Failed to decode request body", "err", err)
			http.Error(w, err.Error(), 400)
			return
		}

		// TODO: validate DTO

		user, err := userDomain.CreateUser(dto.MasterID)
		if err != nil {
			watergun.Logger().Error("Failed to create user", "err", err)
			http.Error(w, err.Error(), 422)
			return
		}

		if err := json.NewEncoder(w).Encode(user); err != nil {
			watergun.Logger().Error("Failed to write response", "err", err)
			http.Error(w, err.Error(), 500)
		}
	}
}

type UserCreationDTO struct {
	MasterID string `json:"master_id"`
}
