package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh"
)

var servers = []string{
	"server1.example.com",
	"server2.example.com",
	"server3.example.com",
	"server4.example.com",
	"server5.example.com",
}

func main() {
	for _, server := range servers {
		exitCode := checkSSHConnectivity(server)
		if exitCode == 0 {
			log.Printf("%s: SSH connectivity successful\n", server)
		} else {
			log.Printf("%s: SSH connectivity failed (exit code: %d)\n", server, exitCode)
		}
	}
}

func checkSSHConnectivity(server string) int {
	// Read the private key file
	keyPath := "/path/to/private_key"
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("Failed to read private key file: %v", err)
	}

	// Parse the private key
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// SSH client configuration
	config := &ssh.ClientConfig{
		User: "<username>",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Establish SSH connection
	client, err := ssh.Dial("tcp", server+":22", config)
	if err != nil {
		if strings.Contains(err.Error(), "unable to authenticate") {
			return 1 // Authentication failed
		} else {
			return 2 // Connection failed
		}
	}
	defer client.Close()

	// Run a command to check SSH connectivity
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create SSH session: %v", err)
	}
	defer session.Close()

	// Customize the command to be executed for connectivity check
	cmd := "echo 'SSH connectivity check'"
	err = session.Run(cmd)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode() // Command execution failed
		} else {
			return 3 // Other error occurred
		}
	}

	return 0 // SSH connectivity successful
}
