package xtream

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/hash"
)

type Object struct {
	fs   *Fs
	name string
	url  string
	size int64
}

func (o *Object) Fs() fs.Info {
	return o.fs
}

func (o *Object) Remote() string {
	return o.name
}

func (o *Object) String() string {
	return o.name
}

func (o *Object) Size() int64 {
	if o.size > 0 {
		return o.size
	}
	return -1
}

func (o *Object) ModTime(ctx context.Context) time.Time {
	return time.Time{}
}

func (o *Object) Hash(ctx context.Context, ty hash.Type) (string, error) {
	return "", hash.ErrUnsupported
}

func (o *Object) Storable() bool {
	return true
}

func (o *Object) SetModTime(ctx context.Context, t time.Time) error {
	return fs.ErrorCantSetModTime
}

func (o *Object) Update(ctx context.Context, in io.Reader, src fs.ObjectInfo, options ...fs.OpenOption) error {
	return fs.ErrorPermissionDenied
}

func (o *Object) Remove(ctx context.Context) error {
	return fs.ErrorPermissionDenied
}

func (o *Object) Open(ctx context.Context, options ...fs.OpenOption) (io.ReadCloser, error) {
	var offset int64

	for _, opt := range options {
		if ro, ok := opt.(*fs.RangeOption); ok {
			offset = ro.Start
		}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", o.url, nil)
	if err != nil {
		return nil, err
	}

	if offset > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", offset))
	}

	// Xtream requires consistent User-Agent
	if ua := o.fs.svc.UserAgent(); ua != "" {
		req.Header.Set("User-Agent", ua)
	}

	client := &http.Client{
		Timeout: 0, // streaming: never timeout
	}

	var resp *http.Response

	for i := 0; i < 4; i++ {
		resp, err = client.Do(req)
		if err == nil && resp.StatusCode < 500 {
			break
		}
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 206 {
		return nil, fmt.Errorf("xtream returned %s", resp.Status)
	}

	if o.size == 0 && resp.ContentLength > 0 {
		o.size = resp.ContentLength
	}

	return resp.Body, nil
}
