package service

import (
	"avito_tech/internal/auth"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Handler struct {
	store *StoreService
	coins *CoinsService
	info  *Info
}

type ErrorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"statusCode"`
	Details    string `json:"details"`
}

func NewHandler(store *StoreService, coins *CoinsService, info *Info) *Handler {

	return &Handler{store: store, coins: coins, info: info}
}

func (h *Handler) BuyItem(w http.ResponseWriter, r *http.Request) {
	itemName := chi.URLParam(r, "item")

	claims, ok := r.Context().Value("jwt_claims").(*auth.Claims)
	if !ok {
		writeError(w, "Invalid token", http.StatusUnauthorized, nil)
		return
	}

	response, err := h.store.BuyItem(r.Context(), claims.UserID, itemName)
	if err != nil {
		writeError(w, "Failed to buy item", http.StatusInternalServerError, err)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		writeError(w, "Failed to encode response", http.StatusInternalServerError, err)
	}
}

func (h *Handler) SendCoins(w http.ResponseWriter, r *http.Request) {
	var request SendCoinsRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeError(w, "Invalid request", http.StatusUnauthorized, err)
	}

	claims, ok := r.Context().Value("jwt_claims").(*auth.Claims)
	if !ok {
		writeError(w, "Invalid token", http.StatusUnauthorized, nil)
		return
	}

	response, err := h.coins.SendCoins(r.Context(), claims.UserID, request.ToUser, request.Amount)
	if err != nil {
		writeError(w, "Failed to send coins", http.StatusInternalServerError, err)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		writeError(w, "Failed to encode response", http.StatusInternalServerError, err)
	}
}

func (h *Handler) GetInfo(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("jwt_claims").(*auth.Claims)
	if !ok {
		writeError(w, "Invalid token", http.StatusUnauthorized, nil)
		return
	}

	response, err := h.info.GetUserInfo(r.Context(), claims.UserID)
	if err != nil {
		writeError(w, "Failed to get user info", http.StatusInternalServerError, err)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		writeError(w, "Failed to encode response", http.StatusInternalServerError, err)
	}
}

func writeError(w http.ResponseWriter, message string, statusCode int, details error) {
	w.WriteHeader(statusCode)

	resp := ErrorResponse{
		Error:      message,
		StatusCode: statusCode,
	}

	if details != nil {
		resp.Details = details.Error()
	}

	_ = json.NewEncoder(w).Encode(resp)
}
