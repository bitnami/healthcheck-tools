package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	// VERSION will be overwritten automatically by the build system
	VERSION    = "devel"
	excludeStr string
)

/* Data regarding default permissions
 * 	defaultPerm.file [string]					- default permissions for files
 * 	defaultPerm.dir [string] 					- default permissions for directories
 * 	defaultPerm.owner [string] 				- default owner
 * 	defaultPerm.group [string] 				- default group
 */
type defaultPermissions struct {
	file  string
	dir   string
	owner string
	group string
}

/* Data regarding search options
 * 	search.hidden [bool] 							- includes hidden files and directories in the search
 * 	search.exclude [*regexp.Regexp] 	- files and/or directories to be excluded
 * 	search.baseDirectory [string]			- base directory
 */
type searchSettings struct {
	hidden        bool
	exclude       *regexp.Regexp
	baseDirectory string
}

func main() {
	var (
		getVersion  bool
		defaultPerm defaultPermissions
		search      searchSettings
		verbose     bool
	)

	flag.StringVar(&defaultPerm.file, "file_default", "rw-rw-r--", "Default file permissions")
	flag.StringVar(&defaultPerm.dir, "dir_default", "rwxrwxr-x", "Default directory permissions")
	flag.StringVar(&defaultPerm.owner, "owner", "bitnami", "Default owner")
	flag.StringVar(&defaultPerm.group, "group", "daemon", "Default group")
	flag.StringVar(&search.baseDirectory, "dir", "/opt/bitnami", "Directory to check")
	flag.StringVar(&excludeStr, "exclude", "(.?!.*)", "Files and/or directories to be excluded")
	flag.BoolVar(&search.hidden, "hidden", false, "Includes hidden files and directories")
	flag.BoolVar(&verbose, "verbose", false, "Print every file and directory")
	flag.BoolVar(&getVersion, "version", false, "Show current version")
	flag.Parse()

	// Unifies the format eliminating the last "/" if exists
	search.baseDirectory = strings.TrimSuffix(search.baseDirectory, "/")
	excludeStr = strings.TrimSuffix(excludeStr, "/")

	search.exclude = regexp.MustCompile(excludeStr)

	// Checks if the default permissions introduced by the user are in the Linux format
	defLinuxFormat := regexp.MustCompile("^[-rwx]{9,10}$")
	if !defLinuxFormat.MatchString(defaultPerm.file) || !defLinuxFormat.MatchString(defaultPerm.dir) {
		log.Fatalf("file_default and dir_default should be in the Linux format (i.e. \"rw-rw-r--\")\n")
		os.Exit(2)
	}

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
	- Include hidden: %t
	- Verbose: %t
==================================================
`, search.baseDirectory, defaultPerm.file, defaultPerm.dir, defaultPerm.owner, defaultPerm.group, search.exclude, search.hidden, verbose)

	fmt.Printf("\x1b[34;1m\n-- Checking permissions --\n\x1b[0m")
	FindRecursive(search.baseDirectory, defaultPerm, search, verbose, 0)
}
