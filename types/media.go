package types

type EpisodeType int

const (
	MangaChapter EpisodeType = iota
	AnimeEpisode
	NovelChapter
)

type Episode struct {
	BaseInfo
	Downloadable
	Type EpisodeType
}


func NewEpisode(t EpisodeType) *Episode {
	return &Episode{
		Type: t,
	}
}

type MediaType int

const (
	Anime MediaType = iota
	Manga
	Novel
)

type Media struct {
	BaseInfo
	Downloadable
	Episodes []*Episode
	Type     MediaType
}

func NewMedia(t MediaType) *Media {
	return &Media{
		Type: t,
	}
}

type MediaCollection struct {
	BaseInfo
	MediaContent []*Media
}
