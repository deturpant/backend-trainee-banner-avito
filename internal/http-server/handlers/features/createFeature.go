package features

import (
	"backend-trainee-banner-avito/internal/entities"
	"backend-trainee-banner-avito/internal/lib/api/response"
	"backend-trainee-banner-avito/internal/lib/logger/errMsg"
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Features interface {
	CreateFeature(ctx context.Context, feature *entities.Feature) error
}
type RequestFeature struct {
	Name string `json:"name" validate:"required"`
}
type ResponseFeature struct {
	response.Response
	ID   int    `json:"feature_id"`
	Name string `json:"name"`
}

// @Summary Создать новую фияу
// @Description Создает новую фичу на основе переданных данных
// @ID create-feature
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Токен администратора" format:"Bearer <admin_token>"
// @Param request body RequestFeature true "Данные для создания фичи"
// @Success 200 {object} ResponseFeature "Созданная фича"
// @Failure 400 {object} banners.Response "Неверные параметры запроса"
// @Failure 500 {object} banners.Response "Ошибка при создании фича"
// @Router /features [post]
func New(log *slog.Logger, featureRepository Features) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.features.createFeature.New"
		log = log.With(
			slog.String("options", loggerOptions),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req RequestFeature
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to decode request body", errMsg.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("Failed to decode request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("Invalid request", errMsg.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}
		feature := entities.Feature{Name: req.Name}
		err = featureRepository.CreateFeature(r.Context(), &feature)
		if err != nil {
			log.Error("Failed to create feature", errMsg.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("Failed to create feature"))
			return
		}
		log.Info("Feature added")
		responseOK(w, r, req.Name, feature.ID)
	}
}
func responseOK(w http.ResponseWriter, r *http.Request, name string, feature_id int) {
	render.JSON(w, r, ResponseFeature{Response: response.OK(),
		Name: name, ID: feature_id})
}
