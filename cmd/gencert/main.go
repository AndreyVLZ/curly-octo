package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"github.com/AndreyVLZ/curly-octo/keygen"
)

const (
	caCertFile     = "cert/ca-cert.pem"
	serverCertFile = "cert/server-cert.pem"
	serverKeyFile  = "cert/server-key.pem"
	clientCertFile = "cert/client-cert.pem"
	clientKeyFile  = "cert/client-key.pem"
)

func main() {
	rootCert, rootCertDER, rootKey, err := keygen.GenerateRootCertificate()
	if err != nil {
		log.Printf("generate root certificate: %v\n", err)

		return
	}

	serverCert, serverKey, err := keygen.GenerateServerCertificate(rootKey, rootCert)
	if err != nil {
		log.Printf("generate server certificate: %v\n", err)

		return
	}

	clientCert, clientKey, err := keygen.GenerateClientCartificate(rootKey, rootCert)
	if err != nil {
		log.Printf("generate client certificate: %v\n", err)

		return
	}

	if err := writeCert(caCertFile, rootCertDER); err != nil {
		log.Printf("write ca cert: %v\n", err)

		return
	}

	if err := writeCert(serverCertFile, serverCert); err != nil {
		log.Printf("write server cert: %v\n", err)

		return
	}

	if err := writeKey(serverKeyFile, serverKey); err != nil {
		log.Printf("write server key: %v\n", err)

		return
	}

	if err := writeCert(clientCertFile, clientCert); err != nil {
		log.Printf("write client cert: %v\n", err)

		return
	}

	if err := writeKey(clientKeyFile, clientKey); err != nil {
		log.Printf("write client key: %v\n", err)

		return
	}

	log.Println("generate certificate OK")
}

func writeCert(filePath string, cert []byte) error {
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})

	return write(filePath, certPEM)
}

func writeKey(filePath string, key *rsa.PrivateKey) error {
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	return write(filePath, keyPEM)
}

func write(filePath string, data []byte) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("file open: %w", err)
	}

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("file write: %w", err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("file close: %w", err)
	}

	return nil
}
