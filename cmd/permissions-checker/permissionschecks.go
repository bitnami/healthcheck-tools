package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
)

type currentPermissions struct {
	permissions           string
	defaultPermissions    string
	owner                 string
	group                 string
	hasCorrectPermissions bool
	hasCorrectOwner       bool
	hasCorrectGroup       bool
}

// FindRecursive iterates in a directory showing permissions in a recursive way
/* path [string] 											- where to look for files and directories
 * defaultPerm [defaultPermissions] 	- data regarding default permissions
 * 	defaultPerm.file [string]					- default permissions for files
 * 	defaultPerm.dir [string] 					- default permissions for directories
 * 	defaultPerm.owner [string] 				- default owner
 * 	defaultPerm.group [string] 				- default group
 * search [searchSettings] 						- data regarding search options
 * 	search.hidden [bool] 							- includes hidden files and directories in the search
 * 	search.exclude [*regexp.Regexp] 	- files and/or directories to be excluded
 * 	search.baseDirectory [string]			- base directory
 * level [int]                        - used to print the output in a hierarchical way
 */
func FindRecursive(path string, defaultPerm defaultPermissions, search searchSettings, verbose bool, level int) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	// In verbose mode, we are printing the output in a hierarchical way
	if verbose {
		level++
	}

	for _, f := range files {
		name := f.Name()

		relativePath := strings.TrimPrefix(strings.Join([]string{path, name}, "/"), search.baseDirectory)
		fullPath := strings.Join([]string{path, name}, "/")

		if search.exclude.MatchString(relativePath) {
			if verbose {
				fmt.Printf(Colorize("yellow", fmt.Sprintf("%sExcluding %s\n", strings.Repeat(" ", level), name)))
			}
		} else {
			if search.hidden || !strings.HasPrefix(name, ".") {
				mode := f.Mode()
				currentOwner, _ := user.LookupId(fmt.Sprint(f.Sys().(*syscall.Stat_t).Uid))
				currentGroup, _ := user.LookupGroupId(fmt.Sprint(f.Sys().(*syscall.Stat_t).Gid))

				currentPerm := currentPermissions{
					permissions:     mode.Perm().String(),
					owner:           currentOwner.Username,
					group:           currentGroup.Name,
					hasCorrectOwner: currentOwner.Username == defaultPerm.owner,
					hasCorrectGroup: currentGroup.Name == defaultPerm.group,
				}

				if mode.IsRegular() {
					currentPerm.defaultPermissions = defaultPerm.file
					currentPerm.hasCorrectPermissions = checkPermissions(currentPerm.permissions, currentPerm.defaultPermissions)
					printOutput(level, "f", fullPath, currentPerm, defaultPerm, verbose)
				} else if mode.IsDir() {
					currentPerm.defaultPermissions = defaultPerm.dir
					currentPerm.hasCorrectPermissions = checkPermissions(currentPerm.permissions, currentPerm.defaultPermissions)
					printOutput(level, "d", fullPath, currentPerm, defaultPerm, verbose)
					FindRecursive(fullPath, defaultPerm, search, verbose, level)
				}
			}
		}
	}
}

// checkPermissions returns true if the permissions are correct (false in another case)
/* currentPermissions [string] - permissions that the file or directory has (i.e. -rwxrwxr-x)
 * defaultPermissions [string] - permissions that should have (i.e. rwxrwxr-x)
 */
func checkPermissions(currentPermissions, defaultPermissions string) bool {
	return strings.Contains(currentPermissions, defaultPermissions)
}

// printOutput Prints the data in different formats according to the situation
/* level [string] 														- used to print the output in a hierarchical way
 * kind [string] 															- "f" if it is file or "d" in case of a directory
 * fullPath [string] 				  								- file/directory full path (if verbose only use the name)
 * currentPerm [currentPermissions]						- data to print relate
 * 	currentPerm.permissions [string] 					- file/directory current permissions
 * 	currentPerm.defaultPermissions [string] 	- file/directory default permissions
 * 	currentPerm.owner [string] 								- file/directory current owner
 * 	currentPerm.group [string] 								- file/directory current group
 * 	currentPerm.hasCorrectPermissions [bool] 	- true if the permissions are correct
 * 	currentPerm.hasCorrectOwner [bool] 				- true if the owner is correct
 * 	currentPerm.hasCorrectGroup [bool] 				- true if the group is correct
 * defaultPerm [defaultPermissions] 					- data regarding default permissions
 * 	defaultPerm.file [string]									- default permissions for files
 * 	defaultPerm.dir [string] 									- default permissions for directories
 * 	defaultPerm.owner [string] 								- default owner
 * 	defaultPerm.group [string] 								- default group
 */
func printOutput(level int, kind, fullPath string, currentPerm currentPermissions, defaultPerm defaultPermissions, verbose bool) {
	if verbose {
		fullPath = filepath.Base(fullPath)
	}
	hierarchy := strings.Repeat(" ", level)

	if currentPerm.hasCorrectPermissions && currentPerm.hasCorrectOwner && currentPerm.hasCorrectGroup { // Everything correct
		if verbose {
			fmt.Printf("%s(%s) %s %s %s %s\n", hierarchy, kind, fullPath, currentPerm.permissions, currentPerm.owner, currentPerm.group)
		}
	} else if currentPerm.hasCorrectPermissions && (!currentPerm.hasCorrectOwner || !currentPerm.hasCorrectGroup) { // Permissions correct, fails owner or group
		fmt.Printf(Colorize("red", fmt.Sprintf("%s(%s) %s %s %s %s (expected %s %s)\n", hierarchy, kind, fullPath, currentPerm.permissions, currentPerm.owner, currentPerm.group, defaultPerm.owner, defaultPerm.group)))
	} else if !currentPerm.hasCorrectPermissions && (currentPerm.hasCorrectOwner && currentPerm.hasCorrectGroup) { // Permissions wrong, owner and group correct
		fmt.Printf(Colorize("red", fmt.Sprintf("%s(%s) %s %s (expected %s) %s %s\n", hierarchy, kind, fullPath, currentPerm.permissions, currentPerm.defaultPermissions, currentPerm.owner, currentPerm.group)))
	} else if !currentPerm.hasCorrectPermissions && (!currentPerm.hasCorrectOwner || !currentPerm.hasCorrectGroup) { // Nothing correct
		fmt.Printf(Colorize("red", fmt.Sprintf("%s(%s) %s %s (expected %s) %s %s (expected %s %s)\n", hierarchy, kind, fullPath, currentPerm.permissions, currentPerm.defaultPermissions, currentPerm.owner, currentPerm.group, defaultPerm.owner, defaultPerm.group)))
	}
}

// Colorize Returns a string using ansi colors
func Colorize(color, s string) string {
	const (
		esc        = "\x1b"
		ansiBlue   = esc + "[34;1m"
		ansiRed    = esc + "[31;1m"
		ansiYellow = esc + "[33m"
		ansiReset  = esc + "[0m"
	)
	result := s

	if strings.Compare(color, "blue") == 0 {
		result = fmt.Sprintf("%s%s%s", ansiBlue, string(s), ansiReset)
	} else if strings.Compare(color, "red") == 0 {
		result = fmt.Sprintf("%s%s%s", ansiRed, string(s), ansiReset)
	} else if strings.Compare(color, "yellow") == 0 {
		result = fmt.Sprintf("%s%s%s", ansiYellow, string(s), ansiReset)
	}
	return result
}
