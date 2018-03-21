package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
)

// FindRecursive iterate in a directory (path) showing permissions in a recursive way
/* path  - where to look for files and directories
 * level - used to print the output in a hierarchical way
 */
func FindRecursive(path string, level int) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	// In verbose mode we are printing the output in a hierarchical way
	if Verbose {
		level++
	}

	for _, f := range files {
		name := f.Name()

		relativePath := strings.TrimPrefix(strings.Join([]string{path, name}, "/"), Directory)
		fullPath := strings.Join([]string{path, name}, "/")

		if m, _ := regexp.MatchString(Exclude, relativePath); m {
			if Verbose {
				fmt.Printf("\x1b[33m%sExcluding %s\n\x1b[0m", strings.Repeat(" ", level), name)
			}
		} else {
			if ShowHidden || !strings.HasPrefix(name, ".") {
				mode := f.Mode()
				currentOwner, _ := user.LookupId(fmt.Sprint(f.Sys().(*syscall.Stat_t).Uid))
				currentGroup, _ := user.LookupGroupId(fmt.Sprint(f.Sys().(*syscall.Stat_t).Gid))
				currentPerm := mode.Perm().String()
				hasCorrectOwner := currentOwner.Username == DefaultOwner
				hasCorrectGroup := currentGroup.Name == DefaultGroup

				if mode.IsRegular() {
					hasCorrectPermissions := checkPermissions(currentPerm, DefaultFilePerm)
					printOutput(level, "f", fullPath, currentPerm, DefaultFilePerm, currentOwner.Username, currentGroup.Name, hasCorrectPermissions, hasCorrectOwner, hasCorrectGroup)
				} else if mode.IsDir() {
					hasCorrectPermissions := checkPermissions(currentPerm, DefaultDirPerm)
					printOutput(level, "d", fullPath, currentPerm, DefaultDirPerm, currentOwner.Username, currentGroup.Name, hasCorrectPermissions, hasCorrectOwner, hasCorrectGroup)
					FindRecursive(fullPath, level)
				}
			}
		}
	}
}

// checkPermissions return true if the permissions are correct (false in other case)
/* currentPermissions [string] - permissions that the file or directory has (i.e. -rwxrwxr-x)
 * defaultPermissions [string] - permissions that should have (i.e. rwxrwxr-x)
 */
func checkPermissions(currentPermissions, defaultPermissions string) bool {
	return strings.Contains(currentPermissions, defaultPermissions)
}

// printOutput Print the data in different formats according to the situation
/* level [string] 							- used to print the output in a hierarchical way
 * kind [string] 								- "f" if it is file or "d" in case of a directory
 * name [string] 				  			- receive the full path of the file/directory (if verbose only use the name)
 * currentPermissions [string] 	- file/directory current permissions
 * defaultPermissions [string] 	- file/directory default permissions
 * currentOwner [string] 				- file/directory current owner
 * currentGroup [string] 				- file/directory current group
 * hasCorrectPermissions [bool] - true if the permissions are correct
 * hasCorrectOwner [bool] 			- true if the owner is correct
 * hasCorrectGroup [bool] 			- true if the group is correct
 */
func printOutput(level int, kind, name, currentPermissions, defaultPermissions, currentOwner, currentGroup string, hasCorrectPermissions, hasCorrectOwner, hasCorrectGroup bool) {
	if Verbose {
		name = filepath.Base(name)
	}

	if hasCorrectPermissions && hasCorrectOwner && hasCorrectGroup { // Everything correct
		if Verbose {
			fmt.Printf("%s(%s) %s %s %s %s\n", strings.Repeat(" ", level), kind, name, currentPermissions, currentOwner, currentGroup)
		}
	} else if hasCorrectPermissions && (!hasCorrectOwner || !hasCorrectGroup) { // Permissions correct, fails owner or group
		fmt.Printf("\x1b[31;1m%s(%s) %s %s %s %s (expected %s %s)\n\x1b[0m", strings.Repeat(" ", level), kind, name, currentPermissions, currentOwner, currentGroup, DefaultOwner, DefaultGroup)
	} else if !hasCorrectPermissions && (hasCorrectOwner && hasCorrectGroup) { // Permissions wrong, owner and group correct
		fmt.Printf("\x1b[31;1m%s(%s) %s %s (expected %s) %s %s\n\x1b[0m", strings.Repeat(" ", level), kind, name, currentPermissions, defaultPermissions, currentOwner, currentGroup)
	} else if !hasCorrectPermissions && (!hasCorrectOwner || !hasCorrectGroup) { // Nothing correct
		fmt.Printf("\x1b[31;1m%s(%s) %s %s (expected %s) %s %s (expected %s %s)\n\x1b[0m", strings.Repeat(" ", level), kind, name, currentPermissions, defaultPermissions, currentOwner, currentGroup, DefaultOwner, DefaultGroup)
	}
}

// ManageInputs Validate and/or format the parameters entered by the user
func ManageInputs() {
	// Unify the format eliminating the last / if exists
	Directory = strings.TrimSuffix(Directory, "/")
	Exclude = strings.TrimSuffix(Exclude, "/")

	// Check if the default permissions introduced by the user are in the Linux format
	if m, _ := regexp.MatchString("^[-rwx]{9,10}$", DefaultFilePerm); !m {
		log.Fatalf("file_default should be in the Linux format (i.e. \"rw-rw-r--\")\n")
		os.Exit(2)
	} else if m, _ := regexp.MatchString("^[-rwx]{9,10}$", DefaultDirPerm); !m {
		log.Fatalf("dir_default should be in the Linux format (i.e. \"rw-rw-r--\")\n")
		os.Exit(2)
	}
}
