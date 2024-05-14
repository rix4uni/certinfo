package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type CertificateDetails struct {
	Host                            string            `json:"host"`
	IssuedTo                        map[string]string `json:"Issued_To"`
	IssuedBy                        map[string]string `json:"Issued_By"`
	ValidityPeriod                  map[string]string `json:"Validity_Period"`
	CertificateSubjectAlternativeName []string         `json:"Certificate_Subject_Alternative_Name"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	port := "443"

	for scanner.Scan() {
		host := scanner.Text()
		
		// Connect to the host
		conn, err := tls.Dial("tcp", net.JoinHostPort(host, port), &tls.Config{
			InsecureSkipVerify: true,
		})
		if err != nil {
			// fmt.Println("Failed to connect:", err)
			continue
		}
		defer conn.Close()

		// Fetch the certificate
		certs := conn.ConnectionState().PeerCertificates
		if len(certs) == 0 {
			fmt.Println("No certificates found.")
			continue
		}
		cert := certs[0]

		// Extract the certificate details
		certDetails := CertificateDetails{
			Host: host,
			IssuedTo: map[string]string{
				"Common_Name_(CN)": cert.Subject.CommonName,
				"Organization_(O)": strings.Join(cert.Subject.Organization, ","),
			},
			IssuedBy: map[string]string{
				"Common_Name_(CN)": cert.Issuer.CommonName,
				"Organization_(O)": strings.Join(cert.Issuer.Organization, ","),
			},
			ValidityPeriod: map[string]string{
				"Issued_On": cert.NotBefore.Format(time.RFC3339),
				"Expires_On": cert.NotAfter.Format(time.RFC3339),
			},
		}

		// Extract the SAN (Subject Alternative Name)
		for _, name := range cert.DNSNames {
			certDetails.CertificateSubjectAlternativeName = append(certDetails.CertificateSubjectAlternativeName, name)
		}

		// Convert to JSON
		b, err := json.MarshalIndent(certDetails, "", "  ")
		if err != nil {
			fmt.Println("Failed to convert to JSON:", err)
			continue
		}
		fmt.Println(string(b))
	}

	if scanner.Err() != nil {
		fmt.Println("Failed to read from stdin:", scanner.Err())
	}
}
