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
	// Ignore Files and/or directories to be ignored
	Ignore string
)

func main() {
	var (
		directory   string
		fileDefault string
		dirDefault  string
		owner       string
		group       string
		getVersion  bool
	)
	flag.StringVar(&directory, "dir", "/opt/bitnami", "Directory to check")
	flag.StringVar(&fileDefault, "file_default", "rw-rw-r--", "Default file permissions")
	flag.StringVar(&dirDefault, "dir_default", "rwxrwxr-x", "Default directory permissions")
	flag.StringVar(&owner, "owner", "bitnami", "Default owner")
	flag.StringVar(&group, "group", "daemon", "Default group")
	flag.StringVar(&Ignore, "ignore", "", "Files and/or directories to be ignored")
	flag.BoolVar(&getVersion, "version", false, "Show current version")
	flag.BoolVar(&ShowHidden, "show_hidden", false, "Show hidden files and directories")
	flag.BoolVar(&Verbose, "verbose", false, "Print every file and directory")
	flag.Parse()

	directory = strings.TrimSuffix(directory, "/")
	Ignore = strings.TrimSuffix(Ignore, "/")

	if getVersion {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	fmt.Printf(`==================================================
	PERMISSIONS CHECKS
==================================================
Starting checks with these parameters:
	- Directory to check: %s
	- Default file permissions: %s
	- Default dir permissions: %s
	- Owner: %s
	- Group: %s
	- Ignore: %s
	- Show hidden: %t
	- Verbose: %t
==================================================
`, directory, fileDefault, dirDefault, owner, group, Ignore, ShowHidden, Verbose)

	fmt.Printf("\x1b[34;1m\n-- Checking permissions --\n\x1b[0m")
	FindRecursive(directory, "", fileDefault, dirDefault, owner, group)
}
