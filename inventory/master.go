package inventory

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"time"
)

// RoleMaster is the name the master's certificate carries (AUTH-DESIGN
// §1.1): signed by the deployment key, whatever the address it is reached at.
const RoleMaster = "maitre"

// MasterClient reaches the master: over https, only a certificate of the
// deployment roots (PEM) for RoleMaster; over http (before TLS), as is.
func MasterClient(rootsPEM []byte, timeout time.Duration) (*http.Client, error) {
	pool := x509.NewCertPool()
	if len(rootsPEM) > 0 && !pool.AppendCertsFromPEM(rootsPEM) {
		return nil, errors.New("no certificate in the deployment roots")
	}
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.TLSClientConfig = &tls.Config{RootCAs: pool, ServerName: RoleMaster, MinVersion: tls.VersionTLS13}
	return &http.Client{Timeout: timeout, Transport: t}, nil
}
