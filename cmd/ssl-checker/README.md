# SSL Health Checker
_ssl-checker_ performs a set of health checks in your web application to detect any possible issues with the SSL configuration.

## Installation

```
$> go get github.com/bitnami-labs/healthcheck-tools/cmd/ssl-checker
```

## Building from source

```
$> git clone https://github.com/bitnami-labs/healthcheck-tools.git
$> make ssl-checker-linux-amd64
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

## List of health checks
The tool will perform the following health checks:

  - Check if the Apache configuration contains SSL certificate-key pairs. It will show where these are defined. 
  - Check if the detected certificates are not corrupted.
  - Check the domain name of the certificates.
  - Check if the certificate-key pairs match.
  - Check the certificate that the web server is returning (this requires you to have a running web server).
  
  ## Useful links
  
  - [Troubleshoot SSL issues (Bitnami Documentation pages)](https://docs.bitnami.com/general/how-to/troubleshoot-ssl-issues/).
