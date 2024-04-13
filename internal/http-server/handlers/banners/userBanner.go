package banners

import (
	"backend-trainee-banner-avito/internal/entities"
	"backend-trainee-banner-avito/internal/lib/api/response"
	"backend-trainee-banner-avito/internal/lib/logger/errMsg"
	"backend-trainee-banner-avito/internal/storage/cache"
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
type Content struct {
	Text string `json:"text"`
	Url  string `json:"url"`
	Name string `json:"name"`
}

// @Summary Получить баннер с указанными параметрами
// @Description Получает баннер с учетом переданных параметров запроса
// @ID get-banner
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Токен пользователя" format:"Bearer <user_token>"
// @Param feature_id query int true "Идентификатор особенности баннера"
// @Param tag_id query int true "Идентификатор тега баннера"
// @Param use_last_revision query bool false "Использовать последнюю ревизию"
// @Success 200 {object} Content "Содержимое баннера"
// @Failure 400 {object} banners.Response "Неверные параметры запроса"
// @Failure 404 {object} banners.Response "Баннер не найден"
// @Router /user_banner [get]
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
		if req.UseLastRevision {
			banner, err := bannerRepo.FindBannerByFeatureTag(r.Context(), req.FeatureID, req.TagID)
			if err != nil {
				log.Error("Failed to find banner", errMsg.Err(err))
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, response.Error("Failed to find banner"))
				return
			}
			cache.StoreBannerInCache(req.FeatureID, req.TagID, *banner)
			responseGetOK(w, r, *banner)
		} else {
			banner, found := cache.GetBannerFromCache(req.FeatureID, req.TagID)
			if found {
				log.Info("Banner found in CACHE")
				responseGetOK(w, r, *banner)
			} else {
				banner, err := bannerRepo.FindBannerByFeatureTag(r.Context(), req.FeatureID, req.TagID)
				if err != nil {
					log.Error("Failed to find banner", errMsg.Err(err))
					render.Status(r, http.StatusNotFound)
					render.JSON(w, r, response.Error("Failed to find banner"))
					return
				}
				cache.StoreBannerInCache(req.FeatureID, req.TagID, *banner)
				responseGetOK(w, r, *banner)
			}
		}

	}
}

func responseGetOK(w http.ResponseWriter, r *http.Request, banner entities.Banner) {
	render.JSON(w, r, banner.Content)
}
