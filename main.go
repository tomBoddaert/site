package main

import (
	"fmt"
	"os"
	"strings"
)

const HELP = `Usage: %s [options]

  Options:
    help        Display this help text

    build       Build the site
      inclts    Include TypeScript (.ts) files in the output
      fmthtml   Format HTML (.html) files from rawPages

    serve       Serve the site

    debug       Enable debug mode (panic)

    author      Information about the author
`

const AUTHOR = `This program was created by:

  Tom Boddaert
    https://tomBoddaert.github.io/

`

var debugMode = false

func main() {
	// Get the string used to run this program
	cmd := os.Args[0]
	// Get the arguments
	args := os.Args[1:]

	// If no arguments were provided, print the
	//  help text and exit
	if len(args) == 0 {
		fmt.Printf(HELP, cmd)
		fmt.Println("\nNo options provided!")
		os.Exit(1)
	}

	// Set default options
	doHelp := false
	doBuild := false
	doServe := false
	doAuthor := false

	// Set options from args
	for _, arg := range args {
		switch strings.ToLower(arg) {
		case "help":
			doHelp = true

		case "build":
			doBuild = true

		case "inclts":
			includeTSFiles = true

		case "fmthtml":
			formatHTMLFiles = true

		case "serve":
			doServe = true

		case "debug":
			debugMode = true

		case "author":
			doAuthor = true

		default:
			fmt.Printf("Unknown option: '%s'!\nUse '%s help'\n", arg, cmd)
			os.Exit(1)
		}
	}

	if doAuthor {
		fmt.Print(AUTHOR)
	}

	if doHelp {
		fmt.Printf(HELP, cmd)
		if doBuild || doServe {
			fmt.Println("\nOther options used with 'help', ignoring and exiting!")
			os.Exit(0)
		}
	}

	if doBuild {
		build()
		fmt.Println("Site built successfully")
	}

	// If 'inclts' was provided without 'build', print a warning
	if includeTSFiles && !doBuild {
		fmt.Println("\n'inclts' option used without 'build', ignoring")
	}

	// If 'fmthtml' was provided without 'build', print a warning
	if formatHTMLFiles && !doBuild {
		fmt.Println("\n'fmthtml' option used without 'build', ignoring")
	}

	if doServe {
		serve()
	}
}

// Error handling
func check(err error) {
	if err != nil {
		// In debug mode, panic, otherwise print some generic
		//  error text and exit with an error code
		if debugMode {
			panic(err)
		} else {
			fmt.Println("There was an error!")
			fmt.Println("Please use debug mode or contact me (use the 'author' subcommand)")
			os.Exit(1)
		}
	}
}
