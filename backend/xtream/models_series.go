package xtream

type Series struct {
	SeriesID int    `json:"series_id"`
	Name     string `json:"name"`
}

type SeriesInfo struct {
	Episodes map[string][]SeriesEpisode `json:"episodes"`
}

type SeriesEpisode struct {
	ID           int    `json:"id"`
	EpisodeNum   int    `json:"episode_num"`
	Title        string `json:"title"`
	ContainerExt string `json:"container_extension"`
	FileSize     int64  `json:"file_size"`
}
