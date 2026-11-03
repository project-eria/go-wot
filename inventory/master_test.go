package inventory

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The master's certificate: of the deployment roots, for RoleMaster only.
func TestMasterClient(t *testing.T) {
	rootPub, rootPriv, _ := ed25519.GenerateKey(rand.Reader)
	root := &x509.Certificate{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "root"}, IsCA: true, BasicConstraintsValid: true,
		KeyUsage: x509.KeyUsageCertSign, NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour)}
	rootDER, err := x509.CreateCertificate(rand.Reader, root, root, rootPub, rootPriv)
	require.NoError(t, err)
	root, _ = x509.ParseCertificate(rootDER)
	serve := func(name string) *httptest.Server {
		pub, priv, _ := ed25519.GenerateKey(rand.Reader)
		leaf := &x509.Certificate{SerialNumber: big.NewInt(2), Subject: pkix.Name{CommonName: name}, DNSNames: []string{name},
			NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour), ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}}
		der, err := x509.CreateCertificate(rand.Reader, leaf, root, pub, rootPriv)
		require.NoError(t, err)
		srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("[]")) }))
		srv.TLS = &tls.Config{Certificates: []tls.Certificate{{Certificate: [][]byte{der}, PrivateKey: priv}}}
		srv.StartTLS()
		t.Cleanup(srv.Close)
		return srv
	}
	roots := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: rootDER})
	c, err := MasterClient(roots, 5*time.Second)
	require.NoError(t, err)

	resp, err := c.Get(serve(RoleMaster).URL)
	require.NoError(t, err, "the master's certificate, reached at an IP")
	resp.Body.Close()
	_, err = c.Get(serve("relais-interieur").URL)
	assert.Error(t, err, "another role of the house")
	_, err = MasterClient([]byte("nothing"), time.Second)
	assert.Error(t, err)
}
