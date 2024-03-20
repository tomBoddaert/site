package main

import (
	"os/exec"
	"strings"
)

func transpileTS(tsArgs []string) {
	args := make([]string, 0, len(tsArgs)+1)
	args = append(args, "tsc")
	args = append(args, tsArgs...)

	cmd := exec.Command("npx", args...)
	out := new(strings.Builder)
	cmd.Stdout = out

	err := cmd.Run()
	output := out.String()

	if err != nil {
		if strings.Contains(output, "error TS18003") {
			logger.Info("TSC found no files to transpile")
			return
		}

		if strings.Contains(output, "This is not the tsc command you are looking for") {
			logger.Error("TypeScript is not installed! Install it with 'npm i typescript'")
			return
		}

		logger.Errorf("TS files not transpiled!")
	}

	if len(output) != 0 || err != nil {
		logger.Infof("TSC output (%v): %v", strings.Join(cmd.Args, " "), output)
	}
}
