package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	// VERSION will be overwritten automatically by the build system
	VERSION = "devel"
	// ShowHidden include hidden files or directories in the loops
	ShowHidden bool
	// Verbose print everything
	Verbose bool
)

func main() {
	var (
		directory   string
		application string
		fileDefault string
		dirDefault  string
		getVersion  bool
	)
	flag.StringVar(&directory, "dir", "/opt/bitnami", "Starting directory")
	flag.StringVar(&application, "application", "", "Application")
	flag.StringVar(&fileDefault, "file_default", "rw-rw-r--", "File default permissions")
	flag.StringVar(&dirDefault, "dir_default", "rwxrwxr-x", "Directory default permissions")
	flag.BoolVar(&getVersion, "version", false, "Show current version")
	flag.BoolVar(&ShowHidden, "show_hidden", false, "Show hidden files and directories")
	flag.BoolVar(&Verbose, "verbose", false, "Print every file and directory")
	flag.Parse()

	if getVersion {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	fmt.Printf(`==================================================
	PERMISSIONS CHECKS
==================================================
Starting checks with these parameters:
	- Directory to check: %s
	- Application: %s
	- Default file permissions: %s
	- Default dir permissions: %s
	- Show hidden: %t
	- Verbose: %t
==================================================
`, directory, application, fileDefault, dirDefault, ShowHidden, Verbose)

	fmt.Printf("\x1b[34;1m\n-- Checking %s permissions --\n\x1b[0m", strings.Title(application))
	FindRecursive(directory, "", fileDefault, dirDefault)
}
