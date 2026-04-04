package pipeline

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
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
