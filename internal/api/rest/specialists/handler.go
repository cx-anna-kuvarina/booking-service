package specialists

import (
	"net/http"

	"booking-service/internal/api/rest/helpers"
	"github.com/rs/zerolog/log"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetSpecialists(resp http.ResponseWriter, req *http.Request) {
	queries := req.URL.Query()

	areaType, city, err := validateGetSpecialistsQueries(queries)
	if err != nil {
		helpers.WriteErrorResponse(resp, helpers.NewErrorResponse(err.Error(), helpers.InvalidQueries), http.StatusBadRequest)
		return
	}

	log.Ctx(req.Context()).Info().Msgf("areatype %s, city %s", areaType, city)
}
