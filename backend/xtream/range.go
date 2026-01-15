package xtream

import (
	"strconv"

	"github.com/rclone/rclone/fs"
)

func rangeFromOptions(opts []fs.OpenOption) string {
	var start, end int64 = -1, -1

	for _, o := range opts {
		switch v := o.(type) {
		case *fs.SeekOption:
			start = v.Offset
		case *fs.RangeOption:
			start = v.Start
			end = v.End
		}
	}

	if start >= 0 && end > start {
		return "bytes=" + itoa(start) + "-" + itoa(end)
	}
	if start >= 0 {
		return "bytes=" + itoa(start) + "-"
	}
	return ""
}

func itoa(i int64) string {
	return strconv.FormatInt(i, 10)
}
