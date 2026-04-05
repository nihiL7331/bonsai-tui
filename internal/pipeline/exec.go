package pipeline

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func RunStreamed(cmdName string, args []string, prefix string, logFn func(string, string)) error {
	cmd := exec.Command(cmdName, args...)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Failed to start %s: %w", cmdName, err)
	}

	go streamLog(stdout, prefix, logFn)
	go streamLog(stderr, prefix+"_ERR", logFn)

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%s script failed: %w", cmdName, err)
	}

	return nil
}

func streamLog(pipe io.ReadCloser, prefix string, logFn func(string, string)) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		logFn(prefix, scanner.Text())
	}
}

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
