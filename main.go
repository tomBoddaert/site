package main

import (
	"os"
	"strings"

	"github.com/charmbracelet/log"
)

const HELP = `Usage: %s [options]

  Options:
    help        Display this help text
      author    Include information about the author

    build       Build the site

    serve       Serve the site

    debug       Enable debug mode

  Environment Variables:
    SITE_ADDRESS    The address and port for serve to listen on
                    (address:port)`

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
		logger.Fatalf(err.Error())
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
		Serve(config)
	}
}

func build(config Config) {
	if len(config.PrebuildCmds) != 0 {
		logger.SetReportTimestamp(true)
		logger.Info("Running prebuild commands...")
		logger.SetReportTimestamp(false)

		RunCmds(config.PrebuildCmds)
	} else {
		logger.Debug("No prebuild commands to run")
	}

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

	CopyRaw(&config)

	logger.SetReportTimestamp(true)
	logger.Info("Building templated pages...")
	logger.SetReportTimestamp(false)

	BuildTemplated(&config)

	logger.SetReportTimestamp(true)
	logger.Info("Replacing old build...")
	logger.SetReportTimestamp(false)

	logger.Debugf("Renaming directory (%v -> %v)", tempDir, outDir)

	check(os.RemoveAll(outDir))
	check(os.Rename(tempDir, outDir))

	if len(config.PostbuildCmds) != 0 {
		logger.SetReportTimestamp(true)
		logger.Info("Running postbuild commands...")
		logger.SetReportTimestamp(false)

		RunCmds(config.PostbuildCmds)
	} else {
		logger.Debug("No postbuild commands to run")
	}

	logger.SetReportTimestamp(true)
	logger.Info("Done")
	logger.SetReportTimestamp(false)
}
