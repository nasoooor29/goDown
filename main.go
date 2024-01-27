package main

import (
	"fmt"
	"goDown/core"
	"goDown/utils"
)

func main() {
	url := "https://novelbuddy.io/novel/my-trillion-dollar-assets-is-exposed-by-my-wifes-bragging"
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
	// siteData := site.GetSiteData()
	// arr, err := site.SearchForMedia("FUCK")
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// 	return
	// }

	m, err := site.ScrapeMediaPage(url)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	site.ScrapeEpisodePage(&m.Episodes[0])
	
	// fmt.Printf("m: %v\n", m)
	// a := utils.ReturnAsJson(m)
	// fmt.Printf("a: %v\n", a)
	// // site.ScrapeMediaPage(url)
	// // site.SearchForMedia("one piece")
	// // site.ScrapeMediaPage(url)
	// epUrl := "https://witanime.pics/episode/tsuki-ga-michibiku-isekai-douchuu-2nd-season-%d8%a7%d9%84%d8%ad%d9%84%d9%82%d8%a9-1/"
	// site.ScrapeEpisodePage(&types.Episode{
	// 	BaseInfo: *types.NewBaseInfo("EPEPEPE"),
	// 	Downloadable: types.Downloadable{
	// 		URL: epUrl,
	// 	},
	// })

}
