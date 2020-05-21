package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bitnami-labs/healthcheck-tools/cmd/smtp-checker/apps"
	"github.com/mmikulicic/multierror"
)

const defaultRecipient = "test@example.com"

var (
	// VERSION will be overwritten automatically by the build system
	VERSION = "devel"
)

func main() {
	var (
		installDir     string
		app            string
		recipient      string
		getVersion     bool
		secureOutput   bool
		passwordOutput string
	)
	flag.StringVar(&installDir, "install_dir", "/opt/bitnami", "Installation Directory")
	flag.StringVar(&app, "application", "", "Application")
	flag.StringVar(&recipient, "mail_recipient", defaultRecipient, fmt.Sprintf("Mail Recipient (%s by default)", defaultRecipient))
	flag.BoolVar(&getVersion, "version", false, "Show current version")
	flag.BoolVar(&secureOutput, "secure_output", false, "Hide SMTTP password in output")
	smtp := apps.NewSMTPSettingsFromFlags(flag.CommandLine)
	flag.Parse()

	if getVersion {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	if app != "" {
		fmt.Printf(`======================================
SMTP CONFIGURATION
======================================
Obtaining SMTP configuration for app: %q
  - Installation Directory: %q

`, app, installDir)

		appConfig, err := ObtainConfigData(installDir, app)
		if err != nil {
			log.Fatalf("Found errors when obtaining the SMTP configuration: %q", err)
		}
		err = appConfig.ValidateSMTPSettings()
		if err != nil {
			log.Fatalf("Found errors when validating the SMTP settings: %q", err)
		}
		smtp = appConfig.GetSMTPSettings()
		fmt.Println("SMTP configuration successfully retrieved!!")
	}

	if smtp.Host == "" || smtp.Port == 0 || smtp.User == "" || smtp.Pass == "" {
		log.Fatalf("Indicate your application using '-application' flag or set the smtp credentials using 'smtp-host', 'smtp-port', '-smtp-user' and '-smtp-password' flags")
	}

	recipientText := recipient
	if recipient == defaultRecipient {
		recipientText = fmt.Sprintf("%s (invalid mail account, use -mail_recipient lag to indicate a valid one)", defaultRecipient)
	}

	if secureOutput {
		passwordOutput = "xxxxxx"
	} else {
		passwordOutput = smtp.Pass
	}

	fmt.Printf(`
======================================
SMTP CHECKS
======================================
Using SMTP credentials:
  - SMTP Host: %q
  - SMTP Port: %d
  - SMTP User: %q
  - SMTP Password: %q
  - Mail Recipient: %q

`, smtp.Host, smtp.Port, smtp.User, passwordOutput, recipientText)

	var errors error

	fmt.Println("-- Check: Connectivity with SMTP server --")
	err := RunConnectivityChecks(smtp.Host, smtp.Port)
	if err != nil {
		errors = multierror.Append(errors, err)
	}

	if smtp.Port == 465 {
		fmt.Println("-- Check: Connectivity with SMTP server via TLS --")
		err = RunTLSConnectivityChecks(smtp.Host, smtp.Port)
		if err != nil {
			errors = multierror.Append(errors, err)
		}
	}

	fmt.Println("-- Check: server time offset --")
	err = RunNTPChecks()
	if err != nil {
		errors = multierror.Append(errors, err)
	}

	fmt.Println("-- Check: Send mail via SMTP --")
	if recipient != defaultRecipient {
		fmt.Printf("\nNote: Remember to check the recipient's mail inbox!\n")
	}
	err = RunSendMailChecks(smtp, recipient)
	if err != nil {
		errors = multierror.Append(errors, err)
	}

	fmt.Printf(`
======================================
SMTP CHECKS FINISHED
======================================

`)
	if errors != nil {
		log.Fatalf("Found errors when checking the SMTP configuration:\n%v", errors)
	}
}
