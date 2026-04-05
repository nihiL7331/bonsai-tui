package packer

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
)

func RenderAtlas(width int, height int, sprites []PlacedSprite) ([]byte, error) {
	atlasRect := image.Rect(0, 0, width, height)
	atlas := image.NewRGBA(atlasRect)

	for _, sprite := range sprites {
		file, err := os.Open(sprite.Path)
		if err != nil {
			return nil, fmt.Errorf("Failed to open %s: %w", sprite.Path, err)
		}

		srcImg, err := png.Decode(file)
		file.Close()
		if err != nil {
			return nil, fmt.Errorf("Failed to decode %s: %w", sprite.Path, err)
		}

		destRect := image.Rect(
			sprite.Rect.X,
			sprite.Rect.Y,
			sprite.Rect.X+sprite.Rect.W,
			sprite.Rect.Y+sprite.Rect.H,
		)

		draw.Draw(atlas, destRect, srcImg, image.Point{0, 0}, draw.Src)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, atlas); err != nil {
		return nil, fmt.Errorf("Failed to encode atlas: %w", err)
	}

	return buf.Bytes(), nil
}
