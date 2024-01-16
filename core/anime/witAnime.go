package anime

import (
	"fmt"
	"goDown/types"
	"goDown/utils"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

type WitAnime types.SiteData

func NewWitAnime(logger *zap.SugaredLogger) *WitAnime {
	return &WitAnime{
		Name:   "witanime",
		Url:    "https://witanime.pics/",
		Logger: logger,
		SiteUrls: types.SiteUrls{
			IndexPageUrl:          "https://witanime.pics/",
			EpisodeArchivePageUrl: "https://witanime.pics/episode/",
			MediaArchivePageUrl:   "https://witanime.pics/%d9%82%d8%a7%d8%a6%d9%85%d8%a9-%d8%a7%d9%84%d8%a7%d9%86%d9%85%d9%8a/",

			// IndexPageUrlRegex:          `^(https?:\/\/)witanime\.pics\/?$`,
			// EpisodeArchivePageUrlRegex: `^(https?:\/\/)witanime\.pics\/episode\/?(?:\/page\/\d+\/?)?$`,
			// MediaArchivePageUrlRegex:   `^(https?:\/\/)?witanime\.pics/%[0-9a-fA-F]{2}(%[0-9a-fA-F]{2})*(/page/\d+)?/?$`,
			SearchUrl: "https://witanime.pics/?search_param=animes&s=[REPLACE]",
			MediaPageRegex:   `^(https?:\/\/)witanime\.pics/anime/[a-zA-Z0-9\-]+$`,
			EpisodePageRegex: `^(https?:\/\/)witanime\.pics/episode/[^/]+/$`,
		},
	}
}

func (w *WitAnime) GetSiteData() *types.SiteData {
	a := types.SiteData(*w)
	return &a
}

func (w *WitAnime) ScrapeEpisodeArchive() ([]string, error) {
	arr := []string{}
	w.Logger.Infof("started scraping the episodes archive")
	next, err := w.ScrapeEpisodeArchivePage(&arr, w.EpisodeArchivePageUrl)
	// next, err := w.ScrapeEpisodeArchivePage(&arr, "https://witanime.pics/episode/page/395/")
	if err != nil {
		w.Logger.Errorw("error while scraping episode archive", "next page url", next, "error", err)
		return nil, err
	}
	for {
		w.Logger.Infof("started scraping page: %v", next)
		nextUrl, err := w.ScrapeEpisodeArchivePage(&arr, next)
		if err != nil {
			return nil, err
		}
		if nextUrl == "" || nextUrl == next {
			break
		}
		next = nextUrl
	}
	return arr, nil
}
func (w *WitAnime) ScrapeEpisodeArchivePage(EpisodeUrlsPtr *[]string, pageUrl string) (string, error) {
	EpisodeUrls := *EpisodeUrlsPtr
	doc, err := utils.GetDocFromUrl(w.Logger, http.MethodGet, pageUrl)
	if err != nil {
		return "", err
	}
	doc.Find(".episodes-card-title a").Each(func(i int, s *goquery.Selection) {
		epUrlDecoded := w.findEpisodeUrl(s)
		if epUrlDecoded == "" {
			return
		}
		w.Logger.Infof("scraped episode link: %v", epUrlDecoded)

		EpisodeUrls = append(EpisodeUrls, epUrlDecoded)
	})
	nextPageUrl, exist := doc.Find(".pagination li:last-child a").First().Attr("href")
	if !exist {
		return "", nil
	}
	if len(EpisodeUrls) == 0 {
		w.Logger.Warnf("could not find any episodes on page: %v", pageUrl)
	}
	return nextPageUrl, nil
}

func (w *WitAnime) findEpisodeUrl(s *goquery.Selection) string {
	epData, exist := s.Attr("onclick")
	if !exist {
		w.Logger.Warnw(
			"could not find the episode card",
			"selector", ".episodes-card-title a",
			"tag", s,
		)
		return ""
	}
	epUrlEncoded, exist := utils.ExtractBetweenSingleQuotes(epData)
	if !exist {
		w.Logger.Warnw("could not get the text from episode onclick function", "data", epData)
		return ""
	}
	epUrlDecoded, err := utils.DecodeAtob(epUrlEncoded)
	if err != nil {
		w.Logger.Warnw("could not decode episode url", "encodedUrl", epUrlEncoded)
		return ""
	}
	return epUrlDecoded
}
func (w *WitAnime) ScrapeLatestEpisodes() ([]string, error) {
	latestEpisodes := []string{}
	doc, err := utils.GetDocFromUrl(w.Logger, http.MethodGet, w.IndexPageUrl)
	if err != nil {
		return nil, err
	}
	episodesSelector := ".episodes-list-content:first-child .episodes-card-title a"
	doc.Find(episodesSelector).Each(
		func(i int, s *goquery.Selection) {
			epLink := s.AttrOr("href", "")
			if epLink == "" {
				w.Logger.Warnw(
					"could not find the new episode",
					"selector", episodesSelector,
					"tag", s,
				)
			}
			w.Logger.Infof("found new episode: %v", epLink)
			latestEpisodes = append(latestEpisodes, epLink)
		})
	if len(latestEpisodes) == 0 {
		w.Logger.Warnw(
			"could not find any episodes",
			"selector", episodesSelector,
		)
		return nil, fmt.Errorf("could not find any episodes")
	}
	return latestEpisodes, nil
}

func (w *WitAnime) ScrapeMediaArchive() ([]string, error) {
	arr := []string{}
	w.Logger.Infof("started scraping the media archive")
	next, err := w.ScrapeMediaArchivePage(&arr, w.MediaArchivePageUrl)
	// next, err := w.ScrapeMediaArchivePage(&arr, "https://witanime.pics/%D9%82%D8%A7%D8%A6%D9%85%D8%A9-%D8%A7%D9%84%D8%A7%D9%86%D9%85%D9%8A/page/50/")
	if err != nil {
		w.Logger.Errorw("error while scraping episode archive", "next page url", next, "error", err)
		return nil, err
	}
	for {
		w.Logger.Infof("started scraping page: %v", next)
		nextUrl, err := w.ScrapeMediaArchivePage(&arr, next)
		if err != nil {
			return nil, err
		}
		if nextUrl == "" {
			break
		}
		next = nextUrl
	}
	return arr, nil
}

func (w *WitAnime) ScrapeMediaArchivePage(MediaUrlsPtr *[]string, pageUrl string) (string, error) {
	MediaUrls := *MediaUrlsPtr
	doc, err := utils.GetDocFromUrl(w.Logger, http.MethodGet, pageUrl)
	if err != nil {
		return "", err
	}
	doc.Find(".anime-card-title a").Each(func(i int, s *goquery.Selection) {
		mediaUrl := s.AttrOr("href", "")
		if mediaUrl == "" {
			w.Logger.Warnw("could not find the media url", "selector", ".anime-card-title a", "selection", s)
			return
		}
		w.Logger.Infof("scraped episode link: %v", mediaUrl)
		MediaUrls = append(MediaUrls, mediaUrl)
	})
	nextPageUrl, exist := doc.Find(".pagination a.next").First().Attr("href")
	if !exist {
		return "", nil
	}
	if len(MediaUrls) == 0 {
		w.Logger.Warnf("could not find any media on page: %v", pageUrl)
	}
	return nextPageUrl, nil
}

func (w *WitAnime) ScrapeLatestMedia() ([]string, error) {
	latestMedia := []string{}
	doc, err := utils.GetDocFromUrl(w.Logger, http.MethodGet, w.IndexPageUrl)
	if err != nil {
		return nil, err
	}

	selectors := []string{".owl-carousel", ".anime-card-title a"}
	doc.Find(selectors[0]).Last().Find(selectors[1]).Each(
		func(i int, s *goquery.Selection) {
			mediaLink := s.AttrOr("href", "")
			if mediaLink == "" {
				w.Logger.Warnw(
					"could not find the new media",
					"selector", selectors,
					"tag", s,
				)
			}
			w.Logger.Infof("found new media: %v", mediaLink)
			latestMedia = append(latestMedia, mediaLink)
		})
	if len(latestMedia) == 0 {
		w.Logger.Warnw(
			"could not find any media",
			"selector", selectors,
		)
		return nil, fmt.Errorf("could not find any media")
	}
	return latestMedia, nil
}

func (w *WitAnime) SearchForMedia(query string) ([]string, error) {
	MediaUrls := []string{}
	escapedquery := url.QueryEscape(query)
	searchUrl := strings.ReplaceAll(w.SearchUrl, "[REPLACE]", escapedquery)
	doc, err := utils.GetDocFromUrl(w.Logger, http.MethodGet, searchUrl)
	if err != nil {
		return nil, err
	}
	
	doc.Find(".anime-card-title a").Each(func(i int, s *goquery.Selection) {
		mediaUrl := s.AttrOr("href", "")
		if mediaUrl == "" {
			w.Logger.Warnw("could not find the media url", "selector", ".anime-card-title a", "selection", s)
			return
		}
		w.Logger.Infof("scraped episode link: %v", mediaUrl)
		MediaUrls = append(MediaUrls, mediaUrl)
	})
	if len(MediaUrls) == 0 {
		w.Logger.Errorw("could not find any media", "query", query)
		return MediaUrls, fmt.Errorf("could not find any media")
	}
	return MediaUrls, nil
}

func (w *WitAnime) ScrapeEpisodePage() {

}
func (w *WitAnime) ScrapeMediaPage(mediaUrl string) (*types.Media, error) {
	doc, err := utils.GetDocFromUrl(w.Logger, http.MethodGet, mediaUrl)
	if err != nil {
		return nil, err
	}
	title := strings.TrimSpace(doc.Find("h1").Text())
	m := types.Media{
		BaseInfo: types.BaseInfo{
			Name: title,
		},
	}
	// fmt.Printf("title: %v\n")
	return &m, nil
}
func (w *WitAnime) DownloadEpisode() {

}
func (w *WitAnime) DownloadMedia() {

}
