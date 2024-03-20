package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
)

const HELP = `Usage: %s [options]

  Options:
    help        Display this help text
      author    Include information about the author

    build       Build the site

    serve       Serve the site

    debug       Enable debug mode`

const AUTHOR = `
This program was created by:

  Tom Boddaert
    https://tomBoddaert.com/`

var logger = log.NewWithOptions(
	os.Stdout,
	log.Options{
		TimeFormat: "15:04:05",
	},
)

func check(err error) {
	logger.Helper()

	if err != nil {
		logger.SetOutput(os.Stderr)
		logger.Fatalf("Error: %v", err)
	}
}

func enableDebug() {
	logger.SetLevel(log.DebugLevel)
	logger.SetReportCaller(true)

	logger.Debug("Debug mode enabled")
}

func main() {
	cmd := os.Args[0]
	args := os.Args[1:]

	doHelp := false
	doAuthor := false
	doBuild := false
	doServe := false
	debugMode := false

	for _, arg := range args {
		switch strings.ToLower(arg) {
		case "help":
			doHelp = true

		case "author":
			doAuthor = true

		case "build":
			doBuild = true

		case "serve":
			doServe = true

		case "debug":
			debugMode = true

		default:
			logger.Fatalf("Unknown option: '%s'! Use '%s help'", arg, cmd)
		}
	}

	if debugMode {
		enableDebug()
	}

	if doHelp {
		logger.Printf(HELP, cmd)

		if doAuthor {
			logger.Print(AUTHOR)
		}

		if doBuild || doServe {
			logger.Warn("Other options used with 'help', ignoring.")
		}

		return
	}

	if doAuthor {
		logger.Warn("'author' used without 'help', ignoring.")
	}

	if !doBuild && !doServe {
		logger.Fatalf("Nothing to do! Use '%s help'!", cmd)
	}

	config := getConfig()

	if doBuild {
		build(config)
	}

	if doServe {
		serve(config)
	}
}

func build(config Config) {
	// Parse the destination file mode
	modeNumber, err := strconv.ParseUint(config.DstMode, 8, 32)
	check(err)
	dstMode := os.FileMode(modeNumber)

	logger.SetReportTimestamp(true)
	logger.Info("Creating temp directory...")
	logger.SetReportTimestamp(false)

	tempDir, err := os.MkdirTemp(".", "site-build-*")
	check(err)
	// Set the output directory to the temp directory
	outDir := config.DstDir
	config.DstDir = tempDir

	logger.Debugf("Created temp directory (%v)", tempDir)

	logger.SetReportTimestamp(true)
	logger.Info("Copying raw pages...")
	logger.SetReportTimestamp(false)

	copyRaw(&config, dstMode)

	logger.SetReportTimestamp(true)
	logger.Info("Building templated pages...")
	logger.SetReportTimestamp(false)

	buildTemplated(&config, dstMode)

	logger.SetReportTimestamp(true)
	logger.Info("Replacing old build...")
	logger.SetReportTimestamp(false)

	logger.Debugf("Renaming directory (%v -> %v)", tempDir, outDir)

	check(os.RemoveAll(outDir))
	check(os.Rename(tempDir, outDir))

	if config.TranspileTS {
		logger.SetReportTimestamp(true)
		logger.Info("Transpiling TS...")
		logger.SetReportTimestamp(false)

		transpileTS(config.TSArgs)
	} else {
		logger.Info("Skipping transpiling TS (set in config)")
	}

	logger.SetReportTimestamp(true)
	logger.Info("Done")
	logger.SetReportTimestamp(false)
}
