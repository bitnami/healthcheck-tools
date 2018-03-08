package wordpress

import (
	"fmt"
	"github.com/bitnami-labs/healthcheck-tools/cmd/smtp-checker/apps"
	"github.com/bitnami-labs/healthcheck-tools/pkg/mysql"
	"github.com/juju/errors"
	"github.com/yvasiyarov/php_session_decoder/php_serialize"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	propValue *regexp.Regexp
)

const (
	configFilePath = "apps/wordpress/htdocs/wp-config.php"
)

// Config is a structure that contains the Mail/SMTP
// configuration data of WP
type Config struct {
	Mail         Mail
	SMTPSettings SMTPSettings
}

// Mail is a structure that contains the Mail
// configuration data of WP
type Mail struct {
	FromMail string
	FromName string
	Mailer   string
}

// SMTPSettings is a structure that contains the SMTP
// configuration data of WP
type SMTPSettings struct {
	Host    string
	Port    int
	User    string
	Pass    string
	Encrypt string
	Auth    bool
	AutoTLS bool
}

// GetSMTPSettings returns a SMTPSettings from Config structure
func (c Config) GetSMTPSettings() *apps.SMTPSettings {
	return &apps.SMTPSettings{
		Host: c.SMTPSettings.Host,
		Port: c.SMTPSettings.Port,
		User: c.SMTPSettings.User,
		Pass: c.SMTPSettings.Pass,
	}
}

// ValidateSMTPSettings checks the SMTPSettings are correct
func (c *Config) ValidateSMTPSettings() error {
	if c.SMTPSettings.Host == "" {
		return errors.New("host: empty string")
	}
	if c.SMTPSettings.Port == 0 {
		return errors.New("port: invalid port")
	}
	if c.SMTPSettings.User == "" {
		return errors.New("user: empty string")
	}
	if c.SMTPSettings.Pass == "" {
		return errors.New("password: empty string")
	}
	if c.SMTPSettings.Host == "smtp.gmail.com" && !c.SMTPSettings.AutoTLS {
		return errors.New("set autotls to true when using smtp.gmail.com as domain")
	}
	return nil
}

func getProperty(buffer []byte, property string) string {
	propValue = regexp.MustCompilePOSIX(fmt.Sprintf(`^[[:space:]]*define\('%s',[[:space:]]+'(.*)'\);`, property))
	return propValue.FindStringSubmatch(string(buffer[:]))[1]
}

func parseWPDatabaseConfig(configFile string) (mysql.Database, error) {
	source, err := ioutil.ReadFile(configFile)
	if err != nil {
		return mysql.Database{}, errors.Errorf("error reading config file: %v", err)
	}
	port, _ := strconv.Atoi(strings.Split(getProperty(source, "DB_HOST"), ":")[1])
	return mysql.Database{
		Host: strings.Split(getProperty(source, "DB_HOST"), ":")[0],
		Port: port,
		User: getProperty(source, "DB_USER"),
		Name: getProperty(source, "DB_NAME"),
		Pass: getProperty(source, "DB_PASSWORD"),
	}, nil
}

func checkPlugin(pluginsInfo string) error {
	val, err := php_serialize.NewUnSerializer(pluginsInfo).Decode()
	if err != nil {
		return errors.Errorf("error while decoding object value: %v", err)
	}
	plugins, ok := val.(php_serialize.PhpArray)
	if !ok {
		return errors.Errorf("unable to convert %v to PhpArray", val)
	}
	for _, v := range plugins {
		plugin, ok := v.(string)
		if !ok {
			return errors.Errorf("unable to convert %v to string", v)
		}
		if plugin == "wp-mail-smtp/wp_mail_smtp.php" {
			return nil
		}
	}
	return errors.New("wp-mail-smtp plugin not installed")
}

func checkPluginOnDatabase(database mysql.Database) error {
	query := mysql.Query{
		Table:  "wp_options",
		Column: "option_value",
		Key:    "option_name",
		Value:  "active_plugins",
	}
	queryResult, err := database.MySQLQuery(query)
	if err != nil {
		return err
	}
	return checkPlugin(queryResult)
}

func obtainSMTP(smtpInfo string, config *Config) error {
	val, err := php_serialize.NewUnSerializer(smtpInfo).Decode()
	if err != nil {
		return errors.Errorf("error while decoding object value: %v", err)
	}
	arrays, ok := val.(php_serialize.PhpArray)
	if !ok {
		return errors.Errorf("unable to convert %v to PhpArray", val)
	}
	for key, array := range arrays {
		switch key {
		case "smtp":
			settings, ok := array.(php_serialize.PhpArray)
			if !ok {
				return errors.Errorf("unable to convert %v to PhpArray", array)
			}
			for s, v := range settings {
				switch s {
				case "host":
					config.SMTPSettings.Host, ok = v.(string)
					if !ok {
						return errors.Errorf("unable to convert %v to string", v)
					}
				case "port":
					config.SMTPSettings.Port, ok = v.(int)
					if !ok {
						return errors.Errorf("unable to convert %v to int", v)
					}
				case "encyrption":
					config.SMTPSettings.Encrypt, ok = v.(string)
					if !ok {
						return errors.Errorf("unable to convert %v to string", v)
					}
				case "user":
					config.SMTPSettings.User, ok = v.(string)
					if !ok {
						return errors.Errorf("unable to convert %v to string", v)
					}
				case "pass":
					config.SMTPSettings.Pass, ok = v.(string)
					if !ok {
						return errors.Errorf("unable to convert %v to string", v)
					}
				case "auth":
					config.SMTPSettings.Auth, ok = v.(bool)
					if !ok {
						return errors.Errorf("unable to convert %v to bool", v)
					}
				case "autotls":
					config.SMTPSettings.AutoTLS, ok = v.(bool)
					if !ok {
						return errors.Errorf("unable to convert %v to bool", v)
					}
				}
			}
		case "mail":
			settings, ok := array.(php_serialize.PhpArray)
			if !ok {
				return errors.Errorf("unable to convert %v to PhpArray", array)
			}
			for s, v := range settings {
				switch s {
				case "from_email":
					config.Mail.FromMail, ok = v.(string)
					if !ok {
						return errors.Errorf("unable to convert %v to string", v)
					}
				case "from_name":
					config.Mail.FromName, ok = v.(string)
					if !ok {
						return errors.Errorf("unable to convert %v to string", v)
					}
				case "mailer":
					config.Mail.Mailer, ok = v.(string)
					if !ok {
						return errors.Errorf("unable to convert %v to string", v)
					}
				}
			}
		}
	}
	return nil
}

func obtainSMTPFromDatabase(database mysql.Database, config *Config) error {
	query := mysql.Query{
		Table:  "wp_options",
		Column: "option_value",
		Key:    "option_name",
		Value:  "wp_mail_smtp",
	}
	queryResult, err := database.MySQLQuery(query)
	if err != nil {
		return err
	}
	return obtainSMTP(queryResult, config)
}

// QueryConfig obtains an ApplicationConfig from by querying the MySQL database
func QueryConfig(installDir string) (apps.ApplicationConfig, error) {
	config := Config{}
	database, err := parseWPDatabaseConfig(filepath.Join(installDir, configFilePath))
	if err != nil {
		return nil, errors.Errorf("error parsing wp-config.php file: %v", err)
	}
	err = checkPluginOnDatabase(database)
	if err != nil {
		return nil, errors.Errorf("error checking wp-mail-smtp plugin: %v", err)
	}
	return &config, obtainSMTPFromDatabase(database, &config)
}
