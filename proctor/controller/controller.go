package controller

import (
	"encoding/json"
	"fmt"
	"strings"
)

type atlasClient interface {
	GetLatestAMIs(string) (map[string]string, error)
}

type awsClient interface {
	CreateKey(name string) (string, error)
	DeleteKey(name string) error
	ListKeys(prefix string) ([]string, error)
	StoreObject(name string, bytes []byte, downloadFileName, contentType string) error
	DeleteObject(name string) error
	URLForObject(name string) string
}

type cliLogger interface {
	Println(indentation int, format string, args ...interface{})
	Green(format string, args ...interface{}) string
}

type Controller struct {
	AtlasClient atlasClient
	AWSClient   awsClient
	Log         cliLogger

	VagrantBoxName string
	Region         string
}

func keypairName(classroomName string) string {
	return "classroom-" + classroomName
}

func (c *Controller) CreateClassroom(name string, number int) error {
	c.Log.Println(0, "Looking up latest AMI for %s", c.Log.Green("%s", c.VagrantBoxName))
	amiMap, err := c.AtlasClient.GetLatestAMIs(c.VagrantBoxName)
	if err != nil {
		return err
	}

	ami, ok := amiMap[c.Region]
	if !ok {
		return fmt.Errorf("Couldn't find AMI in region %s", c.Region)
	}
	c.Log.Println(0, "Found %s", c.Log.Green("%s", ami))

	keypair := keypairName(name)
	c.Log.Println(0, "Creating SSH Keypair %s", c.Log.Green("%s", keypair))
	privateKeyPEMBytes, err := c.AWSClient.CreateKey(keypair)
	if err != nil {
		return err
	}

	s3Name := "keys/" + name
	s3URL := c.AWSClient.URLForObject(s3Name)
	c.Log.Println(0, "Uploading private key to %s", c.Log.Green("%s", s3URL))
	err = c.AWSClient.StoreObject(
		s3Name, []byte(privateKeyPEMBytes),
		"bosh101_ssh_key.pem", "application/x-pem-file")
	return err
}

func (c *Controller) DestroyClassroom(name string) error {
	c.Log.Println(0, "Deleting classroom keypair...")
	err := c.AWSClient.DeleteKey(keypairName(name))

	s3Name := "keys/" + name
	c.Log.Println(0, "Deleting private key from S3...")
	err = c.AWSClient.DeleteObject(s3Name)
	return err
}

func (c *Controller) ListClassrooms(format string) (string, error) {
	keys, err := c.AWSClient.ListKeys("classroom-")
	if err != nil {
		return "", err
	}
	for i, k := range keys {
		keys[i] = strings.TrimPrefix(k, "classroom-")
	}

	if format == "json" {
		jsonBytes, err := json.MarshalIndent(keys, "", "    ")
		return string(jsonBytes), err
	}
	if format == "plain" {
		return strings.Join(keys, "\n"), nil
	}
	return "", fmt.Errorf("expected format to be either 'json' or 'plain'")
}
