package banners

import (
	"backend-trainee-banner-avito/internal/entities"
	"backend-trainee-banner-avito/internal/lib/api/response"
	"backend-trainee-banner-avito/internal/lib/logger/errMsg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"strconv"
)

type RequestGetBanner struct {
	FeatureID       int  `json:"feature_id" validate:"required"`
	TagID           int  `json:"tag_id" validate:"required"`
	UseLastRevision bool `json:"use_last_revision"`
}

func NewGetBannerHandler(log *slog.Logger, bannerRepo Banners) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.banners.userBanner.New"
		log := log.With(
			slog.String("options", loggerOptions),
			slog.String("request_id", chi.URLParam(r, "request_id")))

		featureIDStr := r.URL.Query().Get("feature_id")
		tagIDStr := r.URL.Query().Get("tag_id")
		useLastRevisionStr := r.URL.Query().Get("use_last_revision")

		featureID, err := strconv.Atoi(featureIDStr)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("Invalid feature_id"))
			return
		}

		tagID, err := strconv.Atoi(tagIDStr)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("Invalid tag_id"))
			return
		}

		useLastRevision, err := strconv.ParseBool(useLastRevisionStr)
		if err != nil {
			useLastRevision = false
		}

		req := RequestGetBanner{
			FeatureID:       featureID,
			TagID:           tagID,
			UseLastRevision: useLastRevision,
		}

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("Invalid request", errMsg.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}

		banner, err := bannerRepo.FindBannerByFeatureTag(r.Context(), req.FeatureID, req.TagID)
		if err != nil {
			log.Error("Failed to find banner", errMsg.Err(err))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Error("Failed to find banner"))
			return
		}

		responseGetOK(w, r, *banner)
	}
}

func responseGetOK(w http.ResponseWriter, r *http.Request, banner entities.Banner) {
	render.JSON(w, r, banner.Content)
}
