package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rix4uni/certinfo/banner"
)

type CertificateDetails struct {
	Host                              string            `json:"host"`
	IssuedTo                          map[string]string `json:"Issued_To"`
	IssuedBy                          map[string]string `json:"Issued_By"`
	ValidityPeriod                    map[string]string `json:"Validity_Period"`
	CertificateSubjectAlternativeName []string          `json:"Certificate_Subject_Alternative_Name"`
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

func worker(jobs <-chan string, results chan<- CertificateDetails, wg *sync.WaitGroup, verbose bool, moniter bool) {
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
			if verbose {
				fmt.Printf("Failed to connect to %s: %s\n", hostWithPort, err)
			}
			continue
		}
		defer conn.Close()

		// Fetch the certificate
		certs := conn.ConnectionState().PeerCertificates
		if len(certs) == 0 {
			if verbose {
				fmt.Printf("No certificates found for %s\n", hostWithPort)
			}
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

		// Print the certificate details
		if moniter {
			// Print host, port, and SAN if -moniter flag used
			fmt.Printf("%s:%s [%s]\n", certDetails.Host, port, strings.Join(certDetails.CertificateSubjectAlternativeName, ", "))
		}
	}
}

func main() {
	// Parse flags
	jsonOutput := flag.Bool("json", false, "output in JSON format")
	concurrency := flag.Int("c", 50, "number of concurrent workers")
	silent := flag.Bool("silent", false, "silent mode.")
	versionFlag := flag.Bool("version", false, "Print the version of the tool and exit.")
	verbose := flag.Bool("verbose", false, "enable verbose logging")
	moniter := flag.Bool("moniter", false, "monitor the certificate details in a simple format")
	flag.Parse()

	if *versionFlag {
		banner.PrintBanner()
		banner.PrintVersion()
		return
	}

	if !*silent {
		banner.PrintBanner()
	}

	scanner := bufio.NewScanner(os.Stdin)
	jobs := make(chan string, *concurrency)
	results := make(chan CertificateDetails)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg, *verbose, *moniter)
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
		if *jsonOutput {
			// JSON output
			b, err := json.MarshalIndent(certDetails, "", "  ")
			if err != nil {
				fmt.Println("Failed to convert to JSON:", err)
				continue
			}
			fmt.Println(string(b))
		} else if !*moniter {
			// Default (non-JSON) output
			for _, name := range certDetails.CertificateSubjectAlternativeName {
				fmt.Println(name)
			}
		}
	}

	if scanner.Err() != nil {
		fmt.Println("Failed to read from stdin:", scanner.Err())
	}
}
