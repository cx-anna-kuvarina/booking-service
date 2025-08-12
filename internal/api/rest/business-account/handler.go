package business_account

import (
	"encoding/json"
	"net/http"

	"booking-service/internal/store/business_accounts"
	utoken "booking-service/pkg/utils/token"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	store business_accounts.Store
}

func NewHandler(store business_accounts.Store) *Handler {
	return &Handler{
		store: store,
	}
}

type CreateBusinessAccountRequest struct {
	Name         string          `json:"name"`
	BusinessType string          `json:"businessType"`
	Location     string          `json:"location"`
	Links        json.RawMessage `json:"links"`
	UserID       string          `json:"userId"`
}

type UpdateBusinessAccountRequest struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	BusinessType string          `json:"businessType"`
	Location     string          `json:"location"`
	Links        json.RawMessage `json:"links"`
}

func (h *Handler) CreateBusinessAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateBusinessAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.BusinessType == "" || req.Location == "" || req.UserID == "" {
		http.Error(w, "Name, business type, location and user ID are required", http.StatusBadRequest)
		return
	}

	account := &business_accounts.BusinessAccount{
		ID:           uuid.New().String(),
		Name:         req.Name,
		BusinessType: req.BusinessType,
		Location:     req.Location,
		Links:        req.Links,
	}

	if err := h.store.CreateBusinessAccount(r.Context(), account, req.UserID); err != nil {
		http.Error(w, "Failed to create business account", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

func (h *Handler) UpdateBusinessAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	businessAccountID, ok := vars["id"]
	if !ok {
		http.Error(w, "Business account ID is required", http.StatusBadRequest)
		return
	}

	var req UpdateBusinessAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Extract user ID from JWT (Authorization header)
	token := r.Header.Get("Authorization")
	userID, err := utoken.ExtractUserIDFromJWT(token)
	if err != nil {
		http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
		return
	}

	// Validate user owns the business account using store method
	owns, err := h.store.UserOwnsBusinessAccount(r.Context(), businessAccountID, userID)
	if err != nil {
		http.Error(w, "Failed to validate ownership", http.StatusInternalServerError)
		return
	}
	if !owns {
		http.Error(w, "Forbidden: you do not own this business account", http.StatusForbidden)
		return
	}

	if req.ID != businessAccountID {
		http.Error(w, "ID does not match business account ID", http.StatusBadRequest)
		return
	}

	if req.ID == "" || req.Name == "" || req.BusinessType == "" || req.Location == "" {
		http.Error(w, "ID, name, business type, and location are required", http.StatusBadRequest)
		return
	}

	account := &business_accounts.BusinessAccount{
		ID:           req.ID,
		Name:         req.Name,
		BusinessType: req.BusinessType,
		Location:     req.Location,
		Links:        req.Links,
	}

	if err := h.store.UpdateBusinessAccount(r.Context(), account); err != nil {
		http.Error(w, "Failed to update business account", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

func (h *Handler) DeleteBusinessAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	businessAccountID, ok := vars["id"]
	if !ok {
		http.Error(w, "Business account ID is required", http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	// Extract user ID from JWT (Authorization header)
	token := r.Header.Get("Authorization")
	userID, err := utoken.ExtractUserIDFromJWT(token)
	if err != nil {
		http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
		return
	}

	// Validate user owns the business account using store method
	owns, err := h.store.UserOwnsBusinessAccount(r.Context(), businessAccountID, userID)
	if err != nil {
		http.Error(w, "Failed to validate ownership", http.StatusInternalServerError)
		return
	}
	if !owns {
		http.Error(w, "Forbidden: you do not own this business account", http.StatusForbidden)
		return
	}

	err = h.store.DeleteBusinessAccount(r.Context(), businessAccountID)
	if err != nil {
		if err == business_accounts.NotFoundError {
			http.Error(w, "Business account not found", http.StatusNotFound)
			return
		}
		log.Ctx(ctx).Error().Err(err).Msgf("Failed to delete business account: %s", err.Error())
		http.Error(w, "Failed to delete business account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetBusinessAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	businessAccountID, ok := vars["id"]
	if !ok {
		http.Error(w, "Business account ID is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	account, err := h.store.GetBusinessAccount(r.Context(), businessAccountID)
	if err != nil {
		if err == business_accounts.NotFoundError {
			http.Error(w, "Business account not found", http.StatusNotFound)
			return
		}
		log.Ctx(ctx).Error().Err(err).Msgf("Failed to delete business account: %s", err.Error())
		http.Error(w, "Failed to get business account", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}
