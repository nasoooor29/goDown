package types

import (
	"goDown/utils"
	"path"
	"time"
)

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
	Name         string
	Id           string
	CreationTime time.Time
	DlStatus     DownloadStatus
	WaStatus     WatchStatus
}

// This function will generate most of the data just include the name and metadata
func NewBaseInfo(n string) *BaseInfo {
	return &BaseInfo{
		Name:         n,
		Id:           utils.GenHashBasedOnTime(),
		CreationTime: time.Now(),
		DlStatus:     NotDownloaded,
		WaStatus:     NotWatched,
	}
}

type Downloadable struct {
	Index int
	URL   string
	FileName  string
}

type Quality string

const (
	FHD       Quality = "FHD"
	HD        Quality = "HD"
	SD        Quality = "SD"
	Undefined Quality = "Undefined"
)

func FindQuality(input string) Quality {
	cleanInput := utils.RemoveNonEnglishLetters(input)
	switch cleanInput {
	case "FHD":
		return FHD
	case "HD":
		return HD
	case "SD":
		return SD
	default:
		return Undefined
	}
}

const DefaultPath = "."

type Config struct {
	BasePath  string
	AnimePath string
	MangaPath string
	NovelPath string
}

func NewDefaultConfig() *Config {
	return &Config{
		BasePath:  DefaultPath,
		AnimePath: path.Join(DefaultPath, "Anime"),
		MangaPath: path.Join(DefaultPath, "Manga"),
		NovelPath: path.Join(DefaultPath, "Novel"),
	}
}
