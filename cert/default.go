package cert

import (
	"crypto/x509"

	"github.com/pkg/errors"
)

func DefaultTrustPool() (*x509.CertPool, error) {
	pool, err := LoadSystemRoots()
	if err != nil {
		return nil, errors.Wrap(err, "loading system root certificates pool")
	}
	if !pool.AppendCertsFromPEM([]byte(rootAFIP)) {
		return nil, errors.Errorf("unable to add AFIP root CA certificate to trusted pool")
	}
	return pool, nil
}
