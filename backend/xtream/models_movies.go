package xtream

type VodCategory struct {
	CategoryID string `json:"category_id"`
	Category   string `json:"category_name"`
}

type VodStream struct {
	StreamID     int    `json:"stream_id"`
	Name         string `json:"name"`
	CategoryID   string `json:"category_id"`
	ContainerExt string `json:"container_extension"`
	Year         string `json:"year"`
	FileSize     int64  `json:"file_size"`
}
