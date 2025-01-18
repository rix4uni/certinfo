package main

import (
	"bufio"
	"crypto/tls"
	"encoding/csv"
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
	if strings.HasPrefix(url, "*.") {
		return strings.TrimPrefix(url, "*.")
	}
	return url
}

func worker(jobs <-chan string, results chan<- CertificateDetails, wg *sync.WaitGroup, verbose bool, san bool, issued bool, expires bool, today bool, timeout time.Duration) {
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
			Timeout:   timeout,
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

		// Handle multiple flags
		if today && issued {
			// Match today's date with the issued date
			todayDate := time.Now().Format("2006-01-02")
			issuedDate := cert.NotBefore.Format("2006-01-02")
			if todayDate == issuedDate {
				fmt.Printf("%s:%s [%s]\n", certDetails.Host, port, certDetails.ValidityPeriod["Issued_On"])
			}
		} else if san && issued && expires {
			fmt.Printf("%s:%s [%s] [%s] [%s]\n", certDetails.Host, port, certDetails.ValidityPeriod["Issued_On"], certDetails.ValidityPeriod["Expires_On"], strings.Join(certDetails.CertificateSubjectAlternativeName, ", "))
		} else if san && issued {
			fmt.Printf("%s:%s [%s] [%s]\n", certDetails.Host, port, certDetails.ValidityPeriod["Issued_On"], strings.Join(certDetails.CertificateSubjectAlternativeName, ", "))
		} else if san && expires {
			fmt.Printf("%s:%s [%s] [%s]\n", certDetails.Host, port, certDetails.ValidityPeriod["Expires_On"], strings.Join(certDetails.CertificateSubjectAlternativeName, ", "))
		} else if issued && expires {
			fmt.Printf("%s:%s [%s] [%s]\n", certDetails.Host, port, certDetails.ValidityPeriod["Issued_On"], certDetails.ValidityPeriod["Expires_On"])
		} else if san {
			fmt.Printf("%s:%s [%s]\n", certDetails.Host, port, strings.Join(certDetails.CertificateSubjectAlternativeName, ", "))
		} else if issued {
			fmt.Printf("%s:%s [%s]\n", certDetails.Host, port, certDetails.ValidityPeriod["Issued_On"])
		} else if expires {
			fmt.Printf("%s:%s [%s]\n", certDetails.Host, port, certDetails.ValidityPeriod["Expires_On"])
		} else {
			// Send result to the results channel
			results <- certDetails
		}
	}
}

func main() {
	// Parse flags
	jsonOutput := flag.Bool("json", false, "output in JSON format")
	csvOutput := flag.Bool("csv", false, "output in CSV format")
	concurrency := flag.Int("c", 50, "number of concurrent workers")
	silent := flag.Bool("silent", false, "silent mode.")
	versionFlag := flag.Bool("version", false, "Print the version of the tool and exit.")
	verbose := flag.Bool("verbose", false, "enable verbose logging")
	san := flag.Bool("san", false, "monitor the san certificate details in a simple format")
	issued := flag.Bool("issued", false, "output host, port, and certificate expiration date")
	expires := flag.Bool("expires", false, "output only the expiration date")
	today := flag.Bool("today", false, "filter results to show only certificates issued today (works only with -issued flag)")
	timeoutStr := flag.String("timeout", "3s", "connection timeout duration (e.g. 5s, 10m, 1h)")

	flag.Parse()

	if *versionFlag {
		banner.PrintBanner()
		banner.PrintVersion()
		return
	}

	if !*silent {
		banner.PrintBanner()
	}

	// Parse the timeout value
	timeout, err := time.ParseDuration(*timeoutStr)
	if err != nil {
		fmt.Printf("Error parsing timeout duration: %v\n", err)
		return
	}

	// Check if -today is used without -issued or with other flags
	if *today && (!*issued || *san || *expires || *jsonOutput || *csvOutput) {
		fmt.Println("Error: -today flag can only be used with the -issued flag and no other flags.")
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	jobs := make(chan string, *concurrency)
	results := make(chan CertificateDetails)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg, *verbose, *san, *issued, *expires, *today, timeout)
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
	if *csvOutput {
		// Prepare CSV writer
		writer := csv.NewWriter(os.Stdout)
		defer writer.Flush()

		// Write CSV header
		writer.Write([]string{
			"Host", "IssuedTo_CommonName", "IssuedTo_Organization",
			"IssuedBy_CommonName", "IssuedBy_Organization",
			"IssuedOn", "ExpiresOn", "SubjectAlternativeNames",
		})

		for certDetails := range results {
			// Prepare and write CSV rows
			row := []string{
				certDetails.Host,
				certDetails.IssuedTo["Common_Name_(CN)"],
				certDetails.IssuedTo["Organization_(O)"],
				certDetails.IssuedBy["Common_Name_(CN)"],
				certDetails.IssuedBy["Organization_(O)"],
				certDetails.ValidityPeriod["Issued_On"],
				certDetails.ValidityPeriod["Expires_On"],
				strings.Join(certDetails.CertificateSubjectAlternativeName, ", "),
			}
			writer.Write(row)
		}
	} else {
		for certDetails := range results {
			if *jsonOutput {
				// JSON output
				b, err := json.MarshalIndent(certDetails, "", "  ")
				if err != nil {
					fmt.Println("Failed to convert to JSON:", err)
					continue
				}
				fmt.Println(string(b))
			} else if !*san {
				// Default (non-JSON) output
				for _, name := range certDetails.CertificateSubjectAlternativeName {
					fmt.Println(name)
				}
			}
		}
	}

	if scanner.Err() != nil {
		fmt.Println("Failed to read from stdin:", scanner.Err())
	}
}
