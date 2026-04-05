package pipeline

import (
	"bonsai-tui/internal/config"
	"bonsai-tui/internal/pipeline/packer"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func PackAtlas(cfg config.Config, logFn func(string, string)) (*HotReloadPayload, error) {
	imagesDir := getImagesDir(cfg)
	atlasDir := getAtlasCacheDir(cfg)
	atlasPath := filepath.Join(atlasDir, "atlas.png")

	if !isDirNewer(imagesDir, atlasPath) {
		logFn("ATLAS", "Atlas is up to date. Skipping packing.")
		return nil, nil
	}

	files, err := getSortedImageFiles(imagesDir)
	if err != nil {
		return nil, fmt.Errorf("Failed to read images: %w", err)
	}

	odinOutDir := getOdinAtlasDir(cfg)
	binOutDir := getBinAtlasDir(cfg)

	if len(files) == 0 {
		_, err := generateSpriteMetadata(nil, 2048, 2048, odinOutDir, binOutDir)
		if err != nil {
			return nil, err
		}
		logFn("ATLAS", "No images to pack. Skipping packing.")
		return nil, nil
	}

	logFn("ATLAS", "Packing texture atlas...")

	pckr := packer.NewPacker(2048, 2048)

	placedSprites, err := processImages(files, pckr, logFn)
	if err != nil {
		return nil, fmt.Errorf("Failed to process images: %w", err)
	}

	pngBytes, err := packer.RenderAtlas(2048, 2048, placedSprites)
	if err != nil {
		return nil, fmt.Errorf("Failed to render atlas: %w", err)
	}

	if err := os.WriteFile(atlasPath, pngBytes, 0o644); err != nil {
		return nil, fmt.Errorf("Failed to write 'atlas.png': %w", err)
	}

	metadataBin, err := generateSpriteMetadata(
		placedSprites,
		2048,
		2048,
		odinOutDir,
		binOutDir,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate metadata: %w", err)
	}

	return &HotReloadPayload{
		PngBytes:    pngBytes,
		MetadataBin: metadataBin,
	}, nil
}

// grabs all .png files within 'dir'
// and sorts them alphabetically
func getSortedImageFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if strings.ToLower(filepath.Ext(path)) == ".png" {
			files = append(files, path)
		}
		return nil
	})

	sort.Strings(files)
	return files, err
}

func processImages(files []string, pckr *packer.MaxRectsPacker, logFn func(string, string)) ([]packer.PlacedSprite, error) {
	var placedSprites []packer.PlacedSprite
	padding := 2

	for _, file := range files {
		w, h, err := getImageDimensions(file)
		if err != nil {
			return nil, err
		}

		rect := pckr.Insert(w+padding, h+padding)
		if rect.W == 0 {
			return nil, fmt.Errorf("Atlas is full. Could not fit: %s", filepath.Base(file))
		}

		placedSprites = append(placedSprites, packer.PlacedSprite{
			Name: filepath.Base(file),
			Path: file,
			Rect: packer.Rect{
				X: rect.X + (padding / 2),
				Y: rect.Y + (padding / 2),
				W: w,
				H: h,
			},
		})
	}

	return placedSprites, nil
}

// reads header of a PNG file to get image size
func getImageDimensions(path string) (width int, height int, err error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	config, err := png.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}
	return config.Width, config.Height, nil
}

// removes dimension suffixes
func cleanKeySuffix(cleanKey string) string {
	lastUnderscore := strings.LastIndex(cleanKey, "_")
	if lastUnderscore != -1 {
		suffix := cleanKey[lastUnderscore+1:]

		xIdx := strings.Index(suffix, "x")
		if xIdx != -1 {
			left := suffix[:xIdx]
			right := suffix[xIdx+1:]

			_, errLeft := strconv.Atoi(left)
			_, errRight := strconv.Atoi(right)

			if errLeft == nil && errRight == nil {
				return cleanKey[:lastUnderscore]
			}
		}
	}

	return cleanKey
}
