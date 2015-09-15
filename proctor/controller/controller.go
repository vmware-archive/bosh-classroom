package controller

import "fmt"

type atlasClient interface {
	GetLatestAMI(string) (string, error)
}

type awsClient interface {
	CreateKey(string) (string, error)
	DeleteKey(string) error
}

type cliLogger interface {
	Println(indentation int, format string, args ...interface{})
	Green(format string, args ...interface{}) string
}

type Controller struct {
	AtlasClient atlasClient
	AwsClient   awsClient
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
	_, err = c.AwsClient.CreateKey(fmt.Sprintf("classroom-%s", name))
	return err
}

func (c *Controller) DestroyClassroom(name string) error {
	c.Log.Println(0, "Deleting classroom keypair...")
	err := c.AwsClient.DeleteKey(fmt.Sprintf("classroom-%s", name))
	return err
}
