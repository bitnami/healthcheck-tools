package main

import (
	"log"
	"os"
	"testing"
)

func testEq(a, b []CertificatePairInfo) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestGetActiveCertificatePairs(t *testing.T) {
	apacheRoot := "/opt/bitnami/apache2/"
	apacheConf := "/opt/bitnami/apache2/conf/httpd.conf"
	testData := []struct {
		in  string
		out []CertificatePairInfo
	}{
		{`
SSLRandomSeed startup builtin
SSLRandomSeed connect builtin
`, []CertificatePairInfo{}},
		{`
    SSLCertificateFile "../apps/wordpress/conf/certs/server.crt"
    SSLCertificateKeyFile "../apps/wordpress/conf/certs/server.key"
`, []CertificatePairInfo{{"/opt/bitnami/apache2/conf/httpd.conf",
			"/opt/bitnami/apps/wordpress/conf/certs/server.crt",
			"/opt/bitnami/apps/wordpress/conf/certs/server.key"}}},
		{`
    SSLCertificateFile "../apps/wordpress/conf/certs/server.crt"
   # SSLCertificateKeyFile "../apps/wordpress/conf/certs/server.key"
`, []CertificatePairInfo{}},
		{`
    SSLCertificateKeyFile "../apps/wordpress/conf/certs/server.key"
    SSLCertificateFile "../apps/wordpress/conf/certs/server.crt"
   # SSLCertificateKeyFile "../apps/wordpress/conf/certs/server3.key"
    SSLCertificateFile "../apps/wordpress/conf/certs/server2.crt"
    SSLCertificateKeyFile "../apps/wordpress/conf/certs/server2.key"
`, []CertificatePairInfo{{"/opt/bitnami/apache2/conf/httpd.conf",
			"/opt/bitnami/apps/wordpress/conf/certs/server.crt",
			"/opt/bitnami/apps/wordpress/conf/certs/server.key"},
			{"/opt/bitnami/apache2/conf/httpd.conf",
				"/opt/bitnami/apps/wordpress/conf/certs/server2.crt",
				"/opt/bitnami/apps/wordpress/conf/certs/server2.key"}}},
	}

	t.Run("Check Detected SSL files", func(t *testing.T) {
		for _, tt := range testData {
			detectedCerts := getActiveCertificatePairs(apacheConf, tt.in, apacheRoot)
			if !testEq(tt.out, detectedCerts) {
				t.Errorf("Detected certs incorrect for configuration: %s\n\n expected: %q, got: %q", tt.in,
					tt.out, detectedCerts)
			}
		}
	})
}

var testCertificate = `-----BEGIN CERTIFICATE-----
MIICqDCCAZACCQCz8T3726LYsjANBgkqhkiG9w0BAQUFADAWMRQwEgYDVQQDDAtl
eGFtcGxlLmNvbTAeFw0xMjExMTQxMTE4MjdaFw0yMjExMTIxMTE4MjdaMBYxFDAS
BgNVBAMMC2V4YW1wbGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEA5NHl5TfZtO6zugau2tp5mWIcQYJhuwKTmYeXDLYAGJpoD2SixwPL5c8glneI
Rz1N2EQIZVeaWGbS0FLFlPdOkCkplpW9isYVC4XqKrk5b4HW4+YC+Cup0k+Kd4NM
eZOTUvWr5N6dIpdibkVumBc/pao8VtdwywlCL/PwGRsQtkXrRICzdtRa3MXqTmEF
foyVCGgBRtronlB9x4Plfb8Psk4GrPkjrWYgO8peKrl0O5+F+sYg7Gj95zCH73BQ
ANzCVNrgD9fs9cyx3ru9CUdEoIxAAJwQFkjm7xr6xqhIlSgnQ7B0uOSTNRcXY6rw
s+PxGneec/kRPRgzjC/QHY6n8QIDAQABMA0GCSqGSIb3DQEBBQUAA4IBAQBbyMqF
RDsX8zX1EW5qA8AQ8Jb2XqWrVeSO8blMV3WagJ2airMm3+c/82FCwsd/cZ08UXhA
/Kou0gi/F16tV26PiiUdp590Qao3d8H2qxc1rzzULimZPgxH4iA4vRyMHtyZN6h4
7Fdn7O9xNMPu8siOz8rrzsEdEX5URbOMkDLCZsbTIUWVv2XmqrR0K10d5VuLWeLi
r+4G6c6jpa244WmqT9ClqceJ12G1Wnmezy7ybiW0l5M2iuIKFEiRP5Hj0J15o1I2
pXAbKysAdWRHsJSQOtcgO8Vh9k0wo3tKg4HDp1hbrEzoGzOv92Vjg3lG8X+hzbMJ
MQURotHkD4Gk57wL
-----END CERTIFICATE-----`

var testKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA5NHl5TfZtO6zugau2tp5mWIcQYJhuwKTmYeXDLYAGJpoD2Si
xwPL5c8glneIRz1N2EQIZVeaWGbS0FLFlPdOkCkplpW9isYVC4XqKrk5b4HW4+YC
+Cup0k+Kd4NMeZOTUvWr5N6dIpdibkVumBc/pao8VtdwywlCL/PwGRsQtkXrRICz
dtRa3MXqTmEFfoyVCGgBRtronlB9x4Plfb8Psk4GrPkjrWYgO8peKrl0O5+F+sYg
7Gj95zCH73BQANzCVNrgD9fs9cyx3ru9CUdEoIxAAJwQFkjm7xr6xqhIlSgnQ7B0
uOSTNRcXY6rws+PxGneec/kRPRgzjC/QHY6n8QIDAQABAoIBACo3G131tuGtpFTu
xLW11vdYZXQklNlGuWp63IBI162yVv54B5wF9Ek6tH1uIiNaiREcRBxGVEB4/+3V
R4SbN9Ba98RDbgu7TcipdTFaqOEMqFO1bNjSXWtip14zSBmqA2Ur1AHOnFj0awGD
J8tBhsmOpcEz0Ch1VdO5ApPvLV8jH9wQiMI/Q6yYQMtmzTMCUMYdMqe+LOziIOzL
oqN/WXnKL5E5TiO1bIxSpWPbT+IVn1c3/PShmvmRrLWsFUQlkwXJKMYZPO+rCCfe
b+Q9lMLMnj+vOnM3z16WC3aiiJGCZjVTvQ+x22YrBTRPxZmHO2eZ4H/cUQM7Y/tw
I7RjEM0CgYEA9Kxt1t8bWonzBii3P0rwyx0IECvg63k+pp4BpxpeWQKL7NVdSzk3
AyJVcNjUoZgi2kVPdxzZGLrnZfuZ691xQB3oZF0LwBzQ4GFHkTRCB0s8ZA5lcJaI
9pBu91bhz2VOZSTeQWpdMMURjXVyTXZInU1mwzmjVOIAYmO33shH9gcCgYEA72mX
UoIrFPLkOTSZOb7UbjYH01vf6ThQiYCEWg7mD3CbY7n9oobIcQMzNnt7xN4wOl/V
eKfZ7G56q8enfqm45Dyo9aCBCENVzmwO8wLe5UnvJBNL20KjvtwG8w5A6UZQzC7p
3QS+U2zxVQNEeaE6a8Wrq2d1PlhVAHYw8odgNEcCgYBN38+58xrmrz99d1oTuAt5
6kyVsRGOgPGS4HmQMRFUbT4R7DscZSKASd4945WRtTVqmWLYe4MRnvNlfzYXX0zb
ZmmAAClsRP+qWuwHaEWXwrd+9SIOOqtvJrta1/lZJFpWUOy4j10H18Flb7sosnwc
LPWHL4Iv0xriNfDg5Iga4wKBgQDLJBU59SkJBW+Q+oho7vrg6QeK15IOGbJ8eYfT
woCC6VFwNQh5N1QsUELMH8rNKJpTba18SzAl5ThBOY9tciVnw/C5Og9CK6BLHnUw
zWbDtxAq1BSxXsIB2EAtTBLX3MoB9myJFNVJhE7hi3w2mA8yEu+u6IIa/Ghjk+XE
ZAnFUQKBgQDjMinRZrK5wA09jcetI+dNiLnKHoQG6OaXDDsNCatex0O2F36BvVXE
P78qDz/i5aBMWsLx6VDvWJAkBIpZoNS5UsOn17tFaocGUSkcm48bs8Dn6VvsE8Bd
XMPAHyKuILlKYifBvNq5T22KhqKX7yGmk/AeOOiKr2KeMnh27JYrCA==
-----END RSA PRIVATE KEY-----`

var testKeyNotMatched = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA1DPBqVmPCghOo5SEX0/3AROtpwBcZcd9edP3JNlU4g1NlXJ5
hqvRQyQ0gj4w/yxNAvw7FEgfazacht8xqu0CPdDip5M4p++oRI/EYGyLNHx45ycH
2ARrdgj9230gSvWK1c53+UrbZQVtrwBhmq8HX2tRmRpLPazQra41pfabC/TetvWg
yLU7RtyVyQ+We2XGBngUhPIcGmKttM8pR1h1mk3CosZKLahWL+3f4bthOjsKi1mg
oTeZBYEYIspW2T727S9+ScBZqlNP5Ylb+oMAYxfpwsXrqcuN0slLbnmWIWyidUqR
M06peGF7QcCExmkpibtXMZPZRqzBLiDO6S3fQQIDAQABAoIBAQCGnMaPjZalwJtS
PdazN+DzN+6C8EGs9x3r+MDhCVVFiTbcRSA+hCRduUvuC1G9rfa0nBlGRnJX6u7k
yBWT3oH7gNUnhMb+EZPm2mSy3mh4RTXVPPdG25bn69BOOlQP+S+3biDBbfE7ZnQr
/cvs+nds2oGK5Bky6LwAUXEK9SNqMyaORR4Uo4+3bbFMnD4m8R0aHIEj23y5Av55
MvKBF1/LxZnWwwSPnxWvpAHZV1fMOQqlvn5OihfkioT62NSDF8cJ9hXR4XFnwGzi
lYslMNA1uu77q4hr7qgAVgvA5FFKQgyR2cOcHuU3T95/hAX8j0ijdeJ9TahnSvf6
crXC7L4hAoGBAOy8xu6GZGcWeEpqeKKInNl7+F2TRFnRO3dxLd2WPKGoPVFFUyAs
k2EGaYhKm4d769MJ/SUXtj3eg8AnkO5nG6+w8siTC+dC4+CNWLnLRvN4mGuy3wqX
9cgGCZ6z0mLmGd9qUqdVAXK2Ds0DG+BlBnKjz+lam/JYAHYZOseG+2Z1AoGBAOV3
6YfrlKXEJ8v3BSJO+32e/n5tL5i0t6/dejurVWJzwuUtPf2KrxOjZ32Glv60XaH9
3OHrwUhXwYsIVDj4cpnS69xSmPzrkXI1L98kI39MSZiMflVRY1zz6VAAMpv1UHp7
LybBoNboORpNwDtNSpfx6V9pPMJcsLTpT2Xk5bQdAoGAWPsnppXa/Shj9EyKpUTF
97TyCHIG+d98A4wF2kjS50wuJ/LvqKM4jfp0BidMyCLa48rYXG3KMP0G9l2oywL0
VBglWSB0E+t/bXKTS9pNA8xrNefYj8nINvOWABHE3SlxxhjXkk+QKManT4WAn9o5
DzPlPOeyv+c13S3kewSQT9UCgYAMoDpzRh0Zud9Os7rOlnX0BYmPP0a9KkRpItCU
8+pwzlnM7l5Y0warF/SbzYoFXbtBLIy4yZYK+vklQ0IPGGyF1jswFkNgtz17gT2v
E3f4iyQJhsF0xFOpS6pswnYGassQ0jJX+ZN1/7UUo26OVMRj8+WZYFr9fsgiTCwY
OS4CaQKBgQC/0qN9Vj6+xXFzI5BvXmdhq5Tg2wDPsRkYDOEuqwFCvHxyyhEXDehQ
sEb8efM1aFw8ZGH+mJfShWsfXEbe1nAqIuktyARR4g8MOOCFrcWWBtUDSe9nXw3Q
Up/yZWljQkrfRz/hA1/x3xN8+UCJRBxNXME1rGbb79xb4/5pRZ0CEA==
-----END RSA PRIVATE KEY-----`

