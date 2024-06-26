package banners

import (
	"backend-trainee-banner-avito/internal/lib/api/response"
	"backend-trainee-banner-avito/internal/lib/logger/errMsg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

// @Summary Удалить баннер по его идентификатору
// @Description Удаляет баннер с указанным идентификатором
// @ID delete-banner
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Токен администратора" format:"Bearer <admin_token>"
// @Param id path int true "Идентификатор баннера"
// @Success 204 "Баннер успешно удален"
// @Failure 400 {object} banners.Response "Неверный идентификатор баннера"
// @Failure 500 {object} banners.Response "Ошибка при удалении баннера"
// @Router /banner/{id} [delete]
func NewDeleteBannerHandler(log *slog.Logger, bannerRepo Banners) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.banners.deleteBanner.New"
		log := log.With(
			slog.String("options", loggerOptions),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		idStr := chi.URLParam(r, "id") // Получаем id из URL-пути
		id, err := strconv.Atoi(idStr)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("Invalid banner ID"))
			return
		}

		err = bannerRepo.DeleteBannerByID(r.Context(), id)
		if err != nil {
			log.Error("Failed to delete banner", errMsg.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("Failed to delete banner"))
			return
		}
		log.Info("Banner deleted")
		render.Status(r, http.StatusNoContent)
	}
}
