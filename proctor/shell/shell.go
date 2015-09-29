package shell

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

type Runner struct{}

type ConnectionOptions struct {
	Username      string
	Port          int
	PrivateKeyPEM []byte
}

func (r *Runner) ConnectAndRun(host, command string, options *ConnectionOptions) (string, error) {
	signer, err := ssh.ParsePrivateKey(options.PrivateKeyPEM)
	if err != nil {
		return "", err
	}

	config := &ssh.ClientConfig{
		User: options.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, options.Port), config)
	if err != nil {
		return "", fmt.Errorf("failed to dial: %s", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: ", err)
	}
	defer session.Close()

	var stdoutBytes bytes.Buffer
	session.Stdout = &stdoutBytes
	session.Stderr = os.Stderr
	if err := session.Run(command); err != nil {
		return "", fmt.Errorf("failed while running command: %s", err)
	}
	return stdoutBytes.String(), nil
}

func copy(size int64, mode os.FileMode, fileName string, contents io.Reader, destination string, session *ssh.Session) error {
	defer session.Close()
	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()
		fmt.Fprintf(w, "C%#o %d %s\n", mode, size, fileName)
		io.Copy(w, contents)
		fmt.Fprint(w, "\x00")
	}()
	cmd := fmt.Sprintf("scp -t %s", destination)
	if err := session.Run(cmd); err != nil {
		return err
	}
	return nil
}

func scpAndRun(client ssh.Client) {
	scpSession, err := client.NewSession()
	if err != nil {
		panic("Failed to create SCP session: " + err.Error())
	}
	defer scpSession.Close()

	scriptContents := `#!/bin/bash

echo "this script is located at $dirname $0"
`
	scriptReader := strings.NewReader(scriptContents)
	scpError := copy(int64(len(scriptContents)), os.FileMode(0777), "test-script", scriptReader, "/tmp/scripts/", scpSession)
	if scpError != nil {
		panic(scpError)
	}

	execSession, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer execSession.Close()

	var stdoutBytes bytes.Buffer
	execSession.Stdout = &stdoutBytes
	if err := execSession.Run("/tmp/scripts/test-script"); err != nil {
		panic("Failed to run: " + err.Error())
	}
}
