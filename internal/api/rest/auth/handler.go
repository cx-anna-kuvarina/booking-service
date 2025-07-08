package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"booking-service/internal/api/rest/helpers"
	"booking-service/internal/store/users"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

const googleGetUserURL = "https://www.googleapis.com/oauth2/v2/userinfo"

type Handler struct {
	googleConfig oauth2.Config
	randomState  string
	jwtSecret    string
	jwtExpPeriod time.Duration
	uStore       users.Store
}

func NewHandler(googleConfig oauth2.Config, randomState, jwtSecret string, jwtExpPeriod time.Duration,
	uStore users.Store) *Handler {
	return &Handler{
		googleConfig: googleConfig,
		randomState:  randomState,
		jwtSecret:    jwtSecret,
		jwtExpPeriod: jwtExpPeriod,
		uStore:       uStore,
	}
}

func (h *Handler) GoogleLogin(resp http.ResponseWriter, req *http.Request) {
	url := h.googleConfig.AuthCodeURL(h.randomState)
	http.Redirect(resp, req, url, http.StatusTemporaryRedirect)
}

func (h *Handler) GoogleCallback(resp http.ResponseWriter, req *http.Request) {
	state := req.FormValue("state")
	if state != h.randomState {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Invalid state parameter", helpers.InvalidStateParameter),
			http.StatusBadRequest,
		)
		return
	}

	ctx := context.Background()

	// Exchange code for token
	code := req.FormValue("code")
	token, err := h.googleConfig.Exchange(ctx, code)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("Failed to exchange token by code: %s", code)
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Failed to exchange token", helpers.ExchangeTokenErr),
			http.StatusInternalServerError,
		)
	}

	// Get user info
	client := h.googleConfig.Client(ctx, token)
	userResp, err := client.Get(googleGetUserURL)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get user info")
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Failed to get user info", helpers.AuthUserInfoErr),
			http.StatusInternalServerError,
		)
	}

	defer userResp.Body.Close()

	var userInfo *UserInfo
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}
	defer req.Body.Close()

	if err = json.Unmarshal(body, &userInfo); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to decode user info")
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Error decoding info", helpers.DecodeUserInfoErr),
			http.StatusInternalServerError,
		)
		return
	}

	if userInfo.Email == "" {
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Email is required", helpers.InvalidEmailErr),
			http.StatusBadRequest,
		)
	}

	userID, err := h.uStore.GetUserIdByEmail(ctx, userInfo.Email)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("Failed to get user id by email: %s", userInfo.Email)
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Failed to get user id", helpers.GetUserIdErr),
			http.StatusInternalServerError,
		)
		return
	}

	jwtToken, err := generateJWT(userID, h.jwtSecret, h.jwtExpPeriod)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("Failed to generate JWT token for user: %s", userID)
		helpers.WriteErrorResponse(
			resp,
			helpers.NewErrorResponse("Failed to generate token", helpers.GenerateTokenErr),
			http.StatusInternalServerError,
		)
		return
	}

	helpers.WriteData(ctx, resp, map[string]string{
		"access_token": jwtToken,
	}, http.StatusOK)
}
