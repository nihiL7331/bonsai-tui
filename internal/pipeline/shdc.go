package pipeline

import (
	"bonsai-tui/internal/config"
	"bonsai-tui/internal/engine"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func getShaderFormat(cfg config.Config) string {
	switch cfg.TargetOS {
	case "windows":
		return "glsl300es:hlsl4:glsl430"
	case "darwin":
		return "metal_macos:glsl300es:hlsl4:glsl430"
	case "web":
		return "glsl300es"
	default:
		return "glsl300"
	}
}

func CompileShaders(cfg config.Config, logFn func(string, string)) error {
	shdcPath, err := EnsureShdc(logFn)
	if err != nil {
		return fmt.Errorf("SHDC error: %w", err)
	}

	cacheDir := getShaderCacheDir(cfg)
	_ = os.RemoveAll(cacheDir)
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return fmt.Errorf("Failed to create shader cache: %w", err)
	}

	includeDir := getShaderIncludeDir(cfg)
	includes := map[string]string{
		engine.ShaderVsCoreName: "shader_vs_core.glsl",
		engine.ShaderFsCoreName: "shader_fs_core.glsl",
		engine.ShaderUtilsName:  "shader_utils.glsl",
		engine.ShaderHeaderName: "shader_header.glsl",
	}

	for srcName, dstName := range includes {
		src := filepath.Join(includeDir, srcName)
		dst := filepath.Join(cacheDir, dstName)
		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("Failed to cache include %s: %w", srcName, err)
		}
	}

	shaderFormat := "metal_macos:glsl300es:hlsl4:glsl430"
	if runtime.GOOS == "windows" {
		shaderFormat = "glsl300es:hlsl4:glsl430"
	}

	compileCached := func(srcPath, outPath, prefix string) error {
		filename := filepath.Base(srcPath)
		cachedPath := filepath.Join(cacheDir, filename)

		if err := copyFile(srcPath, cachedPath); err != nil {
			return err
		}

		logFn(prefix, fmt.Sprintf("Compiling: %s", filepath.Base(srcPath)))

		args := []string{
			"-i", cachedPath,
			"-o", outPath,
			"-l", shaderFormat,
			"-f", "sokol_odin",
		}
		return RunStreamed(shdcPath, args, prefix, logFn)
	}

	coreSrc := getShaderSrcDir(cfg)
	coreOut := getShaderOutDir(cfg)

	if !isUpToDate(coreSrc, coreOut) {
		if err := compileCached(coreSrc, coreOut, "CORE SHDC"); err != nil {
			return err
		}
	}

	gameShadersDir := cfg.GetShadersDir()

	err = filepath.WalkDir(gameShadersDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".glsl" && ext != ".vert" && ext != ".frag" {
			return nil
		}

		outPath := strings.TrimSuffix(path, ext) + ".odin"

		if isUpToDate(path, outPath) {
			return nil
		}

		return compileCached(path, outPath, "GAME SHDC")
	})

	return err
}
