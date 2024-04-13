package tags

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

type RequestTag struct {
	Name string `json:"name" validate:"required"`
}
type ResponseTag struct {
	response.Response
	ID   int    `json:"tag_id"`
	Name string `json:"name"`
}

type Tag interface {
	CreateTag(ctx context.Context, tag *entities.Tag) error
}

// @Summary Создать новый тег
// @Description Создает новый тег на основе переданных данных
// @ID create-tag
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Токен администратора" format:"Bearer <admin_token>"
// @Param request body RequestTag true "Данные для создания тега"
// @Success 200 {object} ResponseTag "Созданный тег"
// @Failure 400 {object} banners.Response "Неверные параметры запроса"
// @Failure 500 {object} banners.Response "Ошибка при создании тега"
// @Router /tags [post]
func New(log *slog.Logger, tagRepository Tag) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.features.createTag.New"
		log = log.With(
			slog.String("options", loggerOptions),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req RequestTag
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to decode request body", errMsg.Err(err))
			render.JSON(w, r, response.Error("Failed to decode request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("Invalid request", errMsg.Err(err))
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}
		tag := entities.Tag{Name: req.Name}
		err = tagRepository.CreateTag(r.Context(), &tag)
		if err != nil {
			log.Error("Failed to create tag", errMsg.Err(err))
			render.JSON(w, r, response.Error("Failed to create tag"))
			return
		}
		log.Info("Tag added")
		responseOK(w, r, req.Name, tag.ID)
	}
}
func responseOK(w http.ResponseWriter, r *http.Request, name string, tag_id int) {
	render.JSON(w, r, ResponseTag{Response: response.OK(),
		Name: name, ID: tag_id})
}
