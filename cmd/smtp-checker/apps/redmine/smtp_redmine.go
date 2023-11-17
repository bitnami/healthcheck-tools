package redmine

import (
	"path/filepath"

	"github.com/bitnami/healthcheck-tools/cmd/smtp-checker/apps"
	"github.com/juju/errors"
)

const (
	configFilePath = "apps/redmine/htdocs/config/configuration.yml"
)

// Config is a structure that matches the schema of
// Redmine config/configuration.yml file
type Config struct {
	Default mode `json:"default"`
}

type mode struct {
	EmailDelivery emailDelivery `json:"email_delivery"`
}

type emailDelivery struct {
	DeliveryMethod string       `json:"delivery_method"`
	SMTPSettings   SMTPSettings `json:"smtp_settings"`
}

// SMTPSettings is a structure that contains the SMTP
// configuration data of Redmine
type SMTPSettings struct {
	Address        string `json:"address"`
	Port           int    `json:"port"`
	Username       string `json:"user_name"`
	Password       string `json:"password"`
	Domain         string `json:"domain"`
	Authentication string `json:"authentication"`
	AutoStartTLS   bool   `json:"enable_starttls_auto"`
}

// GetSMTPSettings returns a SMTPSettings from Config structure
func (c Config) GetSMTPSettings() *apps.SMTPSettings {
	return &apps.SMTPSettings{
		Host: c.Default.EmailDelivery.SMTPSettings.Domain,
		Port: c.Default.EmailDelivery.SMTPSettings.Port,
		User: c.Default.EmailDelivery.SMTPSettings.Username,
		Pass: c.Default.EmailDelivery.SMTPSettings.Password,
	}
}

// ValidateSMTPSettings checks the SMTPSettings are correct
func (c *Config) ValidateSMTPSettings() error {
	if c.Default.EmailDelivery.SMTPSettings.Address == "" {
		return errors.New("address: empty string")
	}
	if c.Default.EmailDelivery.SMTPSettings.Domain == "" {
		return errors.New("domain: empty string")
	}
	if c.Default.EmailDelivery.SMTPSettings.Port == 0 {
		return errors.New("port: invalid port")
	}
	if c.Default.EmailDelivery.SMTPSettings.Username == "" {
		return errors.New("user_name: empty string")
	}
	if c.Default.EmailDelivery.SMTPSettings.Password == "" {
		return errors.New("password: empty string")
	}
	if c.Default.EmailDelivery.SMTPSettings.Domain != c.Default.EmailDelivery.SMTPSettings.Address {
		return errors.Errorf("address %s does not match domain %s on smtp_settings", c.Default.EmailDelivery.SMTPSettings.Address, c.Default.EmailDelivery.SMTPSettings.Domain)
	}
	if c.Default.EmailDelivery.SMTPSettings.Domain == "smtp.gmail.com" && !c.Default.EmailDelivery.SMTPSettings.AutoStartTLS {
		return errors.New("use enable_starttls_auto when using smtp.gmail.com as domain")
	}
	return nil
}

// ParseConfig obtains an ApplicationConfig from by parsing a config file
func ParseConfig(installDir string) (apps.ApplicationConfig, error) {
	config := Config{}
	return &config, apps.UnmarshalYAMLFile(filepath.Join(installDir, configFilePath), &config)
}
