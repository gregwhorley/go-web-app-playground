package main

import (
	"github.com/mvmaasakkers/certificates/cert"
	"log"
	"time"
)

var (
	serialNumber, _ = cert.GenerateRandomBigInt()
	notBefore       = time.Now()
	notAfter        = time.Date(2021, time.June, 1, 1, 0, 0, 0, time.Local)
)

func GenerateCertRequest() *cert.Request {
	return &cert.Request{
		Organization:     "Chaos.org",
		Country:          "US",
		Province:         "OR",
		Locality:         "Portland",
		StreetAddress:    "P.O. Box 69420",
		PostalCode:       "12345",
		CommonName:       "example.com",
		SerialNumber:     serialNumber,
		NameSerialNumber: "serial number",
		SubjectAltNames:  nil,
		NotBefore:        notBefore,
		NotAfter:         notAfter,
		BitSize:          4096,
	}
}

func GenerateCertAuthority(r *cert.Request) (caCert []byte, privateKey []byte) {
	caCert, privateKey, err := cert.GenerateCA(r)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func GenerateCert(r *cert.Request, caCert []byte, key []byte) (certPem []byte, certKey []byte) {
	certPem, certKey, err := cert.GenerateCertificate(r, caCert, key)
	if err != nil {
		log.Fatal(err)
	}
	return
}
