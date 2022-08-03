package apps

import (
	"flag"
	"fmt"
	"os"

	"github.com/ghodss/yaml"
)

// Application is a structure that contains the info
// about each supported application and its configuration file.
type Application struct {
	Name       string
	ConfigFile string
}

// SMTPSettings is a structure that contains the SMTP
// credentials to use on the SMTP checks
type SMTPSettings struct {
	Host string `default:"localhost"`
	Port int    `default:"25"`
	User string
	Pass string
}

// NewSMTPSettingsFromFlags creates a SMTPSettings from the provided command line flags
func NewSMTPSettingsFromFlags(fs *flag.FlagSet) *SMTPSettings {
	smtp := SMTPSettings{}
	flag.StringVar(&smtp.Host, "smtp_host", "localhost", "SMTP Host")
	flag.IntVar(&smtp.Port, "smtp_port", 25, "SMTP Port")
	flag.StringVar(&smtp.User, "smtp_user", "", "SMTP User")
	flag.StringVar(&smtp.Pass, "smtp_password", "", "SMTP Password")
	return &smtp
}

// ApplicationConfig defines a common interface for getting information from the application config
type ApplicationConfig interface {
	GetSMTPSettings() *SMTPSettings
	ValidateSMTPSettings() error
}

// UnmarshalYAMLFile reads a config file and unmarshals it into a config struct
func UnmarshalYAMLFile(configFile string, config interface{}) error {
	source, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}
	fmt.Printf("Reading config file: %q\n", configFile)
	if err := yaml.Unmarshal(source, config); err != nil {
		return fmt.Errorf("error parsing config file: %v", err)
	}
	return nil
}
