package security

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/pkg/errors"
	"io/ioutil"
)

type SecurityConfig struct {
	SSLCAPath   string `json:"ssl-ca-path"`
	SSLCertPath string `json:"ssl-cert-path"`
	SSLKeyPath  string `json:"ssl-key-path"`
}

func (s *SecurityConfig) ToTLSConfig() (*tls.Config, error) {
	var tlsConfig *tls.Config
	if len(s.SSLCAPath) != 0 {
		var certificates = make([]tls.Certificate, 0)
		if len(s.SSLCertPath) != 0 && len(s.SSLKeyPath) != 0 {
			// Load the client certificates from disk
			certificate, err := tls.LoadX509KeyPair(s.SSLCertPath, s.SSLKeyPath)
			if err != nil {
				return nil, errors.Errorf("could not load client key pair: %s", err)
			}
			certificates = append(certificates, certificate)
		}

		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(s.SSLCAPath)
		if err != nil {
			return nil, errors.Errorf("could not read ca certificate: %s", err)
		}

		// Append the certificates from the CA
		if !certPool.AppendCertsFromPEM(ca) {
			return nil, errors.New("failed to append ca certs")
		}

		tlsConfig = &tls.Config{
			Certificates: certificates,
			RootCAs:      certPool,
		}
	}

	return tlsConfig, nil
}
