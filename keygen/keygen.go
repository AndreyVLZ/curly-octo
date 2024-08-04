package keygen

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"net"
	"time"
)

func GenerateClientCartificate(rootKey *rsa.PrivateKey, rootCert *x509.Certificate) ([]byte, *rsa.PrivateKey, error) {
	clientKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("gen rsa key: %w", err)
	}

	clientCertTmpl := certTemplate()
	clientCertTmpl.KeyUsage = x509.KeyUsageDigitalSignature
	clientCertTmpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}

	clientCertDER, err := x509.CreateCertificate(rand.Reader, clientCertTmpl, rootCert, &clientKey.PublicKey, rootKey)
	if err != nil {
		return nil, nil, fmt.Errorf("create x509 cert: %w", err)
	}

	return clientCertDER, clientKey, nil
}

func GenerateServerCertificate(rootKey *rsa.PrivateKey, rootCert *x509.Certificate) ([]byte, *rsa.PrivateKey, error) {
	servKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("gen rand: %w", err)
	}

	servCertTmpl := certTemplate()
	servCertTmpl.KeyUsage = x509.KeyUsageDigitalSignature
	servCertTmpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	servCertTmpl.IPAddresses = []net.IP{net.ParseIP("127.0.0.1")}

	serverCertDER, err := x509.CreateCertificate(rand.Reader, servCertTmpl, rootCert, &servKey.PublicKey, rootKey)
	if err != nil {
		return nil, nil, fmt.Errorf("gen cert: %w", err)
	}

	return serverCertDER, servKey, nil
}

func GenerateRootCertificate() (*x509.Certificate, []byte, *rsa.PrivateKey, error) {
	rootKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("gen key: %w", err)
	}

	rootCertTmpl := certTemplate()
	rootCertTmpl.IsCA = true
	rootCertTmpl.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature
	rootCertTmpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
	rootCertTmpl.IPAddresses = []net.IP{net.ParseIP("127.0.0.1")}

	rootCertDER, err := x509.CreateCertificate(rand.Reader, rootCertTmpl, rootCertTmpl, &rootKey.PublicKey, rootKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("gen cert: %w", err)
	}

	rootCert, err := x509.ParseCertificate(rootCertDER)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("gen parse cert: %w", err)
	}

	return rootCert, rootCertDER, rootKey, nil
}

func certTemplate() *x509.Certificate {
	tmpl := x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"ООО 'Ромашка'"},
		},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), // +10 лет
		BasicConstraintsValid: true,
	}

	return &tmpl
}
