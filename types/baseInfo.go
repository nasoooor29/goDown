package types

type DownloadStatus int

const (
	NotDownloaded DownloadStatus = iota
	Downloaded
)

type WatchStatus int

const (
	NotWatched WatchStatus = iota
	CurrentlyWatching
	DoneWatching
)

type BaseInfo struct {
	Id           string
	Name         string
	CreationTime int
	DlStatus     DownloadStatus
	WaStatus     WatchStatus
	MetaData     map[string]string
}

type DownloadableType int

const (
	Image DownloadableType = iota
	Video
	Text
)

type Downloadable struct {
	Path     string
	URL      string
	FileName string
	Type     DownloadableType
}
