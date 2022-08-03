package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// These variables will be overwritten automatically by the build system
var VERSION = "devel"
var BUILD_DATE = ""
var COMMIT = ""

func main() {
	var apacheRoot string
	var apacheConf string
	var hostname string
	var port int
	var getVersion bool
	flag.StringVar(&apacheRoot, "apache-root", "/opt/bitnami/apache2/", "Root of Apache installation")
	flag.StringVar(&apacheConf, "apache-conf", "/opt/bitnami/apache2/conf/httpd.conf",
		"Path to the root Apache configuration file")
	flag.StringVar(&hostname, "hostname", "", "Web application hostname")
	flag.IntVar(&port, "port", 443, "Web application port")
	flag.BoolVar(&getVersion, "version", false, "Show current version")
	flag.Parse()
	if getVersion {
		fmt.Printf("smtp-checker %s\n", VERSION)
		fmt.Printf("Built on: %s\n", BUILD_DATE)
		fmt.Printf("Commit: %s\n", COMMIT)
		os.Exit(0)

		os.Exit(0)
	}
	if hostname == "" {
		log.Fatal("-hostname flag must be set")
	}
	fmt.Printf(`======================================
SSL CHECKS
======================================
Starting checks with these parameters:
  - Apache Root: %q
  - Apache Root configuration: %q
  - Hostname: %q
  - Port: %d
======================================
`, apacheRoot, apacheConf, hostname, port)

	fmt.Println("-- Check: Active SSL Certificates in Apache Configuration --")
	err := RunActiveCertificatesChecks(apacheConf, apacheRoot)
	foundErrors := false
	if err != nil {
		fmt.Fprintf(os.Stderr, "Active Certificate check failed: %q\n", err)
		foundErrors = true
	}
	fmt.Printf("-- End of check --\n\n")

	fmt.Println("-- Check: HTTPS Connection to web server --")
	err = RunHTTPSConnectionChecks(hostname, port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTPS Connection failed: %q\n", err)
		foundErrors = true
	}
	fmt.Printf("-- End of check --\n\n")
	fmt.Println("SSL Checks finished")
	if foundErrors {
		log.Fatalf("Found errors when checking the SSL configuration")
	} else {
		os.Exit(0)
	}
}
