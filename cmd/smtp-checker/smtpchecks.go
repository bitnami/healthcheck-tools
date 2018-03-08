package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"math"
	"net"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/andybalholm/crlf"
	"github.com/beevik/ntp"
	"github.com/bitnami-labs/healthcheck-tools/cmd/smtp-checker/apps"
	"github.com/bitnami-labs/healthcheck-tools/cmd/smtp-checker/apps/redmine"
	"github.com/bitnami-labs/healthcheck-tools/cmd/smtp-checker/apps/wordpress"
	"github.com/juju/errors"
)

const (
	timeout        = 10 * time.Second
	maxClockOffset = 1 * time.Second
)

func absDuration(d time.Duration) time.Duration {
	return time.Duration(math.Abs(float64(d)))
}

// ObtainConfigData obtains the configuration data from
// the app
func ObtainConfigData(installDir string, app string) (appConfig apps.ApplicationConfig, err error) {
	parsers := map[string]func(string) (apps.ApplicationConfig, error){
		"redmine":   redmine.ParseConfig,
		"wordpress": wordpress.QueryConfig,
	}
	parse, ok := parsers[app]
	if !ok {
		var names []string
		for k := range parsers {
			names = append(names, k)
		}
		return nil, errors.Errorf("bad app name %q; currently supported: %s", app, strings.Join(names, ", "))
	}
	return parse(installDir)
}

// RunConnectiviyChecks performs checks on the connectivity
// with SMTP server
func RunConnectivityChecks(hostname string, port int) error {
	smtpServer := fmt.Sprintf("%s:%d", hostname, port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", smtpServer)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	conn.Close()
	fmt.Println("Succesful connectivity!")
	return nil
}

// RunTLSConnectiviyChecks performs checks on the connectivity with SMTP server
func RunTLSConnectivityChecks(hostname string, port int) error {
	smtpServer := fmt.Sprintf("%s:%d", hostname, port)
	conn, err := tls.Dial("tcp", smtpServer, nil)
	if err != nil {
		return err
	}
	conn.Close()
	fmt.Println("Succesful TLS connectivity!")
	return nil
}

// RunNTPChecks performs checks on the Time offset respect a NTP pool
func RunNTPChecks() error {
	rp, err := ntp.QueryWithOptions("pool.ntp.org", ntp.QueryOptions{Timeout: timeout})
	if err != nil {
		return err
	}
	if absDuration(rp.ClockOffset) > maxClockOffset {
		return errors.Errorf("incorrect time offset (>%s, synchronize your server clock via ntp", maxClockOffset)
	}
	h, err := os.Hostname()
	if err != nil {
		h = "localhost"
	}
	fmt.Printf("Time synchronisation of host %s within reasonable bounds!\n", h)
	return nil
}

// RunSendMailChecks performs checks on sending mails via SMTP
func RunSendMailChecks(settings *apps.SMTPSettings, recipient string) error {
	auth := smtp.PlainAuth(
		"",
		settings.User,
		settings.Pass,
		settings.Host,
	)
	sender := settings.User
	smtpServer := fmt.Sprintf("%s:%d", settings.Host, settings.Port)
	var msg bytes.Buffer
	w := crlf.NewWriter(&msg)
	fmt.Fprintf(w, "To: %s\n", recipient)
	fmt.Fprintf(w, `Subject: Testing Mail

This is a testing email body.`)
	if err := smtp.SendMail(smtpServer, auth, sender, []string{recipient}, msg.Bytes()); err != nil {
		return err
	}
	fmt.Println("Mail successfully sent via SMTP!")
	return nil
}
