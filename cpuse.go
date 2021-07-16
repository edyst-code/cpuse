package main

import (
    "bufio"
    "fmt"
    "os"
    "log"
    "strings"
    "syscall"

    "golang.org/x/crypto/ssh"
    "golang.org/x/term"
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

    fmt.Print("Enter password: ")
    bytePassword, err := term.ReadPassword(int(syscall.Stdin))
    if err != nil {
        log.Fatal("Failed to read password: ", err)
    }
    password := string(bytePassword)

    // a config shared between client and server
    config := &ssh.ClientConfig{
        User: userName,
        Auth: []ssh.AuthMethod{
            ssh.Password(password),
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
