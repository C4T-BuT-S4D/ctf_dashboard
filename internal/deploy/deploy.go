package deploy

import (
	"ctf_dashboard/internal/common"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"strings"
)

func ValidPublicKey(key string) error {
	_, _, _, _, err := ssh.ParseAuthorizedKey([]byte(key))
	return err
}

func UploadSSHKey(vulnbox common.Vulnbox, key string, keyfile string) error {
	err := ValidPublicKey(key)
	if err != nil {
		return err
	}

	key = strings.TrimSpace(key) + "\n"

	privateKeyData, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return err
	}
	privateKey, err := ssh.ParsePrivateKey(privateKeyData)
	if err != nil {
		return err
	}
	config := &ssh.ClientConfig{
		User: vulnbox.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(privateKey),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", vulnbox.Host+":22", config)
	if err != nil {
		return nil
	}
	defer func() {
		if err := client.Close(); err != nil {
			logrus.Errorf("Error closing client: %v", err)
		}
	}()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer func() {
		if err := sftpClient.Close(); err != nil {
			logrus.Errorf("Error closing client: %v", err)
		}
	}()

	authkeys, err := sftpClient.OpenFile(".ssh/authorized_keys", os.O_APPEND|os.O_CREATE|os.O_WRONLY)
	if err != nil {
		return err
	}

	if _, err = authkeys.Write([]byte(key)); err != nil {
		return err
	}
	return nil
}
