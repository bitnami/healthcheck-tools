package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"path"
	"regexp"

	"github.com/bitnami-labs/healthcheck-tools/pkg/apache"
)

// CertificatePairInfo contains paths of an active certificate-key path
type CertificatePairInfo struct {
	apacheConfPath string
	certPath       string
	keyPath        string
}

// getActiveCertificatePairsInAllFiles obtains the certificate-key pairs that are being used in a single file
func getActiveCertificatePairs(file, text, apacheRoot string) []CertificatePairInfo {
	res := []CertificatePairInfo{}
	sslcertRe := regexp.MustCompilePOSIX(`^[[:space:]]*SSLCertificateFile[[:space:]]+["]?([^\n"]+)["]?`)
	sslcertMatches := sslcertRe.FindAllStringSubmatch(text, -1)

	sslcertKeyRe := regexp.MustCompilePOSIX(`^[[:space:]]*SSLCertificateKeyFile[[:space:]]+["]?([^\n"]+)["]?`)
	sslcertKeyMatches := sslcertKeyRe.FindAllStringSubmatch(text, -1)

	for index, element := range sslcertMatches {
		newCert := element[1]
		if !path.IsAbs(newCert) {
			newCert = path.Join(apacheRoot, newCert)
		}
		newKey := ""
		if len(sslcertKeyMatches) > index {
			newKey = sslcertKeyMatches[index][1]
			if !path.IsAbs(newKey) {
				newKey = path.Join(apacheRoot, newKey)
			}
			res = append(res, CertificatePairInfo{file, newCert, newKey})
		}
	}
	return res

}

// getActiveCertificatePairsInAllFiles obtains the certificate-key pairs that are being used in each Apache configuration file
func getActiveCertificatePairsInAllFiles(apacheConf map[string]string, apacheRoot string) []CertificatePairInfo {
	res := []CertificatePairInfo{}
	for file, content := range apacheConf {
		activePairs := getActiveCertificatePairs(file, content, apacheRoot)
		res = append(res, activePairs...)
	}
	return res
}

func (cpi CertificatePairInfo) String() string {
	return fmt.Sprintf(`Apache File: %q
Certificate file: %q
Key file: %q`, cpi.apacheConfPath, cpi.certPath, cpi.keyPath)
}

// getEncodedCertificate opens the certificate and returns its byte sequence
func (cpi CertificatePairInfo) getEncodedCertificate() ([]byte, error) {
	return ioutil.ReadFile(cpi.certPath)
}

// getEncodedKey opens the key and returns its byte sequence
func (cpi CertificatePairInfo) getEncodedKey() ([]byte, error) {
	return ioutil.ReadFile(cpi.certPath)
}

// getDecodedCertificateInfo returns the certificate domain name or an error if it cannot be opened or decoded
func (cpi CertificatePairInfo) getCertificateDomainName(encodedCert []byte) (string, error) {
	res := ""
	block, _ := pem.Decode(encodedCert)
	parsedCert, err := x509.ParseCertificate(block.Bytes)
	if err == nil {
		res = parsedCert.Subject.CommonName
	}
	return res, err
}

// printSSLCertsInfo prints on screen the certificate domain name
func (cpi CertificatePairInfo) printCertificateDomain() error {
	encodedCert, err := cpi.getEncodedCertificate()
	if err != nil {
		return err
	}
	domain, err := cpi.getCertificateDomainName(encodedCert)
	if err == nil {
		fmt.Printf("Domain name: %q\n", domain)
	}
	return err
}

// getCertKeyMatchInfo returns, for each active certificate-key pair, whether they match or not
func (cpi CertificatePairInfo) certKeyMatch() (bool) {
	res := true
	_, err := tls.LoadX509KeyPair(cpi.certPath, cpi.keyPath)
	if err != nil {
		res = false
	}
	return res
}

// HTTPSConnectionInfo contains paths of an active certificate-key path
type HTTPSConnectionInfo struct {
	hostname string
	port     int
}

func (httpsConnInfo HTTPSConnectionInfo) String() string {
	return fmt.Sprintf(`Hostname: %q
Port: %d`, httpsConnInfo.hostname, httpsConnInfo.port)
}

// printCertKeyMatchInfo prints, for each active certificate-key pair, whether they match or not
func (cpi CertificatePairInfo) printCertKeyMatchInfo() {
	match := cpi.certKeyMatch()
	fmt.Printf("Certificate and key match: %t\n", match)
}

// getServerCertificateDomain attempts a HTTPS connection to the server and returns the returned certificate domain name
func (httpsConnInfo HTTPSConnectionInfo) getServerCertificateDomain() (string, error) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	connectionString := fmt.Sprintf("%s:%d", httpsConnInfo.hostname, httpsConnInfo.port)
	conn, err := tls.Dial("tcp", connectionString, conf)
	if err != nil {
		return "", err
	}
	conn.Handshake()
	res := conn.ConnectionState().PeerCertificates[0].Subject.CommonName
	conn.Close()
	return res, err
}

// printHTTPSConnectionInfo prints the results of the HTTPS connection attempt to the server
func (httpsConnInfo HTTPSConnectionInfo) printHTTPSConnectionInfo() error {
	fmt.Printf("%s\n", httpsConnInfo)
	domain, err := httpsConnInfo.getServerCertificateDomain()
	if err == nil {
		fmt.Printf("Server certificate domain: %q\n", domain)
	}
	return err
}

// RunActiveCertificatesChecks performs checks on the active certificate key pairs in the Apache configuration
func RunActiveCertificatesChecks(confFile, apacheRoot string) error {
	apacheConf, err := apache.OpenAllApacheConfigurationFiles(confFile, apacheRoot)
	if err != nil {
		return err
	}
	certKeyPairs := getActiveCertificatePairsInAllFiles(apacheConf, apacheRoot)
	if len(certKeyPairs) == 0 {
		fmt.Println("No SSL certificates found in the Apache configuration")
	} else {
		for index, cpi := range certKeyPairs {
			fmt.Printf("Ocurrence #%d\n%s\n", index+1, cpi)
			err = cpi.printCertificateDomain()
			if err != nil {
				return err
			}
			cpi.printCertKeyMatchInfo()
		}
	}
	return err
}

// RunHTTPSConnectionChecks performs checks on the HTTPS connection to web server
func RunHTTPSConnectionChecks(hostname string, port int) error {
	httpsConnection := HTTPSConnectionInfo{hostname, port}
	err := httpsConnection.printHTTPSConnectionInfo()
	return err
}
