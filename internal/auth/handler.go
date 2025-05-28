package auth

import (
	"encoding/json"
	"net/http"
)

type Handle struct {
	authorization *AuthorizationServiceImpl
}

func NewHandle(authorization *AuthorizationServiceImpl) *Handle {
	return &Handle{
		authorization: authorization,
	}
}

func (h *Handle) Authorization(w http.ResponseWriter, r *http.Request) {
	var request AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
	}

	response, err := h.authorization.Login(r.Context(), request.UserName, request.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
