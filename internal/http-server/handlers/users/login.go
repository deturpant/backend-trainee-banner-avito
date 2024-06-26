package users

import (
	"backend-trainee-banner-avito/internal/lib/api/response"
	"backend-trainee-banner-avito/internal/lib/auth"
	"backend-trainee-banner-avito/internal/lib/logger/errMsg"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"time"
)

type ResponseAuthUser struct {
	response.Response
	ID    int    `json:"user_id"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	Token string `json:"token"`
}

// @Summary Войти в систему
// @Description Аутентифицирует пользователя и выдает токен доступа
// @ID login
// @Accept  json
// @Produce  json
// @Param request body RequestUser true "Данные для входа в систему"
// @Success 200 {object} ResponseAuthUser "Аутентифицированный пользователь с токеном доступа"
// @Failure 400 {object} banners.Response "Неверные параметры запроса"
// @Failure 401 {object} banners.Response "Неверные учетные данные"
// @Router /login [post]
func LoginFunc(log *slog.Logger, userRepository User, jwt *auth.JWTManager) http.HandlerFunc {
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
		user, err := userRepository.FindUserByUsername(r.Context(), req.Username)
		if err != nil {
			log.Error("User not found with login")
			render.JSON(w, r, response.Error("Invalid username"))

			return
		}

		errAuth := auth.ComparePasswordHash(req.Password, user.Password)
		if errAuth != nil {
			log.Error("Invalid password")
			render.JSON(w, r, response.Error("Invalid password"))
			return
		}
		token, err := jwt.GenerateToken(user.Username, user.Role, time.Second*600)
		log.Info("User authenticated")
		responseAuthOK(w, r, req.Username, user.ID, user.Role, token)
	}
}
func responseAuthOK(w http.ResponseWriter, r *http.Request, name string, userID int, role string, token string) {
	render.JSON(w, r, ResponseAuthUser{Response: response.OK(),
		Name: name, ID: userID, Role: role, Token: token})
}
