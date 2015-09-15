package controller

import "fmt"

type atlasClient interface {
	GetLatestAMI(string) (string, error)
}

type awsClient interface {
	CreateKey(name string) (string, error)
	DeleteKey(name string) error
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
}

func (c *Controller) CreateClassroom(name string, number int) error {
	c.Log.Println(0, "Looking up latest AMI for %s", c.Log.Green("%s", c.VagrantBoxName))
	ami, err := c.AtlasClient.GetLatestAMI(c.VagrantBoxName)
	if err != nil {
		return err
	}
	c.Log.Println(0, "Found %s", c.Log.Green("%s", ami))

	c.Log.Println(0, "Creating new SSH Keypair for EC2...")
	privateKeyPEMBytes, err := c.AWSClient.CreateKey(fmt.Sprintf("classroom-%s", name))
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
	err := c.AWSClient.DeleteKey(fmt.Sprintf("classroom-%s", name))

	s3Name := "keys/" + name
	c.Log.Println(0, "Deleting private key from S3...")
	err = c.AWSClient.DeleteObject(s3Name)
	return err
}
