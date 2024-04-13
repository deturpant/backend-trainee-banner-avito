package banners

import (
	"backend-trainee-banner-avito/internal/lib/api/response"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type RequestGetBanners struct {
	FeatureID *int `json:"feature_id"`
	TagID     *int `json:"tag_id"`
	Limit     *int `json:"limit"`
	Offset    *int `json:"offset"`
}

// @Summary Получить баннеры с учетом заданных параметров
// @Description Получает список баннеров с учетом переданных параметров запроса
// @ID get-banners
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Токен администратора" format:"Bearer <admin_token>"
// @Param feature_id query int false "Идентификатор особенности баннера"
// @Param tag_id query int false "Идентификатор тега баннера"
// @Param limit query int false "Ограничение на количество возвращаемых баннеров"
// @Param offset query int false "Смещение для пагинации результатов"
// @Success 200 {array} entities.Banner "Список баннеров"
// @Failure 500 {object} banners.Response "Ошибка при получении баннеров"
// @Router /banner [get]
func NewGetBannersHandler(bannerRepo Banners, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := parseGetBannersRequest(r)

		banners, err := bannerRepo.FindBannersByParameters(r.Context(), req)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			logger.Error("Failed to get banners")
			render.JSON(w, r, response.Error("Failed to get banners"))
			return
		}

		render.JSON(w, r, banners)
	}
}

func parseGetBannersRequest(r *http.Request) RequestGetBanners {
	req := RequestGetBanners{}

	// Парсим параметры запроса и преобразуем в числа
	if featureIDStr := r.URL.Query().Get("feature_id"); featureIDStr != "" {
		featureID, _ := strconv.Atoi(featureIDStr)
		req.FeatureID = &featureID
	}

	if tagIDStr := r.URL.Query().Get("tag_id"); tagIDStr != "" {
		tagID, _ := strconv.Atoi(tagIDStr)
		req.TagID = &tagID
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, _ := strconv.Atoi(limitStr)
		req.Limit = &limit
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, _ := strconv.Atoi(offsetStr)
		req.Offset = &offset
	}

	return req
}
