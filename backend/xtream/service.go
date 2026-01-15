package xtream

import (
	"context"
	"fmt"
	"time"

	"github.com/rclone/rclone/fs"
)

type Service struct {
	fs        *Fs
	api       *Client
	cache     *Cache
	movieMap  map[string]string
	seriesMap map[string]string
}

func NewService(api *Client) *Service {
	return &Service{
		api:       api,
		cache:     NewCache(10 * time.Minute),
		movieMap:  make(map[string]string), // FIX
		seriesMap: make(map[string]string), // FIX
	}
}

func (s *Service) Load(ctx context.Context) error {
	// Movies
	fs.Infof(s.fs, "Load VodStreams")
	streams, err := s.VodStreams()
	if err != nil {
		return err
	}

	for _, m := range streams {
		ext := m.ContainerExt
		if ext == "" {
			ext = "mp4"
		}
		remote := m.Name + "." + ext
		s.movieMap[remote] = s.MovieURL(m.StreamID, ext)
	}
	return nil
	fs.Infof(s.fs, "Load Series")

	// Series
	series, err := s.Series()
	if err != nil {
		return err
	}

	for _, srs := range series {
		info, err := s.SeriesInfo(srs.SeriesID)
		if err != nil {
			continue
		}

		for _, eps := range info.Episodes {
			for _, ep := range eps {
				ext := ep.ContainerExt
				if ext == "" {
					ext = "mp4"
				}
				remote := ep.Title + "." + ext
				s.seriesMap[remote] = s.SeriesURL(ep.ID, ext)
			}
		}
	}
	fs.Infof(s.fs, "Load abgeschlossen: %d Movies, %d Series", len(s.movieMap), len(s.seriesMap))
	return nil
}

func (s *Service) VodCategories() ([]VodCategory, error) {
	if v, ok := s.cache.Get("vod_cat"); ok {
		return v.([]VodCategory), nil
	}
	var out []VodCategory
	if err := s.api.call("get_vod_categories", &out); err != nil {
		return nil, err
	}
	s.cache.Set("vod_cat", out)
	return out, nil
}

func (s *Service) VodStreams() ([]VodStream, error) {
	if v, ok := s.cache.Get("vod_streams"); ok {
		return v.([]VodStream), nil
	}
	var out []VodStream
	if err := s.api.call("get_vod_streams", &out); err != nil {
		return nil, err
	}
	s.cache.Set("vod_streams", out)
	return out, nil
}

func (s *Service) Series() ([]Series, error) {
	if v, ok := s.cache.Get("series"); ok {
		return v.([]Series), nil
	}
	var out []Series
	if err := s.api.call("get_series", &out); err != nil {
		return nil, err
	}
	s.cache.Set("series", out)
	return out, nil
}

func (s *Service) SeriesInfo(id int) (*SeriesInfo, error) {
	key := fmt.Sprintf("series_%d", id)
	if v, ok := s.cache.Get(key); ok {
		return v.(*SeriesInfo), nil
	}
	var out SeriesInfo
	if err := s.api.call(fmt.Sprintf("get_series_info&series_id=%d", id), &out); err != nil {
		return nil, err
	}
	s.cache.Set(key, &out)
	return &out, nil
}

func (s *Service) MovieURL(id int, ext string) string {
	return fmt.Sprintf("%s/movie/%s/%s/%d.%s",
		s.api.host, s.api.user, s.api.pass, id, ext)
}

func (s *Service) SeriesURL(id int, ext string) string {
	return fmt.Sprintf("%s/series/%s/%s/%d.%s",
		s.api.host, s.api.user, s.api.pass, id, ext)
}

func (s *Service) ResolveURL(remote string) (string, error) {
	if url, ok := s.movieMap[remote]; ok {
		return url, nil
	}
	if url, ok := s.seriesMap[remote]; ok {
		return url, nil
	}
	return "", fs.ErrorObjectNotFound
}

func (s *Service) NewObject(ctx context.Context, remote string) (fs.Object, error) {
	url, err := s.ResolveURL(remote)
	if err != nil {
		return nil, err
	}
	return &Object{
		fs:   s.fs,
		name: remote,
		url:  url,
	}, nil
}

func (s *Service) UserAgent() string {
	if s.api == nil {
		return ""
	}
	return s.api.userAgent
}
