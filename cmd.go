package main

import (
	"os"
	"os/exec"
	"unicode/utf8"
)

const maxOutputRuneLength uint = 80
const maxOutputByteLength uint = maxOutputRuneLength * utf8.UTFMax

func RunCmds(cmds [][]string) {
	for _, cmd := range cmds {
		runCmd(cmd)
	}
}

func runCmd(cmd []string) {
	logger.Debugf("Running command (%v)", cmd)

	if len(cmd) == 0 {
		logger.Warn("Empty command in config, ignoring")
		return
	}

	prog := cmd[0]
	args := cmd[1:]
	// Limit the outputs in length to avoid allocating large amounts of text
	// that will later get truncated
	stdout := new(MaxStringBuilder)
	stderr := new(MaxStringBuilder)
	stdout.MaxLen = maxOutputByteLength
	stderr.MaxLen = maxOutputByteLength

	runner := exec.Command(prog, args...)
	runner.Stdout = stdout
	runner.Stderr = stderr

	err := runner.Run()

	logger.SetOutput(os.Stderr)
	if err != nil {
		logger.Errorf("A command failed: %q", cmd)

		if stdout.Len() != 0 {
			stdoutStr := stdout.String()
			logger.Infof("Command output: %q", truncate(stdoutStr, maxOutputRuneLength))
		}

		if stderr.Len() != 0 {
			stderrStr := stderr.String()
			logger.Infof("Command error output: %q", truncate(stderrStr, maxOutputRuneLength))
		} else {
			logger.Info("Command produced no error output.")
		}
		logger.SetOutput(os.Stdout)

	} else {
		if stderr.Len() != 0 {
			stderrStr := stderr.String()
			logger.Debugf("Command error output: %v", truncate(stderrStr, maxOutputRuneLength))
		}
		logger.SetOutput(os.Stdout)

		if stdout.Len() != 0 {
			stdoutStr := stdout.String()
			logger.Debugf("Command output: %v", truncate(stdoutStr, maxOutputRuneLength))
		}
	}
}
