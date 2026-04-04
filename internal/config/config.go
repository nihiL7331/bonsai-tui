package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ProjectDir      string `toml:"-"`                 // grabbed at runtime
	AssetsDir       string `toml:"assets_dir"`        // default: assets
	SourceDir       string `toml:"source_dir"`        // default: source
	GameDir         string `toml:"game_dir"`          // default: game
	BuildDir        string `toml:"build_dir"`         // default: build
	BuildWebDir     string `toml:"build_web_dir"`     // default: web
	BuildDesktopDir string `toml:"build_desktop_dir"` // default: desktop
	ShadersDir      string `toml:"shaders_dir"`       // default: shaders
}

func Load() (Config, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		ProjectDir:      pwd,
		AssetsDir:       "assets",
		SourceDir:       "source",
		GameDir:         "game",
		BuildDir:        "build",
		BuildWebDir:     "web",
		BuildDesktopDir: "desktop",
		ShadersDir:      "shaders",
	}

	tomlPath := filepath.Join(pwd, "bonsai.toml")
	_, err = toml.DecodeFile(tomlPath, &cfg)
	if err != nil && !os.IsNotExist(err) {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) GetAssetsDir() string {
	return filepath.Join(c.ProjectDir, c.AssetsDir)
}

func (c Config) GetSourceDir() string {
	return filepath.Join(c.ProjectDir, c.SourceDir)
}

func (c Config) GetGameDir() string {
	return filepath.Join(c.GetSourceDir(), c.GameDir)
}

func (c Config) GetShadersDir() string {
	return filepath.Join(c.GetGameDir(), c.ShadersDir)
}

func (c Config) GetBuildWebDir() string {
	return filepath.Join(c.ProjectDir, c.BuildDir, c.BuildWebDir)
}

func (c Config) GetBuildDektopDir() string {
	return filepath.Join(c.ProjectDir, c.BuildDir, c.BuildDesktopDir)
}
