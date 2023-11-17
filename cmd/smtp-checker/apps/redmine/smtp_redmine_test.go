package redmine

import (
	"log"
	"os"
	"testing"

	"github.com/bitnami/healthcheck-tools/cmd/smtp-checker/apps"
)

var testRedmineConfig = `
# = Redmine configuration file
#
default:
  email_delivery:
    delivery_method: :smtp
    smtp_settings:
      address: "smtp.gmail.com"
      port: 587
      domain: "smtp.gmail.com"
      authentication: :plain
      enable_starttls_auto: true
      user_name: "user@gmail.com"
      password: "XXXXXXXX"
# Production
production:
# Development
development:
`

func TestParseConfig(t *testing.T) {
	t.Run("Check parsed database configuration data", func(t *testing.T) {
		config := Config{}
		tmpConfigFile := createTemporaryFile(testRedmineConfig, "configuration.yaml")
		defer os.Remove(tmpConfigFile.Name())
		err := apps.UnmarshalYAMLFile(tmpConfigFile.Name(), &config)
		if err != nil {
			t.Errorf("Error unmarshaling YAML file: %v", err)
		}
		err = config.ValidateSMTPSettings()
		if err != nil {
			t.Errorf("Error validating SMTP data: %v", err)
		}
	})
}

func createTemporaryFile(content, prefix string) *os.File {
	tmpFile, err := os.CreateTemp("", prefix)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := tmpFile.Write([]byte(content)); err != nil {
		log.Fatal(err)
	}
	return tmpFile
}
