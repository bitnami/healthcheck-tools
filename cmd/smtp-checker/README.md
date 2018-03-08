# SMTP Health Checker
_smtp-checker_ performs a set of health checks in your web application to detect any possible issues with the SMTP configuration.

## Installation

```
$> go get github.com/bitnami-labs/healthcheck-tools/cmd/smtp-checker
```

## Building from source

```
$> git clone https://github.com/bitnami-labs/healthcheck-tools.git
$> cd cmd/smtp-checker
$> make
```

## Basic usage

The tool is executed as follows:

```
$> smtp-checker -application <APPLICATION> -install_dir <STACK INSTALLATION DIRECTORY> -smtp_host <SMTP HOST> -smtp_port <SMTP PORT> -smtp_user <SMTP USER> -smtp_password <SMTP PASSWORD>
```

The tool requires a set of parameters to work properly:

  - *application*: Application used (e.g wordpress). Parameter required.
  - *install_dir*: Stack installation directory. Default value: */opt/bitnami*.

Or:

  - *smtp_host*: SMTP server hostname. Parameter required if application not provided.
  - *smtp_port*: SMTP server port. Parameter required if application not provided.
  - *smtp_user*: SMTP user. Parameter required if application not provided.
  - *smtp_password*: SMTP user's password. Parameter required if application not provided.

Optional parameters.

  - *mail_recipient*: Mail recipient for sending testing mails via SMTP.  Default value: *test@example.com*.
