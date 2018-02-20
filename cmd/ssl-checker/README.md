# SSL Health Checker
_ssl-checker_ performs a set of health checks in your web application to detect any possible issues with the SSL configuration.

## Installation

```
$> go get github.com/bitnami-labs/healthcheck-tools/cmd/ssl-checker/...
```

## Building from source

```
$> git clone https://github.com/bitnami-labs/healthcheck-tools.git
$> cd cmd/ssl-checker
$> make 
```

## Requirements

This tool expects _Apache_ as the web server (_nginx_ will be supported in later versions).

## Basic usage

The tool is executed as follows:

```
$> ssl-checker -apache-root <APACHE FOLDER> -apache-conf <APACHE CONF FILE> -hostname <SERVER IP/HOSTNAME> -port <HTTPS PORT>
```

The tool requires a set of parameters to work properly:

  - *apache-root*: Directory where apache is installed. Default value: */opt/bitnami/apache2*.
  - *apache-conf*: Apache configuration file. Default value: */opt/bitnami/apache/conf/httpd.conf*.
  - *hostname*: Hostname or IP address where the web server is running. Parameter required.
  - *port*: Port where the web server is serving HTTPS requests. Default value: 443 
