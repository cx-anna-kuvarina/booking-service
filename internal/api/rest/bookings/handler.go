package bookings

import (
	"encoding/json"
	"net/http"
	"time"

	"booking-service/internal/store/bookings"

	"github.com/gorilla/mux"
)

type Handler struct {
	store bookings.Store
}

func NewHandler(store bookings.Store) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) CreateBooking(resp http.ResponseWriter, req *http.Request) {
	var createReq bookings.CreateBookingRequest
	if err := json.NewDecoder(req.Body).Decode(&createReq); err != nil {
		http.Error(resp, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if createReq.UserID == "" {
		http.Error(resp, "user_id is required", http.StatusBadRequest)
		return
	}
	if createReq.BusinessID == "" {
		http.Error(resp, "business_id is required", http.StatusBadRequest)
		return
	}
	if createReq.ServiceID == "" {
		http.Error(resp, "service_id is required", http.StatusBadRequest)
		return
	}
	if createReq.StartTime.IsZero() {
		http.Error(resp, "start_time is required", http.StatusBadRequest)
		return
	}
	if createReq.EndTime.IsZero() {
		http.Error(resp, "end_time is required", http.StatusBadRequest)
		return
	}

	// Validate time range
	if createReq.EndTime.Before(createReq.StartTime) {
		http.Error(resp, "end_time must be after start_time", http.StatusBadRequest)
		return
	}

	// Validate booking is not in the past
	if createReq.StartTime.Before(time.Now()) {
		http.Error(resp, "cannot create booking in the past", http.StatusBadRequest)
		return
	}

	booking, err := h.store.CreateBooking(req.Context(), createReq)
	if err != nil {
		http.Error(resp, "failed to create booking", http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(resp).Encode(booking); err != nil {
		http.Error(resp, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetBooking(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	bookingID := vars["id"]

	if bookingID == "" {
		http.Error(resp, "booking ID is required", http.StatusBadRequest)
		return
	}

	booking, err := h.store.GetBooking(req.Context(), bookingID)
	if err != nil {
		http.Error(resp, "failed to get booking", http.StatusInternalServerError)
		return
	}

	if booking == nil {
		http.Error(resp, "booking not found", http.StatusNotFound)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(resp).Encode(booking); err != nil {
		http.Error(resp, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
