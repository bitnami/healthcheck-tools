package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	// VERSION will be overwritten automatically by the build system
	VERSION = "devel"
	// ShowHidden include hidden files or directories in the search
	ShowHidden bool
	// Verbose print everything
	Verbose bool
	// Exclude some files and/or directories (RegExp)
	Exclude string
	// DefaultFilePerm are the default default file permissions
	DefaultFilePerm string
	// DefaultDirPerm are the default directory permissions
	DefaultDirPerm string
	// DefaultOwner is the default owner
	DefaultOwner string
	// DefaultGroup is the default group
	DefaultGroup string
	// Directory to check
	Directory string
)

func main() {
	var (
		getVersion bool
	)
	flag.StringVar(&Directory, "dir", "/opt/bitnami", "Directory to check")
	flag.StringVar(&DefaultFilePerm, "file_default", "rw-rw-r--", "Default file permissions")
	flag.StringVar(&DefaultDirPerm, "dir_default", "rwxrwxr-x", "Default directory permissions")
	flag.StringVar(&DefaultOwner, "owner", "bitnami", "Default owner")
	flag.StringVar(&DefaultGroup, "group", "daemon", "Default group")
	flag.StringVar(&Exclude, "exclude", "(?!.*)", "Files and/or directories to be excluded")
	flag.BoolVar(&getVersion, "version", false, "Show current version")
	flag.BoolVar(&ShowHidden, "show_hidden", false, "Show hidden files and directories")
	flag.BoolVar(&Verbose, "verbose", false, "Print every file and directory")
	flag.Parse()

	ManageInputs()

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
	- Exclude: %s
	- Show hidden: %t
	- Verbose: %t
==================================================
`, Directory, DefaultFilePerm, DefaultDirPerm, DefaultOwner, DefaultGroup, Exclude, ShowHidden, Verbose)

	fmt.Printf("\x1b[34;1m\n-- Checking permissions --\n\x1b[0m")
	FindRecursive(Directory, 0)
}
