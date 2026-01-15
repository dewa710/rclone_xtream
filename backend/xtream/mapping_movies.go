package xtream

import (
	"fmt"
	"strings"
)

func movieFilename(v VodStream) string {
	n := sanitize(v.Name)
	if v.Year != "" {
		return fmt.Sprintf("%s (%s).%s", n, v.Year, v.ContainerExt)
	}
	return fmt.Sprintf("%s.%s", n, v.ContainerExt)
}

func sanitize(s string) string {
	return strings.NewReplacer("/", "-", "\\", "-", ":", "").Replace(s)
}
