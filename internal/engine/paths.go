package engine

import (
	"bonsai-tui/internal/config"
	"path/filepath"
	"runtime"
)

const (
	MainDir           = "bonsai"
	AtlasDir          = "atlas"
	AudioDir          = "audio"
	FontsDir          = "fonts"
	CacheDir          = ".cache"
	ShadersIncludeDir = "shaders/include"
	ShaderVsCoreName  = "shader_vs_core/shader_vs_core.glsl"
	ShaderFsCoreName  = "shader_fs_core/shader_fs_core.glsl"
	ShaderHeaderName  = "shader_header/shader_header.glsl"
	ShaderUtilsName   = "shader_utils/shader_utils.glsl"
	ShadersSrcName    = "shaders/shader.glsl"
	ShadersOutName    = "shaders/shader.odin"
	DefaultFontName   = "core/render/PixelCode_9.ttf"
	SpriteCacheName   = "sprites/sprites.bin"
	FontCacheDir      = "fonts"
)

func GetDesktopBinName() string {
	if runtime.GOOS == "windows" {
		return "game_desktop.exe"
	}
	return "game_desktop.bin"
}

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
		"--shell-file", filepath.Join(cfg.ProjectDir, MainDir, "core/platform/web/index.html"),
		"--preload-file", filepath.Join(cfg.ProjectDir, MainDir, CacheDir, AtlasDir),
		"--preload-file", filepath.Join(cfg.GetAssetsDir(), AudioDir),
		"--preload-file", filepath.Join(cfg.GetAssetsDir(), FontsDir),
		"--preload-file", filepath.Join(cfg.ProjectDir, MainDir, DefaultFontName),
		"--preload-file", filepath.Join(cfg.ProjectDir, MainDir, CacheDir, SpriteCacheName),
		"--preload-file", filepath.Join(cfg.ProjectDir, MainDir, CacheDir, FontCacheDir),
	)

	return flags
}
