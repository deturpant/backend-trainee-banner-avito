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
type Request struct {
	Name string `json:"name" validate:"required"`
}
type Response struct {
	response.Response
	ID   int    `json:"feature_id"`
	Name string `json:"name"`
}

func New(log *slog.Logger, featureRepository Features) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.features.createFeature.New"
		log = log.With(
			slog.String("options", loggerOptions),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req Request
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
		feature := entities.Feature{Name: req.Name}
		err = featureRepository.CreateFeature(r.Context(), &feature)
		if err != nil {
			log.Error("Failed to create feature", errMsg.Err(err))
			render.JSON(w, r, response.Error("Failed to create feature"))
			return
		}
		log.Info("Feature added")
		responseOK(w, r, req.Name)
	}
}
func responseOK(w http.ResponseWriter, r *http.Request, name string) {
	render.JSON(w, r, Response{Response: response.OK(),
		Name: name})
}
