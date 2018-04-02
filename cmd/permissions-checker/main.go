// Package permissionschecks performs a set of health checks in your
// server or local machine to detect any possible issues with the
// permissions configuration comparing the current permissions,
// owner and groups with those that should have.
package permissionschecks

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

// Data regarding default permissions
type defaultPermissions struct {
	file  string // default permissions for files
	dir   string // default permissions for directories
	owner string // default owner
	group string // default group
}

// Data regarding search options
type searchSettings struct {
	hidden        bool           // includes hidden files and directories in the search
	exclude       *regexp.Regexp // files and/or directories to be excluded
	baseDirectory string         // base directory
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

	// Checks if the default permissions introduced by the user are in the Unix format
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

	fmt.Printf(Colorize("blue", "\n-- Checking permissions --\n"))
	FindRecursive(search.baseDirectory, defaultPerm, search, verbose, 0)
}
