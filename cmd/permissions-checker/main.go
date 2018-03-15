package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	// VERSION will be overwritten automatically by the build system
	VERSION = "devel"
	// ShowHidden include hidden files or directories in the loops
	ShowHidden bool
)

func main() {
	var (
		installDir string
		getVersion bool
	)
	flag.StringVar(&installDir, "install_dir", "/opt/bitnami", "Installation Directory")
	flag.BoolVar(&getVersion, "version", false, "Show current version")
	flag.BoolVar(&ShowHidden, "show_hidden", false, "Show hidden files and directories")
	flag.Parse()

	if getVersion {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	fmt.Printf(`==================================================
	PERMISSIONS CHECKS
==================================================
Starting checks with these parameters:
	- Installation directory: %s
	- Show hidden: %t
==================================================

`, installDir, ShowHidden)

	files, err := ioutil.ReadDir(installDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		switch mode := f.Mode(); {
		case mode.IsRegular():
			fileName := f.Name()
			filePerm := mode.Perm()
			if f.Name()[0:1] != "." || ShowHidden {
				PrintPermissions("", "f", fileName, filePerm)
			}
		case mode.IsDir():
			dirName := f.Name()
			dirPerm := mode.Perm()
			if f.Name()[0:1] != "." || ShowHidden {
				PrintPermissions("", "d", dirName, dirPerm)
				FindRecursive(strings.Join([]string{installDir, dirName}, "/"), "")
			}
		}
	}
}
