package controller

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
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
	CreateStack(name string, template string, parameters map[string]string) (string, error)
	DeleteStack(name string) error
	GetStackStatus(name string) (string, error)
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
	Template       string
}

func prefix(classroomName string) string {
	return "classroom-" + classroomName
}

func (c *Controller) CreateClassroom(name string, number int) error {
	const requiredPattern = `^[a-zA-Z][-a-zA-Z0-9]*$`
	regex := regexp.MustCompile(requiredPattern)
	if !regex.MatchString(name) {
		return fmt.Errorf("invalid name: must match pattern %s", requiredPattern)
	}

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

	prefixedName := prefix(name)
	c.Log.Println(0, "Creating SSH Keypair %s", c.Log.Green("%s", prefixedName))
	privateKeyPEMBytes, err := c.AWSClient.CreateKey(prefixedName)
	if err != nil {
		return err
	}

	s3Name := "keys/" + prefixedName
	s3URL := c.AWSClient.URLForObject(s3Name)
	c.Log.Println(0, "Uploading private key to %s", c.Log.Green("%s", s3URL))
	err = c.AWSClient.StoreObject(
		s3Name, []byte(privateKeyPEMBytes),
		"bosh101_ssh_key.pem", "application/x-pem-file")
	if err != nil {
		return err
	}

	c.Log.Println(0, "Creating CloudFormation stack %s", c.Log.Green("%s", prefixedName))
	_, err = c.AWSClient.CreateStack(prefixedName, c.Template, map[string]string{
		"AMI":           ami,
		"KeyName":       prefixedName,
		"InstanceCount": strconv.Itoa(number),
	})

	return err
}

func (c *Controller) DestroyClassroom(name string) error {
	prefixedName := prefix(name)

	c.Log.Println(0, "Deleting CloudFormation stack %s", c.Log.Green("%s", prefixedName))
	err := c.AWSClient.DeleteStack(prefixedName)
	if err != nil {
		return err
	}

	c.Log.Println(0, "Deleting classroom keypair...")
	err = c.AWSClient.DeleteKey(prefixedName)
	if err != nil {
		return err
	}

	s3Name := "keys/" + prefixedName
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

func (c *Controller) DescribeClassroom(name, format string) (string, error) {
	prefixedName := prefix(name)

	status, err := c.AWSClient.GetStackStatus(prefixedName)
	if err != nil {
		return "", err
	}

	if format == "json" {
		var description struct {
			Status string `json:"status"`
		}
		description.Status = status
		descriptionBytes, err := json.Marshal(description)
		return string(descriptionBytes), err
	}
	if format == "plain" {
		return fmt.Sprintf("%s: %s", "status", status), nil
	}
	return "", fmt.Errorf("expected format to be either 'json' or 'plain'")
}
