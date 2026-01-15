package xtream

import (
	"context"
	"io"
	"time"

	"github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/config/configmap"
	"github.com/rclone/rclone/fs/config/configstruct"
	"github.com/rclone/rclone/fs/hash"
)

const (
	defaultUserAgent = "rclone-xtream/1.0"
	defaultRateLimit = 5
)

// Options defines the configuration for this backend
type Options struct {
	Host         string        `config:"host"`
	Username     string        `config:"username"`
	Password     string        `config:"password"`
	UserAgent    string        `config:"user_agent"`
	RateLimit    int           `config:"rate_limit"`
	RetryCount   int           `config:"retry_count"`
	RetryMinWait time.Duration `config:"retry_min_wait"`
	RetryMaxWait time.Duration `config:"retry_max_wait"`
}

func init() {
	fs.Register(&fs.RegInfo{
		Name:        "xtream",
		Description: "Xtream Codes VOD backend (Movies & Series)",
		NewFs:       NewFs,
		Options: []fs.Option{
			{Name: "host", Help: "Xtream server URL (http://host:port)", Required: true},
			{Name: "username", Help: "Xtream username", Required: true},
			{Name: "password", Help: "Xtream password", Required: true},
			{Name: "user_agent", Help: "User-Agent for API and stream requests", Default: defaultUserAgent, Advanced: true},
			{Name: "rate_limit", Help: "Max API calls per second (0 = unlimited)", Default: defaultRateLimit, Advanced: true},
			{Name: "retry_count", Help: "Retries on transient errors", Default: "5", Advanced: true},
			{Name: "retry_min_wait", Help: "Min backoff", Default: "500ms", Advanced: true},
			{Name: "retry_max_wait", Help: "Max backoff", Default: "10s", Advanced: true},
		},
	})
}

// Fs represents the Xtream VOD filesystem
type Fs struct {
	name string
	root string
	opt  Options
	svc  *Service
}

// NewFs constructs an Fs from the path
func NewFs(ctx context.Context, name, root string, m configmap.Mapper) (fs.Fs, error) {
	opt := new(Options)
	if err := configstruct.Set(m, opt); err != nil {
		return nil, err
	}

	api := NewClient(
		opt.Host,
		opt.Username,
		opt.Password,
		opt.UserAgent,
		opt.RateLimit,
		opt.RetryCount,
		opt.RetryMinWait,
		opt.RetryMaxWait,
	)

	svc := NewService(api)
	if err := svc.Load(ctx); err != nil {
		return nil, err
	}

	return &Fs{
		name: name,
		root: root,
		opt:  *opt,
		svc:  svc,
	}, nil
}

func (f *Fs) Name() string             { return f.name }
func (f *Fs) Root() string             { return f.root }
func (f *Fs) String() string           { return "Xtream Codes VOD filesystem" }
func (f *Fs) Precision() time.Duration { return fs.ModTimeNotSupported }
func (f *Fs) Hashes() hash.Set         { return hash.Set(hash.None) }
func (f *Fs) Features() *fs.Features   { return &fs.Features{} }
func (f *Fs) IsReadOnly() bool         { return true }
func (f *Fs) SupportsStreaming() bool  { return true }

func (f *Fs) NewObject(ctx context.Context, remote string) (fs.Object, error) {
	return f.svc.NewObject(ctx, remote)
}

func (f *Fs) Put(ctx context.Context, in io.Reader, src fs.ObjectInfo, options ...fs.OpenOption) (fs.Object, error) {
	return nil, fs.ErrorPermissionDenied
}

func (f *Fs) Mkdir(ctx context.Context, dir string) error { return fs.ErrorPermissionDenied }
func (f *Fs) Rmdir(ctx context.Context, dir string) error { return fs.ErrorPermissionDenied }
