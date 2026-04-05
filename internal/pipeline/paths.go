package pipeline

import (
	"bonsai-tui/internal/config"
	"bonsai-tui/internal/engine"
	"path/filepath"
)

func getOdinAtlasDir(cfg config.Config) string {
	return filepath.Join(getEngineDir(cfg), engine.GeneratedDir, "sprites.bin")
}

func getBinAtlasDir(cfg config.Config) string {
	return filepath.Join(getAtlasCacheDir(cfg), "metadata.bin")
}

func getEngineDir(cfg config.Config) string {
	return filepath.Join(cfg.ProjectDir, engine.MainDir)
}

func getImagesDir(cfg config.Config) string {
	return filepath.Join(cfg.GetAssetsDir(), engine.ImagesDir)
}

func getCacheDir(cfg config.Config) string {
	return filepath.Join(getEngineDir(cfg), engine.CacheDir)
}

func getAtlasCacheDir(cfg config.Config) string {
	return filepath.Join(getCacheDir(cfg), engine.AtlasDir)
}

func getShaderCacheDir(cfg config.Config) string {
	return filepath.Join(getCacheDir(cfg), engine.ShaderDir)
}

func getFontCacheDir(cfg config.Config) string {
	return filepath.Join(getCacheDir(cfg), engine.FontsDir)
}

func getShaderIncludeDir(cfg config.Config) string {
	return filepath.Join(getEngineDir(cfg), engine.ShaderIncludeDir)
}

func getShaderSrcDir(cfg config.Config) string {
	return filepath.Join(cfg.ProjectDir, engine.MainDir, engine.ShaderSrcName)
}

func getShaderOutDir(cfg config.Config) string {
	return filepath.Join(cfg.ProjectDir, engine.MainDir, engine.ShaderOutName)
}
