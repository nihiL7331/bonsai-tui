package pipeline

import (
	"fmt"
	"os"
	"path/filepath"
)

func RunUtils(utilsDir string, logFn func(string, string)) error {
	if _, err := os.Stat(utilsDir); os.IsNotExist(err) {
		return nil
	}

	return filepath.WalkDir(utilsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		switch filepath.Ext(path) {
		case ".py":
			return runPy(path, logFn)
		case ".odin":
			logFn("ODIN", fmt.Sprintf("Running script: %s", path))
			return RunStreamed("odin", []string{"run", path, "-file"}, "ODIN", logFn)
		}

		return nil
	})
}

func runPy(path string, logFn func(string, string)) error {
	logFn("PYTHON", fmt.Sprintf("Running script: %s", path))

	err := RunStreamed("python3", []string{path}, "PYTHON", logFn)
	if err != nil { // fallback to 'python' command
		return RunStreamed("python", []string{path}, "PYTHON", logFn)
	}
	return nil
}
