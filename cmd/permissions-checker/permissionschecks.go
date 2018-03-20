package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
	"regexp"
	"strings"
	"syscall"
)

// FindRecursive iterate in a directory (path) showing permissions in a recursive way
/* path  - where to look for files and directories
 * level - used to print the output in a hierarchical way
 */
func FindRecursive(path, level string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	// In verbose mode we are printing the output in a hierarchical way
	if Verbose {
		level = level + "  "
	}

	for _, f := range files {
		name := f.Name()
		relativePath := strings.TrimPrefix(strings.Join([]string{path, name}, "/"), Directory)
		fullPath := strings.Join([]string{path, name}, "/")

		match, _ := regexp.MatchString(Exclude, relativePath)
		if match {
			if Verbose {
				fmt.Printf("\x1b[33m%sExcluding %s\n\x1b[0m", level, name)
			}
		} else {
			mode := f.Mode()
			currentOwner, _ := user.LookupId(fmt.Sprint(f.Sys().(*syscall.Stat_t).Uid))
			currentGroup, _ := user.LookupGroupId(fmt.Sprint(f.Sys().(*syscall.Stat_t).Gid))
			currentPerm := mode.Perm().String()
			hasCorrectOwner := checkOwner(currentOwner.Username, DefaultOwner)
			hasCorrectGroup := checkOwner(currentGroup.Name, DefaultGroup)

			if mode.IsRegular() && (f.Name()[0:1] != "." || ShowHidden) {
				hasCorrectPermissions := checkPermissions(currentPerm, DefaultFilePerm)
				printOutput(level, "f", name, fullPath, currentPerm, DefaultFilePerm, currentOwner.Username, currentGroup.Name, hasCorrectPermissions, hasCorrectOwner, hasCorrectGroup)
			} else if mode.IsDir() && (f.Name()[0:1] != "." || ShowHidden) {
				hasCorrectPermissions := checkPermissions(currentPerm, DefaultDirPerm)
				printOutput(level, "d", name, fullPath, currentPerm, DefaultDirPerm, currentOwner.Username, currentGroup.Name, hasCorrectPermissions, hasCorrectOwner, hasCorrectGroup)
				FindRecursive(fullPath, level)
			}
		}
	}
}

// checkPermissions return true if the permissions are correct (false in other case)
/* currentPermissions [string] - permissions that the file or directory has
 * defaultPermissions [string] - permissions that should have
 */
func checkPermissions(currentPermissions, defaultPermissions string) bool {
	return strings.Contains(currentPermissions, defaultPermissions)
}

// checkOwner return true if the owner is correct (false in other case)
/* currentOwner [string] - owner of the file or directory
 * defaultOwner [string] - owner that should have the file or directory
 */
func checkOwner(currentOwner, defaultOwner string) bool {
	if currentOwner == defaultOwner {
		return true
	}
	return false
}

// checkGroup return true if the group is correct (false in other case)
/* currentGroup [string] - group of the file or directory
 * defaultGroup [string] - group that should have the file or directory
 */
func checkGroup(currentGroup, defaultGroup string) bool {
	if currentGroup == defaultGroup {
		return true
	}
	return false
}

// printOutput Print the data in different formats according to the situation
/* level [string] 							- used to print the output in a hierarchical way
 * kind [string] 								- "f" if it is file or "d" in case of a directory
 * name [string] 								- file/directory name
 * fullPath [string] 						- full path of the file/directory
 * currentPermissions [string] 	- file/directory current permissions
 * defaultPermissions [string] 	- file/directory default permissions
 * currentOwner [string] 				- file/directory current owner
 * currentGroup [string] 				- file/directory current group
 * hasCorrectPermissions [bool] - true if the permissions are correct
 * hasCorrectOwner [bool] 			- true if the owner is correct
 * hasCorrectGroup [bool] 			- true if the group is correct
 */
func printOutput(level, kind, name, fullPath, currentPermissions, defaultPermissions, currentOwner, currentGroup string, hasCorrectPermissions, hasCorrectOwner, hasCorrectGroup bool) {
	if !Verbose {
		name = fullPath
	}

	if hasCorrectPermissions && hasCorrectOwner && hasCorrectGroup { // Everything correct
		if Verbose {
			fmt.Printf("%s(%s) %s %s %s %s\n", level, kind, name, currentPermissions, currentOwner, currentGroup)
		}
	} else if hasCorrectPermissions && (!hasCorrectOwner || !hasCorrectGroup) { // Permissions correct, fails owner or group
		fmt.Printf("\x1b[31;1m%s(%s) %s %s %s %s (expected %s %s)\n\x1b[0m", level, kind, name, currentPermissions, currentOwner, currentGroup, DefaultOwner, DefaultGroup)
	} else if !hasCorrectPermissions && (hasCorrectOwner && hasCorrectGroup) { // Permissions wrong, owner and group correct
		fmt.Printf("\x1b[31;1m%s(%s) %s %s (expected %s) %s %s\n\x1b[0m", level, kind, name, currentPermissions, defaultPermissions, currentOwner, currentGroup)
	} else if !hasCorrectPermissions && (!hasCorrectOwner || !hasCorrectGroup) { // Nothing correct
		fmt.Printf("\x1b[31;1m%s(%s) %s %s (expected %s) %s %s (expected %s %s)\n\x1b[0m", level, kind, name, currentPermissions, defaultPermissions, currentOwner, currentGroup, DefaultOwner, DefaultGroup)
	}
}
