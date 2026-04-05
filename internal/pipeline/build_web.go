package pipeline

import (
	"bonsai-tui/internal/config"
	"bonsai-tui/internal/engine"
	"path/filepath"
)

func GetEmscriptenFlags(cfg config.Config) []string {
	flags := []string{
		"-sWASM_BIGINT",
		"-sWARN_ON_UNDEFINED_SYMBOLS=0",
		"-sALLOW_MEMORY_GROWTH",
		"-sINITIAL_MEMORY=67108864",
		"-sMAX_WEBGL_VERSION=2",
		"-sASSERTIONS",
	}

	flags = append(flags,
		"--shell-file", filepath.Join(getEngineDir(cfg), "core/platform/web/index.html"),
		"--preload-file", filepath.Join(cfg.GetAssetsDir(), engine.AudioDir),
		"--preload-file", filepath.Join(cfg.GetAssetsDir(), engine.FontsDir),
		"--preload-file", filepath.Join(getEngineDir(cfg), engine.DefaultFontName),
		"--preload-file", filepath.Join(getCacheDir(cfg), engine.SpriteCacheName),
		"--preload-file", filepath.Join(getCacheDir(cfg), engine.FontCacheDir),
		"--preload-file", filepath.Join(getCacheDir(cfg), engine.AtlasDir),
	)

	return flags
}
