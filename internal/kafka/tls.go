package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"

	"github.com/philipparndt/go-logger"
)

func NewTLSConfig(certFile, keyFile, caFile string, insecure bool) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	if insecure {
		logger.Warn("Using insecure TLS for Kafka connection")
	}

	config := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: insecure,
	}

	if caFile != "" {
		caCertBytes, err := os.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}
		caCertPool, err := x509.SystemCertPool()
		if err != nil {
			logger.Warn("Could not load system CA certificates, creating new cert pool")
			caCertPool = x509.NewCertPool()
		}
		if ok := caCertPool.AppendCertsFromPEM(caCertBytes); !ok {
			return nil, errors.New("failed to append CA certificate to pool")
		}

		config.RootCAs = caCertPool
	}

	return config, nil
}
