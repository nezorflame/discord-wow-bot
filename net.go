package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/pkg/sftp"

	"sync"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// All the connections
var (
	SSHConn  *ssh.Client
	SFTPConn *sftp.Client
)

// GetJSONResponse - function for getting the GET request response in form of JSON
func GetJSONResponse(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if strings.Contains(r.Status, "404") {
		return nil, errors.New(r.Status)
	}
	if strings.Contains(r.Status, "403") {
		time.Sleep(1 * time.Second)
		return GetJSONResponse(url)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// PostJSONResponse - function for getting the POST request response in form of JSON
// value transmitted is a link for a Google URL Shortener
func PostJSONResponse(url, value string) ([]byte, error) {
	var jsonStr = []byte(`{"longUrl": "` + value + `"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if strings.Contains(r.Status, "404") {
		return nil, errors.New(r.Status)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// ExecuteCommand - function which runs the command in SSH session
func ExecuteCommand(command string) (output string, err error) {
	var (
		session *ssh.Session
		wg      sync.WaitGroup
		bOutput []byte
	)

	session, err = SSHConn.NewSession()
	if err != nil {
		glog.Errorf("Failed to create session: %s", err)
		return
	}
	defer session.Close()

	if err = SetupPty(session, false); err != nil {
		glog.Error(err)
		return
	}

	wg.Add(1)
	go func(c string) {
		bOutput, err = session.Output(c)
		output = string(bOutput)
		wg.Done()
	}(command)
	wg.Wait()

	return
}

// ConnectToServer - function for connecting to the server through SSH
func ConnectToServer(user, address string) (err error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			SSHAgent(),
		},
	}
	if SSHConn, err = ssh.Dial("tcp", address, config); err != nil {
		return
	}
	SFTPConn, err = sftp.NewClient(SSHConn)
	return
}

// SSHAgent - function which returns PublicKeysCallback
func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err != nil {
		glog.Errorf("Unable to get publick keys callback: %s", err)
	} else if sshAgent == nil {
		glog.Error("Unable to get publick keys callback: SSH agent is nil")
	} else if err == nil && sshAgent != nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

// SetupPty - function which setups pseudoterminal to server
func SetupPty(session *ssh.Session, piped bool) error {
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if piped {
		stdin, err := session.StdinPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stdin for session: %v", err)
		}
		go io.Copy(stdin, os.Stdin)

		stdout, err := session.StdoutPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stdout for session: %v", err)
		}
		go io.Copy(os.Stdout, stdout)

		stderr, err := session.StderrPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stderr for session: %v", err)
		}
		go io.Copy(os.Stderr, stderr)
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return fmt.Errorf("Request for pseudo terminal failed: %s", err)
	}

	return nil
}

// DownloadFile - function which downloads the file through SFTP from the server
func DownloadFile(srcPath, dstPath string) error {
	srcFile, err := SFTPConn.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	srcFile.WriteTo(dstFile)

	return nil
}

// UploadFile - function which uploads the file through SFTP to the server
func UploadFile(srcPath, dstPath string) error {
	b, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return err
	}

	dstFile, err := SFTPConn.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	dstFile.Write(b)

	return nil
}
