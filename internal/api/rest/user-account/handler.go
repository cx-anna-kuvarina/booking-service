package user_account

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"booking-service/internal/api/rest/helpers"
	"booking-service/internal/store/users"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	usersStore users.Store
}

func NewHandler(usersStore users.Store) *Handler {
	return &Handler{
		usersStore: usersStore,
	}
}

func (h *Handler) GetUserAccount(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserIDFromJWT(r.Header.Get("Authorization"))

	if err != nil {
		helpers.WriteErrorResponse(w, helpers.NewErrorResponse(err.Error(), helpers.InvalidTokenErr), http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	user, err := h.usersStore.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			log.Ctx(ctx).Error().Err(err).Msgf("user %s not found", userID)
			helpers.WriteErrorResponse(w, helpers.NewErrorResponse(err.Error(), helpers.NotFoundErr), http.StatusNotFound)
			return
		}

		log.Ctx(ctx).Error().Err(err).Msgf("failed to get user %s", userID)
		helpers.WriteErrorResponse(w, helpers.NewErrorResponse(err.Error(), helpers.UsersStoreErr), http.StatusInternalServerError)
		return
	}

	helpers.WriteData(ctx, w, user, http.StatusOK)
}

func (h *Handler) UpdateUserAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, ok := vars["id"]
	if !ok {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	userIDToken, err := helpers.GetUserIDFromJWT(r.Header.Get("Authorization"))
	if err != nil {
		helpers.WriteErrorResponse(w, helpers.NewErrorResponse(err.Error(), helpers.InvalidTokenErr), http.StatusUnauthorized)
		return
	}

	if userIDToken != userID {
		http.Error(w, "User ID does not match", http.StatusForbidden)
		return
	}

	ctx := r.Context()
	uStore, err := h.usersStore.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			log.Ctx(ctx).Error().Err(err).Msgf("user %s not found", userID)
			helpers.WriteErrorResponse(w, helpers.NewErrorResponse(err.Error(), helpers.NotFoundErr), http.StatusNotFound)
			return
		}

		log.Ctx(ctx).Error().Err(err).Msgf("failed to get user %s", userID)
		helpers.WriteErrorResponse(w, helpers.NewErrorResponse(err.Error(), helpers.UsersStoreErr), http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		helpers.WriteErrorResponse(w, helpers.NewErrorResponse(err.Error(), helpers.ReadBodyErr), http.StatusBadRequest)
		return
	}

	var user users.User
	if err := json.Unmarshal(body, &user); err != nil {
		helpers.WriteErrorResponse(w, helpers.NewErrorResponse(err.Error(), helpers.DecodeUserInfoErr), http.StatusBadRequest)
		return
	}

	if user.Email != uStore.Email {
		helpers.WriteErrorResponse(w, helpers.NewErrorResponse("email is not allowed to be changed", helpers.InvalidEmailErr), http.StatusBadRequest)
		return
	}

	user.ID = userID

	if err := h.usersStore.UpdateUser(ctx, &user); err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			log.Ctx(ctx).Error().Err(err).Msgf("user %s not found", userID)
			helpers.WriteErrorResponse(w, helpers.NewErrorResponse(err.Error(), helpers.NotFoundErr), http.StatusNotFound)
			return
		}

		log.Ctx(ctx).Error().Err(err).Msgf("failed to update user %s", userID)
		helpers.WriteErrorResponse(w, helpers.NewErrorResponse(err.Error(), helpers.UsersStoreErr), http.StatusInternalServerError)
		return
	}

	helpers.WriteData(ctx, w, user, http.StatusOK)
}

func (h *Handler) DeleteUserAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, ok := vars["id"]
	if !ok {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userIDToken, err := helpers.GetUserIDFromJWT(r.Header.Get("Authorization"))
	if err != nil {
		helpers.WriteErrorResponse(w, helpers.NewErrorResponse(err.Error(), helpers.InvalidTokenErr), http.StatusUnauthorized)
		return
	}

	if userIDToken != userID {
		http.Error(w, "User ID does not match", http.StatusForbidden)
		return
	}

	ctx := r.Context()

	err = h.usersStore.DeleteUser(ctx, userID)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			log.Ctx(ctx).Error().Err(err).Msgf("user %s not found", userID)
			helpers.WriteErrorResponse(w, helpers.NewErrorResponse(err.Error(), helpers.NotFoundErr), http.StatusNotFound)
			return
		}

		log.Ctx(ctx).Error().Err(err).Msgf("failed to delete user %s", userID)
		helpers.WriteErrorResponse(w, helpers.NewErrorResponse(err.Error(), helpers.UsersStoreErr), http.StatusInternalServerError)
		return
	}

	helpers.WriteData(ctx, w, nil, http.StatusOK)
}
