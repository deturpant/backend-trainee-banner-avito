package banners

import (
	"backend-trainee-banner-avito/internal/lib/api/response"
	"backend-trainee-banner-avito/internal/lib/logger/errMsg"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type RequestUpdateBanner struct {
	TagIDs    []int                  `json:"tag_ids" validate:"required"`
	FeatureID int                    `json:"feature_id" validate:"required"`
	Content   map[string]interface{} `json:"content" validate:"required"`
	IsActive  bool                   `json:"is_active" validate:"required"`
}

// @Summary Обновить информацию о баннере
// @Description Обновляет информацию о баннере с указанным идентификатором
// @ID update-banner
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Токен администратора" format:"Bearer <admin_token>"
// @Param id path int true "Идентификатор баннера"
// @Param request body RequestUpdateBanner true "Данные для обновления баннера"
// @Success 200 {object} ResponseBanner "Обновленный баннер"
// @Failure 400 {object} banners.Response "Неверный идентификатор баннера или неверные данные запроса"
// @Failure 404 {object} banners.Response "Баннер не найден"
// @Failure 500 {object} banners.Response "Ошибка при обновлении баннера"
// @Router /banner/{id} [patch]
func NewUpdateBannerHandler(bannerRepo Banners, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bannerID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			logger.Error("Invalid banner ID")
			render.JSON(w, r, response.Error("Invalid banner ID"))
			return
		}

		var req RequestUpdateBanner
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			render.Status(r, http.StatusBadRequest)
			logger.Error("Failed to parse request body")
			render.JSON(w, r, response.Error("Failed to parse request body"))
			return
		}
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			logger.Error("Invalid request", errMsg.Err(err))
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}
		banner, err := bannerRepo.FindBannerById(r.Context(), bannerID)
		if err != nil {
			render.Status(r, http.StatusNotFound)
			logger.Error("Failed to find banner")
			render.JSON(w, r, response.Error("Banner not found"))
			return
		}

		banner.TagIDs = req.TagIDs
		banner.FeatureID = req.FeatureID
		banner.Content = req.Content
		banner.IsActive = req.IsActive
		banner.UpdatedAt = time.Now()

		if err := bannerRepo.UpdateBanner(r.Context(), &banner); err != nil {
			render.Status(r, http.StatusInternalServerError)
			logger.Error("Failed to update banner")
			render.JSON(w, r, response.Error("Failed to update banner"))
			return
		}
		responseOK(w, r, banner)
	}
}
