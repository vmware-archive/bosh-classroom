package controller

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
