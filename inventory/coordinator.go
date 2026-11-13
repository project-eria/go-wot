package inventory

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"time"
)

// RoleCoordinator is the name the coordinator's certificate carries (AUTH-DESIGN
// §1.1): signed by the deployment key, whatever the address it is reached at.
const RoleCoordinator = "coordinateur"

// CoordinatorClient reaches the coordinator: over https, only a certificate of the
// deployment roots (PEM) for RoleCoordinator; over http (before TLS), as is.
func CoordinatorClient(rootsPEM []byte, timeout time.Duration) (*http.Client, error) {
	pool := x509.NewCertPool()
	if len(rootsPEM) > 0 && !pool.AppendCertsFromPEM(rootsPEM) {
		return nil, errors.New("no certificate in the deployment roots")
	}
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.TLSClientConfig = &tls.Config{RootCAs: pool, ServerName: RoleCoordinator, MinVersion: tls.VersionTLS13}
	return &http.Client{Timeout: timeout, Transport: t}, nil
}
