package xtream

import (
	"context"
	"strings"
	"time"

	"github.com/rclone/rclone/fs"
)

func (f *Fs) List(ctx context.Context, dir string) (fs.DirEntries, error) {
	p := strings.Split(dir, "/")

	switch len(p) {
	case 1:
		if dir == "" {
			return fs.DirEntries{
				fs.NewDir("Movies", time.Time{}),
				fs.NewDir("Series", time.Time{}),
			}, nil
		}
		if dir == "Movies" {
			return f.listMovieCategories()
		}
		if dir == "Series" {
			return f.listSeriesRoot()
		}
	case 2:
		if p[0] == "Movies" {
			return f.listMoviesInCategory(p[1])
		}
		if p[0] == "Series" {
			return f.listSeriesSeasons(p[1])
		}
	case 3:
		if p[0] == "Series" {
			return f.listSeriesEpisodes(p[1], p[2])
		}
	}
	return nil, fs.ErrorDirNotFound
}
