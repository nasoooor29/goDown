package novel

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

type NovelBuddyIo struct {
	types.SiteData
	LatestMediaPageUrl string
}

func NewNovelBuddyIo(logger *zap.SugaredLogger) *NovelBuddyIo {
	d := types.SiteData{
		Name:   "NovelBuddy",
		Url:    "https://novelbuddy.io",
		Logger: logger,
		SiteUrls: types.SiteUrls{
			IndexPageUrl: "https://novelbuddy.io/home",
			// EpisodeArchivePageUrl: "https://novelbuddy.io/episode/", // website not support indivisiual episode archive
			MediaArchivePageUrl: "https://novelbuddy.io/az-list",
			SearchUrl:           "https://novelbuddy.io/search?q=[REPLACE]",
			MediaPageRegex:      `^(https?:\/\/)novelbuddy\.io/novel/[a-zA-Z0-9\-]+/$`,
			EpisodePageRegex:    `^(https?:\/\/)novelbuddy\.io/novel/[a-zA-Z0-9\-]+/[a-zA-Z0-9\-]+`,
		},
	}
	return &NovelBuddyIo{
		SiteData:           d,
		LatestMediaPageUrl: "https://novelbuddy.io/latest",
	}
}

func (w *NovelBuddyIo) GetSiteData() *types.SiteData {
	a := types.SiteData(w.SiteData)
	return &a
}

func (w *NovelBuddyIo) ScrapeEpisodeArchive() ([]string, error) {
	siteData := w.GetSiteData()
	return nil, types.FeatureIsNotSupported(siteData, types.ScrapeEpisodeArchive)
}
func (w *NovelBuddyIo) ScrapeEpisodeArchivePage(EpisodeUrlsPtr *[]string, pageUrl string) (string, error) {
	siteData := w.GetSiteData()
	return "", types.FeatureIsNotSupported(siteData, types.ScrapeEpisodeArchivePage)
}

func (w *NovelBuddyIo) ScrapeLatestEpisodes() ([]string, error) {
	siteData := w.GetSiteData()
	return nil, types.FeatureIsNotSupported(siteData, types.ScrapeLatestEpisodes)
}

