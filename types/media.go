package types

type MediaType int

const (
	Anime MediaType = iota
	Manga
	Novel
)

type Episode struct {
	Downloadable
	BaseInfo           //auto filled
	Type     MediaType // auto filled

	// Not so sure about this decision
	WatchMediaServers    map[Quality][]MediaServer
	DownloadMediaServers map[Quality][]MediaServer
}

type MediaServerRank int

const (
	Fast MediaServerRank = iota
	Medium
	Slow
	Unkown
)

type MediaServer struct {
	Name           string
	Url            string
	Rank           MediaServerRank
	EpisodeQuality Quality
}

func NewEpisode(t MediaType, name string) *Episode {
	return &Episode{
		BaseInfo: *NewBaseInfo(name),
		Type:     t,
	}
}

type Media struct {
	BaseInfo              // auto filled
	Type        MediaType // auto filled
	Summary     string
	ThumnailUrl string
	Tags        []string
	Episodes    []Episode
	MetaData    map[string]string
}

func NewMedia(t MediaType, name string) *Media {
	return &Media{
		BaseInfo: *NewBaseInfo(name),
		Type:     t,
	}
}

type MediaCollection struct {
	BaseInfo
	MediaContent []*Media
}
