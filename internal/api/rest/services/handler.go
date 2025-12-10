package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"booking-service/internal/api/rest/helpers"
	"booking-service/internal/store/business_accounts"
	"booking-service/internal/store/services"

	"github.com/gorilla/mux"
)

type Handler struct {
	servicesStore         services.Store
	businessAccountsStore business_accounts.Store
}

func NewHandler(servicesStore services.Store, businessAccountsStore business_accounts.Store) *Handler {
	return &Handler{
		servicesStore:         servicesStore,
		businessAccountsStore: businessAccountsStore,
	}
}

func (h *Handler) CreateService(resp http.ResponseWriter, req *http.Request) {
	var createReq services.CreateServiceRequest
	if err := json.NewDecoder(req.Body).Decode(&createReq); err != nil {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Invalid request body", helpers.InvalidRequest),
			http.StatusBadRequest,
		)
		return
	}

	// Validate required fields
	if err := validateCreateServiceRequest(createReq); err != nil {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse(err.Error(), helpers.ValidationError),
			http.StatusBadRequest,
		)
		return
	}

	// Verify business account exists and user owns it
	// Note: In a real implementation, you'd get the user ID from JWT token
	// For now, we'll assume the business account ID is valid
	_, err := h.businessAccountsStore.GetBusinessAccount(req.Context(), createReq.BusinessAccountID)
	if err != nil {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Business account not found", helpers.NotFound),
			http.StatusNotFound,
		)
		return
	}

	service, err := h.servicesStore.CreateService(req.Context(), createReq)
	if err != nil {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Failed to create service", helpers.InternalError),
			http.StatusInternalServerError,
		)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(resp).Encode(service); err != nil {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Failed to encode response", helpers.InternalError),
			http.StatusInternalServerError,
		)
		return
	}
}

func (h *Handler) GetService(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	serviceID := vars["id"]

	if serviceID == "" {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Service ID is required", helpers.ValidationError),
			http.StatusBadRequest,
		)
		return
	}

	service, err := h.servicesStore.GetService(req.Context(), serviceID)
	if err != nil {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Failed to get service", helpers.InternalError),
			http.StatusInternalServerError,
		)
		return
	}

	if service == nil {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Service not found", helpers.NotFound),
			http.StatusNotFound,
		)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(resp).Encode(service); err != nil {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Failed to encode response", helpers.InternalError),
			http.StatusInternalServerError,
		)
		return
	}
}

func (h *Handler) UpdateService(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	serviceID := vars["id"]

	if serviceID == "" {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Service ID is required", helpers.ValidationError),
			http.StatusBadRequest,
		)
		return
	}

	var updateReq services.UpdateServiceRequest
	if err := json.NewDecoder(req.Body).Decode(&updateReq); err != nil {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Invalid request body", helpers.InvalidRequest),
			http.StatusBadRequest,
		)
		return
	}

	// Validate update request
	if err := validateUpdateServiceRequest(updateReq); err != nil {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse(err.Error(), helpers.ValidationError),
			http.StatusBadRequest,
		)
		return
	}

	// Verify service exists and user owns it
	existingService, err := h.servicesStore.GetService(req.Context(), serviceID)
	if err != nil {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Failed to get service", helpers.InternalError),
			http.StatusInternalServerError,
		)
		return
	}

	if existingService == nil {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Service not found", helpers.NotFound),
			http.StatusNotFound,
		)
		return
	}

	// Note: In a real implementation, you'd verify the user owns this service
	// by checking the business account ownership

	service, err := h.servicesStore.UpdateService(req.Context(), serviceID, updateReq)
	if err != nil {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Failed to update service", helpers.InternalError),
			http.StatusInternalServerError,
		)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(resp).Encode(service); err != nil {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Failed to encode response", helpers.InternalError),
			http.StatusInternalServerError,
		)
		return
	}
}

