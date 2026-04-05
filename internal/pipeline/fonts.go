package pipeline

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// returns: base name, is pixel, native size (0 if not pixel)
func parseFontFilename(filename string) (string, bool, uint8) {
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)

	parts := strings.Split(base, "_")
	if len(parts) > 1 {
		if size, err := strconv.ParseUint(parts[len(parts)-1], 10, 8); err == nil {
			cleanName := strings.Join(parts[:len(parts)-1], "_")
			return cleanName, true, uint8(size)
		}
	}

	return base, false, 0
}
