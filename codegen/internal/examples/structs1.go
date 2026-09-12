package codegen

import (
	hibiketorrent "seanime/internal/extension/hibike/torrent"
	"seanime/internal/media"
)

//type Struct1 struct {
//	Struct2
//}
//
//type Struct2 struct {
//	Text string `json:"text"`
//}

//type Struct3 []string

type Struct4 struct {
	Torrents    []hibiketorrent.AnimeTorrent `json:"torrents"`
	Destination string                       `json:"destination"`
	SmartSelect struct {
		Enabled               bool  `json:"enabled"`
		MissingEpisodeNumbers []int `json:"missingEpisodeNumbers"`
	} `json:"smartSelect"`
	Media *media.Anime `json:"media"`
}