func (w *NovelBuddyIo) ScrapeMediaArchive() ([]string, error) {
	arr := []string{}
	w.Logger.Infof("started scraping the media archive")
	next, err := w.ScrapeMediaArchivePage(&arr, w.MediaArchivePageUrl)
	// next, err := w.ScrapeMediaArchivePage(&arr, "https://novelbuddy.io/az-list?page=155")
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

func (w *NovelBuddyIo) findMediaLinks(doc *goquery.Document, MediaUrls *[]string) error {
	doc.Find(".book-item .title a:first-child").Each(func(i int, s *goquery.Selection) {
		mediaUrl := s.AttrOr("href", "")
		if mediaUrl == "" {
			w.Logger.Warnw("could not find the media url", "selector", ".book-item .title a:first-child", "selection", s)
			return
		}
		extendedMediaUrl := w.Url + mediaUrl
		w.Logger.Infof("scraped episode link: %v", extendedMediaUrl)
		*MediaUrls = append(*MediaUrls, extendedMediaUrl)
	})
	if len(*MediaUrls) == 0 {
		return fmt.Errorf("could not find any media")
	}
	return nil
}

func (w *NovelBuddyIo) ScrapeMediaArchivePage(MediaUrlsPtr *[]string, pageUrl string) (string, error) {
	doc, err := utils.GetDocFromUrl(w.Logger, http.MethodGet, pageUrl)
	if err != nil {
		return "", err
	}
	err = w.findMediaLinks(doc, MediaUrlsPtr)
	if err != nil {
		return "", err
	}
	nextPageUrl, exist := doc.Find("a.btn.link.active + a.btn.link").First().Attr("href")
	if !exist {
		return "", nil
	}
	return w.Url + nextPageUrl, nil
}

func (w *NovelBuddyIo) ScrapeLatestMedia() ([]string, error) {
	latestMedia := []string{}
	doc, err := utils.GetDocFromUrl(w.Logger, http.MethodGet, w.LatestMediaPageUrl)
	if err != nil {
		return nil, err
	}
	err = w.findMediaLinks(doc, &latestMedia)
	if err != nil {
		return nil, err
	}

	return latestMedia, nil
}

func (w *NovelBuddyIo) SearchForMedia(query string) ([]string, error) {
	MediaUrls := []string{}
	escapedquery := url.QueryEscape(query)
	searchUrl := strings.ReplaceAll(w.SearchUrl, "[REPLACE]", escapedquery)
	doc, err := utils.GetDocFromUrl(w.Logger, http.MethodGet, searchUrl)
	if err != nil {
		return nil, err
	}
	err = w.findMediaLinks(doc, &MediaUrls)
	if err != nil {
		w.Logger.Errorw("could not find any media", "query", query)
		return nil, err
	}
	return MediaUrls, nil
}

func (w *NovelBuddyIo) ScrapeMediaPage(mediaUrl string) (*types.Media, error) {
	doc, err := utils.GetDocFromUrl(w.Logger, http.MethodGet, mediaUrl)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(doc.Find("div.name h1").Text())
	if name == "" {
		w.Logger.Errorw("could not find the title", "mediaUrl", mediaUrl)
		return nil, fmt.Errorf("could not find the title")
	}
	media := types.NewMedia(types.Novel, name)

	if summary := strings.TrimSpace(doc.Find("div.summary p.content").Text()); summary == "" {
		w.Logger.Warnw("could not find the anime summary", "mediaUrl", mediaUrl)
	} else {
		media.Summary = summary
	}

	if img := strings.TrimSpace(doc.Find("#cover img:first-child").AttrOr("data-src", "")); img == "" {
		w.Logger.Warnw("could not find the anime thumbnail", "mediaUrl", mediaUrl)
	} else {
		media.ThumnailUrl = img
	}
	media.MetaData = w.getMediaMetaData(doc)
	media.Tags = w.getMediaTags(doc)
	media.Episodes = w.getMediaPageEpisodes(doc)
	return media, nil
}

func (w *NovelBuddyIo) getMediaMetaData(doc *goquery.Document) map[string]string {
	metaData := map[string]string{}
	doc.Find("div.detail div.meta p").Each(func(i int, s *goquery.Selection) {
		k := strings.TrimSpace(strings.ReplaceAll(s.Find("strong").First().Text(), ":", ""))
		s.Find("strong").First().Remove()
		v := strings.TrimSpace(s.Text())
		if k == "Genres" {
			return
		}
		if k != "" {
			metaData[k] = v
		}
	})
	return metaData
}

func (w *NovelBuddyIo) getMediaTags(doc *goquery.Document) []string {
	tags := []string{}
	doc.Find("div.detail div.meta p").Each(func(i int, s *goquery.Selection) {
		metaName := strings.TrimSpace(strings.ReplaceAll(s.Find("strong").First().Text(), ":", ""))
		if metaName != "Genres" {
			return
		}
		s.Find("a").Each(func(i int, s *goquery.Selection) {
			tag := strings.TrimSpace(strings.ReplaceAll(s.Text(), ",", ""))
			if tag != "" {
				tags = append(tags, tag)
			}
		})
	})
	return tags
}

func (w *NovelBuddyIo) getMediaPageEpisodes(doc *goquery.Document) []types.Episode {
	episodes := []types.Episode{}
	doc.Find("ul#chapter-list a").Each(func(i int, s *goquery.Selection) {
		epName := s.Find(".chapter-title").First().Text()
		epUrl, exist := s.Attr("href")
		if !exist {
			w.Logger.Warnf("could not find episode url, site: %v", w.Url)
			return
		}
		if epName == "" {
			w.Logger.Errorw("could not find episode name", "selection", s)
		}
		ep := types.NewEpisode(types.Novel, epName)
		ep.Downloadable = types.Downloadable{
			FileName: fmt.Sprintf("%v.%v", utils.PadNumber(5, i), "txt"),
			Index:    i,
			URL:      w.Url + epUrl,
		}
		episodes = append(episodes, *ep)
	})
	return episodes
}

func (w *NovelBuddyIo) ScrapeEpisodePage(epData *types.Episode) error {
	cl := utils.NewClient()
	h, err := cl.Head(epData.URL)
	if err != nil {
		w.Logger.Errorw("could not get episode page", "err", err)
		return err
	}
	if h.StatusCode != http.StatusOK {
		return fmt.Errorf("response is not ok, res: %v", h.Status)
	}
	return nil
}

func (w *NovelBuddyIo) DownloadEpisode() {

}
func (w *NovelBuddyIo) DownloadMedia() {

}
