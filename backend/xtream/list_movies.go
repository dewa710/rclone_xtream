package xtream

import (
	"path"
	"sort"
	"strings"
	"time"

	"github.com/rclone/rclone/fs"
)

func (f *Fs) listMovieCategories() (fs.DirEntries, error) {
	seen := map[string]bool{}
	out := fs.DirEntries{}

	for remote := range f.svc.movieMap {
		// Movies/Kategorie/Film.mp4
		parts := strings.Split(remote, "/")
		if len(parts) < 3 {
			continue
		}
		cat := parts[1]
		if !seen[cat] {
			seen[cat] = true
			out = append(out, fs.NewDir(path.Join("Movies", cat), time.Time{}))
		}
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Remote() < out[j].Remote() })
	return out, nil
}

func (f *Fs) listMoviesInCategory(cat string) (fs.DirEntries, error) {
	out := fs.DirEntries{}

	for remote, url := range f.svc.movieMap {
		// Movies/Kategorie/Filename.ext
		parts := strings.Split(remote, "/")
		if len(parts) != 3 {
			continue
		}
		if parts[1] != cat {
			continue
		}

		out = append(out, &Object{
			fs:   f,
			name: remote,
			url:  url,
		})
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Remote() < out[j].Remote() })
	return out, nil
}
