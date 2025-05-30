package service

import (
	"avito_tech/internal/auth"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type Handler struct {
	store *StoreService
}

func NewHandler(store *StoreService) *Handler {
	return &Handler{store: store}
}

func (h *Handler) BuyItem(w http.ResponseWriter, r *http.Request) {
	itemName := chi.URLParam(r, "item")

	log.Printf("%s", itemName)

	claims, ok := r.Context().Value("jwt_claims").(*auth.Claims)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	log.Printf("%s", claims.UserID)

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
