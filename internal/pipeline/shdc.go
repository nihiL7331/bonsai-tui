package pipeline

import (
	"bonsai-tui/internal/config"
	"bonsai-tui/internal/engine"
	"fmt"
	"io"
	"net/http"
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

func getPlatformDetails() (url string, filename string, err error) {
	baseURL := "https://raw.githubusercontent.com/floooh/sokol-tools-bin/master/bin"

	target := runtime.GOOS + "/" + runtime.GOARCH

	var path string
	var exec string
	switch target {
	case "windows/amd64", "windows/arm64":
		path = "/win32/"
		exec = "sokol-shdc.exe"
		break
	case "linux/amd64":
		path = "/linux/"
		exec = "sokol-shdc"
		break
	case "darwin/amd64":
		path = "/osx/"
		exec = "sokol-shdc"
		break
	case "linux/arm64":
		path = "/linux_arm64/"
		exec = "sokol-shdc"
		break
	case "darwin/arm64":
		path = "/osx_arm64/"
		exec = "sokol-shdc"
		break
	default:
		return "", "", fmt.Errorf("Unsupported platform: %s", target)
	}

	return baseURL + path + exec, exec, nil
}

func EnsureShdc() (string, error) {
	url, filename, err := getPlatformDetails()
	if err != nil {
		return "", err
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("Cache directory not found: %w", err)
	}

	installDir := filepath.Join(cacheDir, "bonsai", "bin")
	destPath := filepath.Join(installDir, filename)

	if _, err := os.Stat(destPath); err == nil {
		return destPath, nil
	}

	os.MkdirAll(installDir, 0755)

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return "", fmt.Errorf("Failed to download sokol-shdc: %v", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer out.Close()
	io.Copy(out, resp.Body)

	os.Chmod(destPath, 0o755)

	return destPath, nil
}

func CompileShaders(cfg config.Config, logFn func(string, string)) error {
	shdcPath, err := EnsureShdc()
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
