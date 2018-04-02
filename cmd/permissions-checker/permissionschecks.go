package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
)

// Data regarding current permissions of the item (file or directory)
type currentPermissions struct {
	permissions           os.FileMode // item permissions
	defaultPermissions    string      // default permissions according to item type
	owner                 string      // item owner
	group                 string      // item group
	hasCorrectPermissions bool        // true if item permissions = default permissions
	hasCorrectOwner       bool        // true if item owner = default owner
	hasCorrectGroup       bool        // true if item group = default group
}

// FindPermissions iterates in a directory showing permissions in a recursive way
func FindPermissions(dir string, defaultPerm defaultPermissions, search searchSettings, verbose bool) {
	level := 0
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// In verbose mode, we are printing the output in a hierarchical way
		if verbose {
			level = strings.Count(path, "/") - strings.Count(dir, "/")
		}
		name := info.Name()
		if search.exclude.MatchString(path) {
			if verbose {
				fmt.Printf(Colorize("yellow", fmt.Sprintf("%sExcluding %s\n", strings.Repeat("  ", level), name)))
			}
		} else {
			if search.hidden || !strings.HasPrefix(name, ".") {
				mode := info.Mode()
				currentOwner, _ := user.LookupId(fmt.Sprint(info.Sys().(*syscall.Stat_t).Uid))
				currentGroup, _ := user.LookupGroupId(fmt.Sprint(info.Sys().(*syscall.Stat_t).Gid))

				currentPerm := currentPermissions{
					permissions:     mode.Perm(),
					owner:           currentOwner.Username,
					group:           currentGroup.Name,
					hasCorrectOwner: currentOwner.Username == defaultPerm.owner,
					hasCorrectGroup: currentGroup.Name == defaultPerm.group,
				}

				if mode.IsRegular() {
					currentPerm.defaultPermissions = defaultPerm.file
					currentPerm.hasCorrectPermissions = checkPermissions(currentPerm.permissions, currentPerm.defaultPermissions)
					printOutput(level, "f", path, currentPerm, defaultPerm, verbose)
				} else if mode.IsDir() {
					currentPerm.defaultPermissions = defaultPerm.dir
					currentPerm.hasCorrectPermissions = checkPermissions(currentPerm.permissions, currentPerm.defaultPermissions)
					printOutput(level, "d", path, currentPerm, defaultPerm, verbose)
				} else {
					if verbose {
						fmt.Printf(Colorize("yellow", fmt.Sprintf("%s(o) %s %s %s %s\n", strings.Repeat("  ", level), name, currentPerm.permissions, currentPerm.owner, currentPerm.group)))
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", dir, err)
	}
}

// checkPermissions returns true if the permissions are correct (false in another case)
func checkPermissions(currentPermissions os.FileMode, defaultPermissions string) bool {
	return strings.Contains(currentPermissions.String(), defaultPermissions)
}

// printOutput prints the data in different formats according to the situation
func printOutput(level int, kind, fullPath string, currentPerm currentPermissions, defaultPerm defaultPermissions, verbose bool) {
	if verbose {
		fullPath = filepath.Base(fullPath)
	}
	hierarchy := strings.Repeat("  ", level)

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

// Colorize returns a string using ansi colors
func Colorize(color, s string) string {
	const (
		esc        = "\x1b"
		ansiBlue   = esc + "[34;1m"
		ansiRed    = esc + "[31;1m"
		ansiYellow = esc + "[33;1m"
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
