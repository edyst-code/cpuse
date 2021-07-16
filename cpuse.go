package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Get host
	fmt.Print("Enter remote host: ")
	remoteHostName, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Failed to read host: ", err)
	}
	remoteHostName = strings.TrimRight(remoteHostName, "\n")

	remotePort := 22
	remoteHost := fmt.Sprintf("%s:%d", remoteHostName, remotePort)

	fmt.Print("Enter username: ")
	userName, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Failed to read username: ", err)
	}
	userName = strings.TrimRight(userName, "\n")

	// use private key
	fmt.Print("Enter path to private key: ")
	keyFile, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Failed to read private key path: ", err)
	}
	keyFile = strings.TrimRight(keyFile, "\n")

	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Fatalf("Unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("Unable to parse private key: %v", err)
	}

	// a config shared between client and server
	config := &ssh.ClientConfig{
		User: userName,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	log.Printf("Connecting to %s...\n", remoteHost)
	client, err := ssh.Dial("tcp", remoteHost, config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	defer client.Close()
}
