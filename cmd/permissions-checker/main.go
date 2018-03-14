package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

var (
	// VERSION will be overwritten automatically by the build system
	VERSION = "devel"
)

func main() {
	var (
		installDir string
		app        string
		getVersion bool
	)
	flag.StringVar(&installDir, "install_dir", "/opt/bitnami", "Installation Directory")
	flag.StringVar(&app, "application", "", "Application")
	flag.BoolVar(&getVersion, "version", false, "Show current version")
	flag.Parse()

	if getVersion {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	fmt.Printf(`==================================================
	PERMISSIONS CHECKS
==================================================
Starting checks with these parameters:
	- Application: %s
	- Installation directory: %s
==================================================

`, app, installDir)

	c1 := exec.Command("ls", "-lR", installDir)
	c2 := exec.Command("awk", "FNR > 1 {k=0;for(i=0;i<=8;i++)k+=((substr($1,i+2,1)~/[rwx]/)*2^(8-i));if(k)printf(\"%0o \",k);print $1\" \"$3\" \"$4\" \"$9}")
	c2.Stdin, _ = c1.StdoutPipe()
	c2.Stdout = os.Stdout
	_ = c2.Start()
	_ = c1.Run()
	_ = c2.Wait()
}
