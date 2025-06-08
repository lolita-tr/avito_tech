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

func NewHandler(store *StoreService, coins *CoinsService, info *Info) *Handler {

	return &Handler{store: store, coins: coins, info: info}
}

func (h *Handler) BuyItem(w http.ResponseWriter, r *http.Request) {
	itemName := chi.URLParam(r, "item")

	claims, ok := r.Context().Value("jwt_claims").(*auth.Claims)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	response, err := h.store.BuyItem(r.Context(), claims.UserID, itemName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) SendCoins(w http.ResponseWriter, r *http.Request) {
	var request SendCoinsRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	claims, ok := r.Context().Value("jwt_claims").(*auth.Claims)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	response, err := h.coins.SendCoins(r.Context(), claims.UserID, request.ToUser, request.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) GetInfo(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("jwt_claims").(*auth.Claims)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	response, err := h.info.GetUserInfo(r.Context(), claims.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
