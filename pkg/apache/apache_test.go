package apache

import (
	"testing"
)

func testEq(a, b []string) bool {

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

func TestGetIncludes(t *testing.T) {
	apacheRoot := "/opt/bitnami/apache2/"
	testData := []struct {
		in  string
		out []string
	}{
		{`
SSLRandomSeed startup builtin
#    Include "../apps/example.conf"
`, []string{}},
		{`
    Include "../apps/wordpress/conf/httpd-app.conf"
    Include "../apps/example.conf"
`, []string{"/opt/bitnami/apps/wordpress/conf/httpd-app.conf",
			"/opt/bitnami/apps/example.conf"}},
	}

	t.Run("Check Detected Apache include files", func(t *testing.T) {
		for _, tt := range testData {
			detectedIncludes := GetIncludes(tt.in, apacheRoot)
			if !testEq(tt.out, detectedIncludes) {
				t.Errorf("Detected includes incorrect for configuration: %s\n\n expected: %q, got: %q", tt.in,
					tt.out, detectedIncludes)
			}
		}
	})
}
