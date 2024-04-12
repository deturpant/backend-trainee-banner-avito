package banners

import (
	"backend-trainee-banner-avito/internal/entities"
	"backend-trainee-banner-avito/internal/lib/api/response"
	"backend-trainee-banner-avito/internal/lib/logger/errMsg"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"time"
)

type Banners interface {
	CreateBanner(ctx context.Context, banner *entities.Banner) error
	FindBannerByFeatureTag(ctx context.Context, featureID, tagID int) (*entities.Banner, error)
	DeleteBannerByID(ctx context.Context, id int) error
}
type BannerTags interface {
	CreateBannerTag(ctx context.Context, bannerTag *entities.BannerTag) error
}
type RequestBanner struct {
	TagIDs    []int                  `json:"tag_ids" validate:"required"`
	FeatureID int                    `json:"feature_id" validate:"required"`
	Content   map[string]interface{} `json:"content" validate:"required"`
	IsActive  bool                   `json:"is_active" validate:"required"`
}

type ResponseBanner struct {
	response.Response
	ID        int                    `json:"banner_id"`
	TagIDs    []int                  `json:"tag_ids"`
	FeatureID int                    `json:"feature_id"`
	Content   map[string]interface{} `json:"content"`
	IsActive  bool                   `json:"is_active"`
}

func New(log *slog.Logger, bannerRepository Banners, bannerTagsRepository BannerTags) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.banners.createBanner.New"
		log = log.With(
			slog.String("options", loggerOptions),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req RequestBanner
		err := json.NewDecoder(r.Body).Decode(&req)
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
		banner := entities.Banner{
			TagIDs:    req.TagIDs,
			FeatureID: req.FeatureID,
			Content:   req.Content,
			IsActive:  req.IsActive,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = bannerRepository.CreateBanner(r.Context(), &banner)
		if err != nil {
			log.Error("Failed to create banner", errMsg.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("Failed to create banner"))
			return
		}
		log.Info("Banner added")
		for _, tagID := range req.TagIDs {
			bannerTag := entities.BannerTag{
				BannerID: banner.ID,
				TagID:    tagID,
			}
			err = bannerTagsRepository.CreateBannerTag(r.Context(), &bannerTag)
			if err != nil {
				render.Status(r, http.StatusBadRequest)
				log.Error("Failed to create banner tag", errMsg.Err(err))
				render.JSON(w, r, response.Error("Failed to create banner tag"))
				return
			}
		}
		log.Info("Banner-tags added")
		responseOK(w, r, banner)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, banner entities.Banner) {
	render.JSON(w, r, ResponseBanner{
		Response:  response.OK(),
		ID:        banner.ID,
		TagIDs:    banner.TagIDs,
		FeatureID: banner.FeatureID,
		Content:   banner.Content,
		IsActive:  banner.IsActive,
	})
}
