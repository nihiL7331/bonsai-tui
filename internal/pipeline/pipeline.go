package pipeline

import (
	"bonsai-tui/internal/config"
	"fmt"
	"os/exec"
)

func checkDependencies(logFn func(string, string)) (string, string, error) {
	requiredTools := map[string]string{
		"odin": "Odin compiler not found in PATH. Please install it from https://odin-lang.org/docs/install",
		"emcc": "Could not find Emscripten SDK.\n Please install it (https://emscripten.org/docs/getting_started/downloads.html)\n and set the 'EMSDK' environment variable to its installation folder",
	}

	for bin, errStr := range requiredTools {
		if _, err := exec.LookPath(bin); err != nil {
			return "", "", fmt.Errorf(errStr)
		}
	}

	shdcPath, err := EnsureShdc(logFn)
	if err != nil {
		return "", "", fmt.Errorf("Failed to set up sokol-shdc: %w", err)
	}

	msdfPath, err := EnsureMsdfGen(logFn)
	if err != nil {
		return "", "", fmt.Errorf("Failed to set up msdf-atlas-gen: %w", err)
	}

	return shdcPath, msdfPath, nil
}

func prepareResources(cfg config.Config, logFn func(string, string)) error {
	logFn("PIPELINE", "Running pre-build assembly tasks...")

	shdcPath, msdfPath, err := checkDependencies(logFn)
	if err != nil {
		return fmt.Errorf("Dependency check failed: %w", err)
	}

	if err := RunPreBuildHooks(cfg.ScriptsDir, logFn); err != nil {
		return fmt.Errorf("Build-time scripts failed: %w", err)
	}

	if err := UpdateManifest(cfg.ProjectDir, logFn); err != nil {
		return fmt.Errorf("Failed to update manifest: %w", err)
	}

	if _, err := PackAtlas(cfg, logFn); err != nil {
		return fmt.Errorf("Failed to pack atlas: %w", err)
	}

	if err := GenerateAssets(cfg, msdfPath, logFn); err != nil {
		return fmt.Errorf("Failed to generate assets: %w", err)
	}

	if err := CompileShaders(cfg, shdcPath, logFn); err != nil {
		return fmt.Errorf("Failed to compile shaders: %w", err)
	}

	logFn("PIPELINE", "All resources prepared successfully.")
	return nil
}
