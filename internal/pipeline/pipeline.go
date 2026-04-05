package pipeline

import (
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