func TestGetCertificateDomainName(t *testing.T) {
	t.Run("Check Detected domain", func(t *testing.T) {
		cpi := CertificatePairInfo{"/opt/bitnami/apache2/conf/httpd.conf",
			"/opt/bitnami/apps/wordpress/conf/certs/server.crt",
			"/opt/bitnami/apps/wordpress/conf/certs/server.key"}
		checkResult, err := cpi.getCertificateDomainName([]byte(testCertificate))
		if err != nil {
			t.Errorf("Error obtaining certificate domain: %s", err)
		}
		if checkResult != "example.com" {
			t.Errorf("Incorrect domain detected, expected: example.com, got: %s", checkResult)
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

func TestCertKeyMatch(t *testing.T) {
	t.Run("Check Cert and Key Match", func(t *testing.T) {

		tmpCert := createTemporaryFile(testCertificate, "cert")
		tmpKey := createTemporaryFile(testKey, "key")
		tmpKeyNotMatched := createTemporaryFile(testKeyNotMatched, "key")

		defer os.Remove(tmpCert.Name())
		defer os.Remove(tmpKey.Name())
		defer os.Remove(tmpKeyNotMatched.Name())

		correctKeyPair := CertificatePairInfo{"/opt/bitnami/apache2/httpd.conf", tmpCert.Name(), tmpKey.Name()}
		incorrectKeyPair := CertificatePairInfo{"/opt/bitnami/apache2/httpd.conf", tmpCert.Name(),
			tmpKeyNotMatched.Name()}

		checkResult := correctKeyPair.certKeyMatch()

		if !checkResult {
			t.Errorf("Incorrect certificate match detected, expected: true, got: %t", checkResult)
		}

		checkResult = incorrectKeyPair.certKeyMatch()
		if checkResult {
			t.Errorf("Incorrect certificate match detected, expected: false, got: %t", checkResult)
		}

		if err := tmpCert.Close(); err != nil {
			log.Fatal(err)
		}
		if err := tmpKey.Close(); err != nil {
			log.Fatal(err)
		}
		if err := tmpKeyNotMatched.Close(); err != nil {
			log.Fatal(err)
		}
	})
}

func TestGetServerCertificateDomain(t *testing.T) {
	httpsConnection := HTTPSConnectionInfo{"bitnami.com", 443}
	t.Run("Check HTTPS Connection", func(t *testing.T) {
		checkResult, err := httpsConnection.getServerCertificateDomain()
		if err != nil {
			t.Errorf("Error creating HTTPS request: %s", err)
		}
		if checkResult != "bitnami.com" {
			t.Errorf("Incorrect HTTPS Server certificate detected, expected: bitnami.com, got: %s",
				checkResult)
		}
	})
}
