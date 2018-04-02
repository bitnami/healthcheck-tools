package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	search.baseDirectory = strings.TrimSuffix(filepath.Clean(search.baseDirectory), "/")
	excludeStr = strings.TrimSuffix(filepath.Clean(excludeStr), "/")

	search.exclude = regexp.MustCompile(excludeStr)

	// Checks if the default permissions introduced by the user are in the Linux format
	unixPermRegexp := regexp.MustCompile("^[-rwx]{9}$")
	checks := []struct {
		flag  string
		value string
	}{
		{"file_default", defaultPerm.file},
		{"dir_default", defaultPerm.dir},
	}
	for _, c := range checks {
		if !unixPermRegexp.MatchString(c.value) {
			log.Fatalf(`%s should be in the Unix format (i.e. "rw-rw-r--")`, c.flag)
		}
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
