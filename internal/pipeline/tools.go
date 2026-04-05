package pipeline

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func getShdcPlatformDetails() (url string, filename string, err error) {
	baseURL := "https://raw.githubusercontent.com/floooh/sokol-tools-bin/master/bin"

	target := runtime.GOOS + "/" + runtime.GOARCH

	var path string
	var exec string
	switch target {
	case "windows/amd64", "windows/arm64":
		path = "/win32/"
		exec = "sokol-shdc.exe"
	case "linux/amd64":
		path = "/linux/"
		exec = "sokol-shdc"
	case "darwin/amd64":
		path = "/osx/"
		exec = "sokol-shdc"
	case "linux/arm64":
		path = "/linux_arm64/"
		exec = "sokol-shdc"
	case "darwin/arm64":
		path = "/osx_arm64/"
		exec = "sokol-shdc"
	default:
		return "", "", fmt.Errorf("Unsupported platform: %s", target)
	}

	return baseURL + path + exec, exec, nil
}

func EnsureShdc(logFn func(string, string)) (string, error) {
	url, filename, err := getShdcPlatformDetails()
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

	if err := os.MkdirAll(installDir, 0755); err != nil {
		return "", fmt.Errorf("Failed to create install directory: %w", err)
	}

	logFn("TOOLS", "Downloading sokol-shdc...")

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Failed to request sokol-shdc: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Failed to download sokol-shdc: %v", err)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", fmt.Errorf("Failed to write binary to disk: %w", err)
	}

	os.Chmod(destPath, 0o755)

	return destPath, nil
}

func getMsdfGenPlatformDetails() (url string, filename string, err error) {
	baseURL := "https://github.com/nihiL7331/bonsai-tools-bin/releases/download/v1.0/"

	var path string
	var exec string
	switch runtime.GOOS {
	case "windows":
		path = "msdf-atlas-gen-windows.exe"
		exec = "msdf-atlas-gen.exe"
	case "linux":
		path = "msdf-atlas-gen-linux"
		exec = "msdf-atlas-gen"
	case "darwin":
		path = "msdf-atlas-gen-macos-arm64"
		exec = "msdf-atlas-gen"
	default:
		return "", "", fmt.Errorf("Unsupported platform: %s", runtime.GOOS)
	}

	return baseURL + path, exec, nil
}

func EnsureMsdfGen(logFn func(string, string)) (string, error) {
	url, filename, err := getMsdfGenPlatformDetails()
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

	if err := os.MkdirAll(installDir, 0755); err != nil {
		return "", fmt.Errorf("Failed to create install directory: %w", err)
	}

	logFn("TOOLS", "Downloading msdf-atlas-gen...")

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Failed to request msdf-atlas-gen: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Failed to download msdf-atlas-gen, status code: %d", resp.StatusCode)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", fmt.Errorf("Failed to write binary to disk: %w", err)
	}

	os.Chmod(destPath, 0o755)

	return destPath, nil
}
