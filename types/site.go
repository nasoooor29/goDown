package types

import "go.uber.org/zap"

type SiteInterface interface {
	GetSiteData() *SiteData
	SiteDownloader
	SiteScraper
}

type SiteDownloader interface {
	DownloadEpisode()
	DownloadMedia()
}

type SiteScraper interface {
	// ScrapeIndexPage() 		// this page is used to get the latest episodes of the website
	ScrapeEpisodeArchive() ([]string, error)
	ScrapeEpisodeArchivePage(*[]string, string) (string, error)
	ScrapeLatestEpisodes() ([]string, error)
	ScrapeEpisodePage(*Episode) error
	ScrapeMediaArchive() ([]string, error)
	ScrapeMediaArchivePage(*[]string, string) (string, error)
	ScrapeLatestMedia() ([]string, error)
	ScrapeMediaPage(string) (*Media, error)
	SearchForMedia(string) ([]string, error)
}

type SiteUrls struct {
	IndexPageUrl          string
	EpisodeArchivePageUrl string
	MediaArchivePageUrl   string

	IndexPageUrlRegex          string
	EpisodeArchivePageUrlRegex string
	MediaArchivePageUrlRegex   string

	MediaPageRegex   string
	EpisodePageRegex string

	SearchUrl string
}

func (su *SiteUrls) GetSiteUrls() *SiteUrls {
	return su
}

type SiteData struct {
	Name     string
	Url      string
	MetaData map[string]string
	Logger   *zap.SugaredLogger
	SiteUrls
}

func (s *SiteData) GetSiteData() *SiteData {
	return s
}

type SiteFeature string

const (
	ScrapeEpisodeArchive     SiteFeature = "Scrape Episode Archive"
	ScrapeEpisodeArchivePage SiteFeature = "Scrape Episode Archive Page"
	ScrapeLatestEpisodes     SiteFeature = "Scrape Latest Episodes"
	ScrapeEpisodePage        SiteFeature = "Scrape Episode Page"
	ScrapeMediaArchive       SiteFeature = "Scrape Media Archive"
	ScrapeMediaArchivePage   SiteFeature = "Scrape Media Archive Page"
	ScrapeLatestMedia        SiteFeature = "Scrape Latest Media"
	ScrapeMediaPage          SiteFeature = "Scrape Media Page"
)
