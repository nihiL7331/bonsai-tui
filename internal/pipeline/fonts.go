package pipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type MSDFJson struct {
	Atlas struct {
		Width         int     `json:"width"`
		Height        int     `json:"height"`
		DistanceRange float64 `json:"distanceRange"`
	} `json:"atlas"`
	Glyphs []struct {
		Unicode     int     `json:"unicode"`
		Advance     float64 `json:"advance"`
		PlaneBounds struct {
			Left, Bottom, Right, Top float64
		} `json:"planeBounds"`
		AtlasBounds struct {
			Left, Bottom, Right, Top float64
		} `json:"atlasBounds"`
	} `json:"glyphs"`
}

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

func GenerateFontAtlas(msdfExe string, ttfPath string, outDir string) (*MSDFJson, bool, uint8, error) {
	filename := filepath.Base(ttfPath)
	cleanName, isPixel, nativeSize := parseFontFilename(filename)

	outPng := filepath.Join(outDir, cleanName+".png")
	outJson := filepath.Join(outDir, cleanName+".json")

	args := []string{
		"-font", ttfPath,
		"-format", "png",
		"-imageout", outPng,
		"-json", outJson,
		"-yorigin", "top",
	}

	if isPixel {
		args = append(args,
			"-type", "hardmask",
			"-size", fmt.Sprintf("%d", nativeSize),
		)
	} else {
		args = append(args,
			"-type", "mtsdf",
			"-size", "32",
			"-pxrange", "4",
		)
	}

	cmd := exec.Command(msdfExe, args...)

	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, false, 0, fmt.Errorf("msdf-atlas-gen failed: %s\nOutput: %s", err, string(output))
	}

	jsonBytes, err := os.ReadFile(outJson)
	if err != nil {
		return nil, false, 0, fmt.Errorf("Failed to read generated JSON: %w", err)
	}

	var data MSDFJson
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		return nil, false, 0, fmt.Errorf("Failed to parse MSDF JSON: %w", err)
	}

	return &data, isPixel, nativeSize, nil
}
