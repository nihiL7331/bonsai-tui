package pipeline

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
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

// creates a text file containing standard ASCII chars (32-126)
func EnsureDefaultCharset(outDir string) (string, error) {
	charsetPath := filepath.Join(outDir, "default_charset.txt")

	if _, err := os.Stat(charsetPath); err == nil {
		return charsetPath, nil
	}

	var sb strings.Builder
	for i := 32; i <= 126; i++ {
		sb.WriteByte(byte(i))
	}

	if err := os.WriteFile(charsetPath, []byte(sb.String()), 0o644); err != nil {
		return "", fmt.Errorf("Failed to write default charset: %w", err)
	}

	return charsetPath, nil
}

// iterates over the font dir, generates PNG atlases via
// msdf-atlas-gen and packs JSON into binary format
func buildFonts(fontSrcDir string, fontDataOutDir string, msdfExe string, logFn func(string, string)) ([]*HotReloadPayload, error) {
	var payloads []*HotReloadPayload

	if _, err := os.Stat(fontSrcDir); os.IsNotExist(err) {
		logFn("FONTS", "No fonts directory found, skipping font packing.")
		return payloads, nil
	}

	entries, err := os.ReadDir(fontSrcDir)
	if err != nil {
		return nil, fmt.Errorf("Failed to read '%s' directory: %w", fontSrcDir, err)
	}

	if err := os.MkdirAll(fontDataOutDir, 0o755); err != nil {
		return nil, fmt.Errorf("Failed to create '%s' directory: %w", fontDataOutDir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext == ".ttf" || ext == ".otf" {
			fullPath := filepath.Join(fontSrcDir, entry.Name())
			cleanName, _, _ := parseFontFilename(entry.Name())

			payload, err := packFont(msdfExe, fullPath, cleanName, fontDataOutDir, logFn)
			if err != nil {
				logFn("ERROR", fmt.Sprintf("Failed to process %s: %w", entry.Name(), err))
				continue
			}

			payloads = append(payloads, payload)
		}
	}

	return payloads, nil
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

func GenerateFontAtlas(msdfExe string, ttfPath string, outDir string, charsetPath string) (*MSDFJson, bool, uint8, error) {
	filename := filepath.Base(ttfPath)
	cleanName, isPixel, nativeSize := parseFontFilename(filename)

	outPng := filepath.Join(outDir, cleanName+".png")
	outJson := filepath.Join(outDir, cleanName+".json")

	args := []string{
		"-font", ttfPath,
		"-charset", charsetPath,
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

// checks the cache, generates atlas if needed,
// returns the hot-reload payload
func packFont(msdfExe string, fontPath string, cleanName string, outDir string, logFn func(string, string)) (*HotReloadPayload, error) {
	binPath := filepath.Join(outDir, cleanName+".bin")
	pngPath := filepath.Join(outDir, cleanName+".png")
	jsonPath := filepath.Join(outDir, cleanName+".json")

	if srcStat, err := os.Stat(fontPath); err == nil {
		binStat, errBin := os.Stat(binPath)
		pngStat, errPng := os.Stat(pngPath)

		if errBin == nil && errPng == nil {
			if binStat.ModTime().After(srcStat.ModTime()) && pngStat.ModTime().After(srcStat.ModTime()) {
				logFn("FONTS", fmt.Sprintf("Using cached font: %s", cleanName))

				binBytes, _ := os.ReadFile(binPath)
				pngBytes, _ := os.ReadFile(pngPath)
				return &HotReloadPayload{PngBytes: pngBytes, MetadataBin: binBytes}, nil
			}
		}
	}

	logFn("FONTS", fmt.Sprintf("Packing font: %s...", cleanName))

	charsetPath, err := EnsureDefaultCharset(outDir)
	if err != nil {
		return nil, fmt.Errorf("Failed to setup charset: %w", err)
	}

	msdfData, isPixel, nativeSize, err := GenerateFontAtlas(msdfExe, fontPath, outDir, charsetPath)
	if err != nil {
		return nil, fmt.Errorf("msdf-atlas-gen failed: %w", err)
	}

	binBytes, err := PackFontMetadata(msdfData, isPixel, nativeSize)
	if err != nil {
		return nil, fmt.Errorf("Failed to pack metadata: %w", err)
	}

	if err := os.WriteFile(binPath, binBytes, 0o644); err != nil {
		return nil, fmt.Errorf("Failed to write bin file: %w", err)
	}

	pngBytes, err := os.ReadFile(pngPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read generated PNG: %w", err)
	}

	os.Remove(jsonPath)

	return &HotReloadPayload{
		PngBytes:    pngBytes,
		MetadataBin: binBytes,
	}, nil
}

// takes the parsed MSDF JSON and converts it to .bin
func PackFontMetadata(data *MSDFJson, isPixel bool, nativeSize uint8) ([]byte, error) {
	buf := new(bytes.Buffer)

	isPixelByte := uint8(0)
	if isPixel {
		isPixelByte = 1
	}
	buf.WriteByte(isPixelByte)
	buf.WriteByte(nativeSize)

	scale := float64(32.0)
	if isPixel {
		scale = float64(nativeSize)
	}

	atlasWidth := float64(data.Atlas.Width)
	atlasHeight := float64(data.Atlas.Height)

	for _, glyph := range data.Glyphs {
		unicode := uint32(glyph.Unicode)

		var u0, v0, u1, v1 float32
		var w, h float32
		var xOffset, yOffset float32

		if glyph.AtlasBounds.Right != glyph.AtlasBounds.Left {
			u0 = float32(glyph.AtlasBounds.Left / atlasWidth)
			v0 = float32(glyph.AtlasBounds.Top / atlasHeight)
			u1 = float32(glyph.AtlasBounds.Right / atlasWidth)
			v1 = float32(glyph.AtlasBounds.Bottom / atlasHeight)

			w = float32(math.Abs(glyph.AtlasBounds.Right - glyph.AtlasBounds.Left))
			h = float32(math.Abs(glyph.AtlasBounds.Bottom - glyph.AtlasBounds.Top))

			xOffset = float32(glyph.PlaneBounds.Left * scale)
			yOffset = float32(glyph.PlaneBounds.Top * scale)
		}

		advance := float32(glyph.Advance * scale)

		binary.Write(buf, binary.LittleEndian, unicode)

		binary.Write(buf, binary.LittleEndian, u0)
		binary.Write(buf, binary.LittleEndian, v0)
		binary.Write(buf, binary.LittleEndian, u1)
		binary.Write(buf, binary.LittleEndian, v1)

		binary.Write(buf, binary.LittleEndian, w)
		binary.Write(buf, binary.LittleEndian, h)

		binary.Write(buf, binary.LittleEndian, xOffset)
		binary.Write(buf, binary.LittleEndian, yOffset)
		binary.Write(buf, binary.LittleEndian, advance)
	}

	return buf.Bytes(), nil
}
