package main

import (
	"fmt"
	"goDown/core"
	"goDown/utils"
)

func main() {
	url := "https://witanime.pics/anime/tsuki-ga-michibiku-isekai-douchuu-2nd-season/"
	logger, err := utils.NewLogger(utils.BuildTypeDev)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	sm := core.NewSiteManager(logger)
	site, err := sm.FindSite(logger, url)
	if err != nil {
		return
	}
	// site.ScrapeMediaPage(url)
	site.SearchForMedia("isekai")

	// u := site.GetSiteData().EpisodeArchivePageUrl
	// site.ScrapeMediaArchive()

	// site.ScrapeLatestMedia()
	// site.ScrapeLatestEpisodes()

	// utils.SaveToHTML(url, "temp/wit2.html")
	// doc, err := utils.HtmlToDoc("temp/wit2.html")
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// 	return
	// }
	// title := doc.Find("h1").Text()
	// fmt.Printf("title: %v\n", title)
	// links := doc.Find(".episodes-card-title a")
	// fmt.Printf("links: %v\n", links)
	// links.Each(func(i int, s *goquery.Selection) {
	// 	fmt.Println(s.Attr("onclick"))
	// })

	// aa, err := utils.DecodeAtob("aHR0cHM6Ly93aXRhbmltZS5waWNzL2FuaW1lL2hpbWVzYW1hLWdvdW1vbi1uby1qaWthbi1kZXN1Lw==")
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// 	return
	// }
	// fmt.Printf("aa: %v\n", aa)
}
