package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	remoteHostName, userName, keyFile := getArgs()
	cmdToRun := "ps aux | grep gunicorn"

	if keyFile == "" {
		// use private key
		var err error

		fmt.Print("Enter path to private key: ")
		keyFile, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal("Failed to read private key path: ", err)
		}
		keyFile = strings.TrimRight(keyFile, "\n")
	}

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

	remoteHost := fmt.Sprintf("%s:%d", remoteHostName, 22)
	log.Printf("Connecting to %s...\n", remoteHost)
	client, err := ssh.Dial("tcp", remoteHost, config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	defer client.Close()

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(cmdToRun); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())
}

func getArgs() (string, string, string) {
	if len(os.Args) > 3 {
		return os.Args[1], os.Args[2], os.Args[3]
	}
	if len(os.Args) > 2 {
		return os.Args[1], os.Args[2], ""
	}

	log.Fatal("usage incorrect")
	return "", "", ""
}
