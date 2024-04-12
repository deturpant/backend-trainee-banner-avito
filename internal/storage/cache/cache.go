package cache

import (
	"backend-trainee-banner-avito/internal/entities"
	"strconv"
	"sync"
	"time"
)

type BannerCache struct {
	sync.RWMutex
	Banners map[string]cachedBanner
}

type cachedBanner struct {
	Banner    entities.Banner
	UpdatedAt time.Time
}

var cache BannerCache

func init() {
	cache = BannerCache{
		Banners: make(map[string]cachedBanner),
	}
}

func GetBannerFromCache(featureID, tagID int) (*entities.Banner, bool) {
	cache.RLock()
	defer cache.RUnlock()

	key := generateCacheKey(featureID, tagID)
	cached, found := cache.Banners[key]
	if !found {
		return nil, false
	}

	if time.Since(cached.UpdatedAt) > 10*time.Minute {
		delete(cache.Banners, key)
		return nil, false
	}

	return &cached.Banner, true
}
func StoreBannerInCache(featureID, tagID int, banner entities.Banner) {
	cache.Lock()
	defer cache.Unlock()

	key := generateCacheKey(featureID, tagID)
	cache.Banners[key] = cachedBanner{
		Banner:    banner,
		UpdatedAt: time.Now(),
	}
}

func generateCacheKey(featureID, tagID int) string {
	return strconv.Itoa(featureID) + "-" + strconv.Itoa(tagID)
}
