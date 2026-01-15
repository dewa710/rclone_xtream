package xtream

import (
	"path"
	"sort"
	"strings"
	"time"

	"github.com/rclone/rclone/fs"
)

func (f *Fs) listSeriesRoot() (fs.DirEntries, error) {
	seen := map[string]bool{}
	out := fs.DirEntries{}

	for remote := range f.svc.seriesMap {
		parts := strings.Split(remote, "/")
		if len(parts) < 2 {
			continue
		}
		name := parts[1]
		if !seen[name] {
			seen[name] = true
			out = append(out, fs.NewDir(path.Join("Series", name), time.Time{}))
		}
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Remote() < out[j].Remote() })
	return out, nil
}

func (f *Fs) listSeriesSeasons(name string) (fs.DirEntries, error) {
	seen := map[string]bool{}
	out := fs.DirEntries{}

	for remote := range f.svc.seriesMap {
		parts := strings.Split(remote, "/")
		if len(parts) < 3 || parts[1] != name {
			continue
		}
		season := parts[2]
		if !seen[season] {
			seen[season] = true
			out = append(out, fs.NewDir(path.Join("Series", name, season), time.Time{}))
		}
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Remote() < out[j].Remote() })
	return out, nil
}

func (f *Fs) listSeriesEpisodes(name, season string) (fs.DirEntries, error) {
	out := fs.DirEntries{}

	for remote, url := range f.svc.seriesMap {
		parts := strings.Split(remote, "/")
		if len(parts) != 4 {
			continue
		}
		if parts[1] != name || parts[2] != season {
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
