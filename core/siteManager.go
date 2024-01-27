package core

import (
	"fmt"
	"goDown/core/anime"
	"goDown/core/novel"
	"goDown/types"
	"goDown/utils"

	"go.uber.org/zap"
)

type SiteManager struct {
	Sites []types.SiteInterface
}

func NewSiteManager(logger *zap.SugaredLogger) *SiteManager {
	logger.Info("created site manager")
	AvaliableSites := []types.SiteInterface{
		anime.NewWitAnime(logger),
		novel.NewNovelBuddyIo(logger),
	}
	return &SiteManager{
		Sites: AvaliableSites,
	}
}

func (sm *SiteManager) FindSite(logger *zap.SugaredLogger, url string) (types.SiteInterface, error) {
	for _, site := range sm.Sites {
		ok, err := utils.MatchUrlHosts(url, site.GetSiteData().Url)
		if err != nil {
			logger.Warnw("url matching error",
				"user url", url,
				// "Site url", site.GetSiteData().Url,
				"error", err,
			)
			// return nil, err
		}
		if ok {
			logger.Infow(fmt.Sprintf("site with %v was found", url), "site name", site.GetSiteData().Name, "url", url)
			return site, nil
		}
	}
	logger.Warnw(fmt.Sprintf("site with url: %v, not found", url), "url", url)
	return nil, fmt.Errorf("site not found")
}
