package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// FindRecursive iterate in the directory showing permissions in a recursive way
func FindRecursive(currentPath, level string) {
	files, err := ioutil.ReadDir(currentPath)
	if err != nil {
		log.Fatal(err)
	}
	level = level + "  "
	for _, f := range files {
		switch mode := f.Mode(); {
		case mode.IsRegular():
			fileName := f.Name()
			filePerm := mode.Perm()
			if f.Name()[0:1] != "." || ShowHidden {
				PrintPermissions(level, "f", fileName, filePerm)
			}
		case mode.IsDir():
			dirName := f.Name()
			dirPerm := mode.Perm()
			if f.Name()[0:1] != "." || ShowHidden {
				PrintPermissions(level, "d", dirName, dirPerm)
				FindRecursive(strings.Join([]string{currentPath, dirName}, "/"), level)
			}
		}
	}
}

// PrintPermissions print permissions
func PrintPermissions(level, kind, name string, permissions os.FileMode) {
	fmt.Printf("%s(%s) %s %s\n", level, kind, name, permissions)
}
