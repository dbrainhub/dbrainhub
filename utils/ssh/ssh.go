package ssh

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
)

type SSHClient struct {
	inner *ssh.Client
}

// NewSSHClient create an SSHClient
func NewSSHClient(inner *ssh.Client) *SSHClient {
	return &SSHClient{inner: inner}
}

// Connect connects to remote host and returns an SSHClient.
// `addr` should be of the form "ip:port".
func Connect(addr string, user string, priKey []byte) (*SSHClient, error) {
	signer, err := ssh.ParsePrivateKey(priKey)
	if err != nil {
		return nil, err
	}
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}

	return &SSHClient{inner: sshClient}, nil
}

// Close closes SSHClient. Should be called every time when finished.
func (client *SSHClient) Close() error {
	return client.inner.Close()
}

// Run runs cmd on the remote host and returns its stdout and stderr combined as []byte.
func (client *SSHClient) Run(cmd string) ([]byte, error) {
	ss, err := client.inner.NewSession()
	if err != nil {
		return nil, err
	}
	defer ss.Close()
	return ss.CombinedOutput(cmd)
}

// Scp copies local file to remote.
func (client *SSHClient) Scp(localPath, remotePath string) error {
	sftpClient, err := sftp.NewClient(client.inner)
	if err != nil {
		return err
	}
	remoteFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return err
	}
	localFile, err := os.Open(localPath)
	if err != nil {
		return err
	}
	_, err = io.Copy(remoteFile, localFile)
	return err
}

// Chmod changes file/directory permission.
func (client *SSHClient) Chmod(remotePath string, mode os.FileMode) error {
	sftpClient, err := sftp.NewClient(client.inner)
	if err != nil {
		return err
	}
	return sftpClient.Chmod(remotePath, mode)
}
