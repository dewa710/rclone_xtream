package xtream

import "fmt"

func seasonDir(s string) string {
	return fmt.Sprintf("Season %02s", s)
}

func episodeFile(season string, e SeriesEpisode) string {
	return fmt.Sprintf(
		"S%02sE%02d - %s.%s",
		season, e.EpisodeNum, sanitize(e.Title), e.ContainerExt,
	)
}
