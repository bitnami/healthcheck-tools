package main

import (
	"os"
	"strconv"
	"testing"

	"github.com/bitnami/healthcheck-tools/cmd/smtp-checker/apps"
)

func TestRunTLSConnectivityChecks(t *testing.T) {
	t.Run("Check connectivity with SMTP server via TLS", func(t *testing.T) {
		err := RunTLSConnectivityChecks("smtp.gmail.com", 465)
		if err != nil {
			t.Errorf("error connecting to smtp server via tls: %v", err)
		}
	})
}

func TestRunNTPChecks(t *testing.T) {
	t.Run("Check time offset via NTP", func(t *testing.T) {
		err := RunNTPChecks()
		if err != nil {
			t.Errorf("error checking time offset via NTP: %v", err)
		}
	})
}

func TestRunSendMailChecks(t *testing.T) {
	t.Run("Check mail delivery via SMTP", func(t *testing.T) {
		port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
		smtp := apps.SMTPSettings{
			Host: os.Getenv("SMTP_HOST"),
			Port: port,
			User: os.Getenv("SMTP_USER"),
			Pass: os.Getenv("SMTP_PASS"),
		}
		err := RunSendMailChecks(&smtp, "test@example.com")
		if err != nil {
			t.Errorf("error checking mail delivery via SMTP: %v", err)
		}
	})
}
