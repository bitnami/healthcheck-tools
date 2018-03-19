package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
	"strings"
	"syscall"
)

// FindRecursive iterate in a directory (path) showing permissions in a recursive way
func FindRecursive(path, level, fileDefPermissions, dirDefPermissions string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	level = level + "  "
	for _, f := range files {
		name := f.Name()
		fullPath := strings.Join([]string{path, name}, "/")
		owner, _ := user.LookupId(fmt.Sprint(f.Sys().(*syscall.Stat_t).Uid))
		group, _ := user.LookupGroupId(fmt.Sprint(f.Sys().(*syscall.Stat_t).Gid))
		switch mode := f.Mode(); {
		case mode.IsRegular():
			filePerm := mode.Perm().String()
			if f.Name()[0:1] != "." || ShowHidden {
				CheckPermissions(level, "f", name, fullPath, filePerm, fileDefPermissions, dirDefPermissions, owner.Username, group.Name)
			}
		case mode.IsDir():
			dirPerm := mode.Perm().String()
			if f.Name()[0:1] != "." || ShowHidden {
				CheckPermissions(level, "d", name, fullPath, dirPerm, fileDefPermissions, dirDefPermissions, owner.Username, group.Name)
				FindRecursive(fullPath, level, fileDefPermissions, dirDefPermissions)
			}
		}
	}
}

// CheckPermissions check if permissions are expected
func CheckPermissions(level, kind, name, path, permissions, fileDefPermissions, dirDefPermissions, owner, group string) {
	var defaultPermissions string
	if kind == "f" {
		defaultPermissions = fileDefPermissions
	} else if kind == "d" {
		defaultPermissions = dirDefPermissions
	}

	if strings.Contains(permissions, defaultPermissions) {
		if Verbose {
			fmt.Printf("%s(%s) %s %s %s %s\n", level, kind, name, permissions, owner, group)
		}
	} else {
		if !Verbose {
			fmt.Printf("\x1b[31;1m(%s) %s %s (expected %s) %s %s\n\x1b[0m", kind, path, permissions, defaultPermissions, owner, group)
		} else {
			fmt.Printf("\x1b[31;1m%s(%s) %s %s (expected %s) %s %s\n\x1b[0m", level, kind, name, permissions, defaultPermissions, owner, group)
		}
	}
}
