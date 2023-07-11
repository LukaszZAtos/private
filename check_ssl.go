package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const (
	port             = 443                         // Replace with the appropriate port
	notificationDays = 14                          // Number of days before expiration to receive notifications
	exitCodeOK       = 0                            // Exit code for successful execution
	exitCodeWarning  = 1                            // Exit code for certificate expiration warning
	exitCodeError    = 2                            // Exit code for connection or other errors
)

func main() {
	// Check if hostname argument is provided
	if len(os.Args) < 2 {
		log.Fatal("Please provide a hostname argument.")
	}

	hostname := os.Args[1]

	// Create a deadline for the connection to complete
	dialer := net.Dialer{Timeout: 5 * time.Second}

	// Dial the TCP connection
	conn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port))
	if err != nil {
		log.Printf("Failed to connect to %s:%d - %v", hostname, port, err)
		os.Exit(exitCodeError)
	}
	defer conn.Close()

	// Create the TLS configuration
	config := tls.Config{
		InsecureSkipVerify: true,
		ServerName:         hostname,
	}

	// Create the TLS connection
	tlsConn := tls.Client(conn, &config)

	// Handshake with the server
	if err := tlsConn.Handshake(); err != nil {
		log.Printf("Failed TLS handshake with %s:%d - %v", hostname, port, err)
		os.Exit(exitCodeError)
	}

	// Retrieve the peer certificates
	certificates := tlsConn.ConnectionState().PeerCertificates

	// Check each certificate for expiration
	for _, cert := range certificates {
		expirationDate := cert.NotAfter
		daysUntilExpiration := int(expirationDate.Sub(time.Now()).Hours() / 24)

		if daysUntilExpiration <= notificationDays {
			log.Printf("Certificate for %s expires in %d days on %s", cert.Subject.CommonName, daysUntilExpiration, expirationDate)
			os.Exit(exitCodeWarning)
		}
	}

	log.Println("All certificates are valid and not expiring soon.")
	os.Exit(exitCodeOK)
}
