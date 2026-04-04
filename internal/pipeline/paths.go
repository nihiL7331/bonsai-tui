package pipeline

import (
	"bonsai-tui/internal/config"
	"bonsai-tui/internal/engine"
	"path/filepath"
)

func getEngineDir(cfg config.Config) string {
	return filepath.Join(cfg.ProjectDir, engine.MainDir)
}

func getCacheDir(cfg config.Config) string {
	return filepath.Join(getEngineDir(cfg), engine.CacheDir)
}

func getShaderCacheDir(cfg config.Config) string {
	return filepath.Join(getCacheDir(cfg), "shaders")
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
