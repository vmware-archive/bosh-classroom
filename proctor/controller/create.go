package controller

import (
	"fmt"
	"regexp"
	"strconv"
)

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
