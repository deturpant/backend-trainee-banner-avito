package users

import (
	"backend-trainee-banner-avito/internal/entities"
	"backend-trainee-banner-avito/internal/lib/api/response"
	"backend-trainee-banner-avito/internal/lib/auth"
	"backend-trainee-banner-avito/internal/lib/logger/errMsg"
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type User interface {
	CreateUser(ctx context.Context, user *entities.User) error
	FindUserByUsername(ctx context.Context, username string) (entities.User, error)
}

type RequestUser struct {
	Username string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type ResponseUser struct {
	response.Response
	ID   int    `json:"user_id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

// @Summary Создать нового пользователя
// @Description Создает нового пользователя на основе переданных данных
// @ID create-user
// @Accept  json
// @Produce  json
// @Param request body RequestUser true "Данные для создания пользователя"
// @Success 200 {object} ResponseUser "Созданный пользователь"
// @Failure 400 {object} banners.Response "Неверные параметры запроса"
// @Failure 500 {object} banners.Response "Ошибка при создании пользователя"
// @Router /users [post]
func New(log *slog.Logger, userRepository User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.features.createUser.New"
		log = log.With(
			slog.String("options", loggerOptions),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req RequestUser
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
		hashPass, err := auth.HashPassword(req.Password)
		user := entities.User{Username: req.Username, Password: hashPass, Role: "user"}
		err = userRepository.CreateUser(r.Context(), &user)
		if err != nil {
			log.Error("Failed to create user", errMsg.Err(err))
			render.JSON(w, r, response.Error("Failed to create user"))
			return
		}
		log.Info("User added")
		responseOK(w, r, req.Username, user.ID, user.Role)
	}
}
func responseOK(w http.ResponseWriter, r *http.Request, name string, userID int, role string) {
	render.JSON(w, r, ResponseUser{Response: response.OK(),
		Name: name, ID: userID, Role: role})
}
