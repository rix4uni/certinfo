package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type CertificateDetails struct {
	Host                            string            `json:"host"`
	IssuedTo                        map[string]string `json:"Issued_To"`
	IssuedBy                        map[string]string `json:"Issued_By"`
	ValidityPeriod                  map[string]string `json:"Validity_Period"`
	CertificateSubjectAlternativeName []string         `json:"Certificate_Subject_Alternative_Name"`
}

// cleanURL removes the http:// or https:// prefix from a URL.
func cleanURL(url string) string {
	if strings.HasPrefix(url, "http://") {
		return strings.TrimPrefix(url, "http://")
	}
	if strings.HasPrefix(url, "https://") {
		return strings.TrimPrefix(url, "https://")
	}
	return url
}

func worker(jobs <-chan string, results chan<- CertificateDetails, wg *sync.WaitGroup) {
	defer wg.Done()

	for hostWithPort := range jobs {
		// Clean up the URL
		hostWithPort = cleanURL(hostWithPort)

		// Default to port 443 if not specified
		host, port, err := net.SplitHostPort(hostWithPort)
		if err != nil {
			host = hostWithPort
			port = "443"
		}

		// Set up a custom dialer with timeout
		dialer := &net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 30 * time.Second,
		}

		// Connect to the host with the custom dialer
		conn, err := tls.DialWithDialer(dialer, "tcp", net.JoinHostPort(host, port), &tls.Config{
			InsecureSkipVerify: true,
		})
		if err != nil {
			// fmt.Printf("Failed to connect to %s: %s\n", hostWithPort, err)
			continue
		}
		defer conn.Close()

		// Fetch the certificate
		certs := conn.ConnectionState().PeerCertificates
		if len(certs) == 0 {
			fmt.Printf("No certificates found for %s\n", hostWithPort)
			continue
		}
		cert := certs[0]

		// Extract the certificate details
		certDetails := CertificateDetails{
			Host: hostWithPort,
			IssuedTo: map[string]string{
				"Common_Name_(CN)": cert.Subject.CommonName,
				"Organization_(O)": strings.Join(cert.Subject.Organization, ","),
			},
			IssuedBy: map[string]string{
				"Common_Name_(CN)": cert.Issuer.CommonName,
				"Organization_(O)": strings.Join(cert.Issuer.Organization, ","),
			},
			ValidityPeriod: map[string]string{
				"Issued_On":  cert.NotBefore.Format(time.RFC3339),
				"Expires_On": cert.NotAfter.Format(time.RFC3339),
			},
		}

		// Extract the SAN (Subject Alternative Name)
		for _, name := range cert.DNSNames {
			certDetails.CertificateSubjectAlternativeName = append(certDetails.CertificateSubjectAlternativeName, name)
		}

		// Send result to the results channel
		results <- certDetails
	}
}

func main() {
	workers := 32
	scanner := bufio.NewScanner(os.Stdin)
	jobs := make(chan string, workers)
	results := make(chan CertificateDetails)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	// Read input from stdin and send jobs to workers
	go func() {
		for scanner.Scan() {
			host := scanner.Text()
			jobs <- host
		}
		close(jobs)
	}()

	// Collect results from workers
	go func() {
		wg.Wait()
		close(results)
	}()

	// Process results
	for certDetails := range results {
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
