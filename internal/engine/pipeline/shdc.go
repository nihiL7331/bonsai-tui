package pipeline

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

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