func (h *Handler) DeleteService(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	serviceID := vars["id"]

	if serviceID == "" {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Service ID is required", helpers.ValidationError),
			http.StatusBadRequest,
		)
		return
	}

	// Verify service exists and user owns it
	existingService, err := h.servicesStore.GetService(req.Context(), serviceID)
	if err != nil {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Failed to get service", helpers.InternalError),
			http.StatusInternalServerError,
		)
		return
	}

	if existingService == nil {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Service not found", helpers.NotFound),
			http.StatusNotFound,
		)
		return
	}

	// Note: In a real implementation, you'd verify the user owns this service
	// by checking the business account ownership

	err = h.servicesStore.DeleteService(req.Context(), serviceID)
	if err != nil {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Failed to delete service", helpers.InternalError),
			http.StatusInternalServerError,
		)
		return
	}

	resp.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListServices(resp http.ResponseWriter, req *http.Request) {
	// Parse query parameters
	queries := req.URL.Query()

	limitStr := queries.Get("limit")
	offsetStr := queries.Get("offset")
	businessAccountID := queries.Get("business_account_id")
	category := queries.Get("category")
	isActiveStr := queries.Get("is_active")

	// Set default values
	limit := 20
	offset := 0

	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Build request
	listReq := services.ListServicesRequest{
		Limit:  limit,
		Offset: offset,
	}

	if businessAccountID != "" {
		listReq.BusinessAccountID = &businessAccountID
	}

	if category != "" {
		listReq.Category = &category
	}

	if isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			listReq.IsActive = &isActive
		}
	}

	// Validate request
	if err := validateListServicesRequest(listReq); err != nil {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse(err.Error(), helpers.ValidationError),
			http.StatusBadRequest,
		)
		return
	}

	result, err := h.servicesStore.ListServices(req.Context(), listReq)
	if err != nil {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Failed to list services", helpers.InternalError),
			http.StatusInternalServerError,
		)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(resp).Encode(result); err != nil {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Failed to encode response", helpers.InternalError),
			http.StatusInternalServerError,
		)
		return
	}
}

func (h *Handler) GetServicesByBusinessAccount(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	businessAccountID := vars["business_account_id"]

	if businessAccountID == "" {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Business account ID is required", helpers.ValidationError),
			http.StatusBadRequest,
		)
		return
	}

	services, err := h.servicesStore.GetServicesByBusinessAccount(req.Context(), businessAccountID)
	if err != nil {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Failed to get services", helpers.InternalError),
			http.StatusInternalServerError,
		)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(resp).Encode(services); err != nil {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Failed to encode response", helpers.InternalError),
			http.StatusInternalServerError,
		)
		return
	}
}

// TODO: not added to docs yet
// Validation functions
func validateCreateServiceRequest(req services.CreateServiceRequest) error {
	if req.BusinessAccountID == "" {
		return helpers.NewValidationError("business_account_id is required")
	}
	if req.Name == "" {
		return helpers.NewValidationError("name is required")
	}
	if req.DurationMinutes <= 0 {
		return helpers.NewValidationError("duration_minutes must be greater than 0")
	}
	if req.Price < 0 {
		return helpers.NewValidationError("price cannot be negative")
	}
	if req.Currency != "" && len(req.Currency) != 3 {
		return helpers.NewValidationError("currency must be a 3-character code")
	}
	return nil
}

func validateUpdateServiceRequest(req services.UpdateServiceRequest) error {
	if req.Name != nil && *req.Name == "" {
		return helpers.NewValidationError("name cannot be empty")
	}
	if req.DurationMinutes != nil && *req.DurationMinutes <= 0 {
		return helpers.NewValidationError("duration_minutes must be greater than 0")
	}
	if req.Price != nil && *req.Price < 0 {
		return helpers.NewValidationError("price cannot be negative")
	}
	if req.Currency != nil && len(*req.Currency) != 3 {
		return helpers.NewValidationError("currency must be a 3-character code")
	}
	return nil
}

func validateListServicesRequest(req services.ListServicesRequest) error {
	if req.Limit <= 0 {
		return helpers.NewValidationError("limit must be greater than 0")
	}
	if req.Limit > 100 {
		return helpers.NewValidationError("limit cannot exceed 100")
	}
	if req.Offset < 0 {
		return helpers.NewValidationError("offset cannot be negative")
	}
	return nil
}
