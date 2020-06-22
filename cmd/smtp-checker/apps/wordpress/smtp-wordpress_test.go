package wordpress

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var testWPConfig = `
<?php
/**
 * The base configuration for WordPress
 */

// ** MySQL settings - You can get this info from your web host ** //
/** The name of the database for WordPress */
define('DB_NAME', 'bitnami_wordpress');

/** MySQL database username */
define('DB_USER', 'bn_wordpress');

/** MySQL database password */
define('DB_PASSWORD', 'XXXXXXXX');

/** MySQL hostname */
define('DB_HOST', 'localhost:3306');

/** Database Charset to use in creating database tables. */
define('DB_CHARSET', 'utf8');

/** The Database Collate type. Don't change this if in doubt. */
define('DB_COLLATE', '');
`

var pluginsData = `a:1:{i:0;s:29:"wp-mail-smtp/wp_mail_smtp.php";}`
var smtpData = `a:5:{s:4:"mail";a:4:{s:10:"from_email";s:16:"user@example.com";s:9:"from_name";s:13:"user\'s Blog!";s:6:"mailer";s:4:"smtp";s:11:"return_path";b:0;}s:4:"smtp";a:7:{s:7:"autotls";b:1;s:4:"host";s:14:"smtp.gmail.com";s:4:"port";i:587;s:10:"encryption";s:4:"none";s:4:"user";s:9:"smtp-user";s:4:"pass";s:8:"XXXXXXXX";s:4:"auth";b:1;}s:5:"gmail";a:2:{s:9:"client_id";s:0:"";s:13:"client_secret";s:0:"";}s:7:"mailgun";a:2:{s:7:"api_key";s:0:"";s:6:"domain";s:0:"";}s:8:"sendgrid";a:1:{s:7:"api_key";s:0:"";}}`

func TestParseWPDatabaseConfig(t *testing.T) {
	t.Run("Check parsed database configuration data", func(t *testing.T) {
		tmpConfigFile := createTemporaryFile(testWPConfig, "wp-config.php")
		defer os.Remove(tmpConfigFile.Name())
		database, err := parseWPDatabaseConfig(tmpConfigFile.Name())
		if err != nil {
			t.Errorf("Error parsing wp-config.php file: %v", err)
		}
		if database.Host != "localhost" {
			t.Errorf("Incorrect database host detected, expected: localhost, got: %s", database.Host)
		}
		if database.Port != 3306 {
			t.Errorf("Incorrect database port detected, expected: 3306, got: %d", database.Port)
		}
		if database.Name != "bitnami_wordpress" {
			t.Errorf("Incorrect database name detected, expected: bitnami_wordpress, got: %s", database.Name)
		}
		if database.User != "bn_wordpress" {
			t.Errorf("Incorrect database user detected, expected: bn_wordpress, got: %s", database.User)
		}
		if database.Pass != "XXXXXXXX" {
			t.Errorf("Incorrect database password detected, expected: XXXXXXXX, got: %s", database.Pass)
		}
	})
}

func TestCheckPlugin(t *testing.T) {
	t.Run("Check installed plugins", func(t *testing.T) {
		err := checkPlugin(pluginsData)
		if err != nil {
			t.Errorf("Error checking installed plugins: %v", err)
		}
	})
}

func TestObtainSMTP(t *testing.T) {
	t.Run("Check obtained SMTP data", func(t *testing.T) {
		config := Config{}
		err := obtainSMTP(smtpData, &config)
		if err != nil {
			t.Errorf("Error obtaining SMTP data: %v", err)
		}
		err = config.ValidateSMTPSettings()
		if err != nil {
			t.Errorf("Error validating SMTP data: %v", err)
		}
		smtp := config.GetSMTPSettings()
		if smtp.Host != "smtp.gmail.com" {
			t.Errorf("Incorrect SMTP host detected, expected: smtp.gmail.com, got: %s", smtp.Host)
		}
		if smtp.Port != 587 {
			t.Errorf("Incorrect SMTP port detected, expected: 587, got: %d", smtp.Port)
		}
		if smtp.Pass != "XXXXXXXX" {
			t.Errorf("Incorrect SMTP password detected, expected: XXXXXXXX, got: %s", smtp.Pass)
		}
		if smtp.User != "smtp-user" {
			t.Errorf("Incorrect SMTP user detected, expected: smtp-user, got: %s", smtp.User)
		}
	})
}

func createTemporaryFile(content, prefix string) *os.File {
	tmpFile, err := ioutil.TempFile("", prefix)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := tmpFile.Write([]byte(content)); err != nil {
		log.Fatal(err)
	}
	return tmpFile
}
