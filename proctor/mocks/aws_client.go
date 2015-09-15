package mocks

type AWSClient struct {
	CreateKeyCall struct {
		Receives struct {
			KeyName string
		}
		Returns struct {
			PrivateKeyPEM string
			Error         error
		}
	}
	DeleteKeyCall struct {
		Receives struct {
			KeyName string
		}
		Returns struct {
			Error error
		}
	}
}

func (c *AWSClient) CreateKey(keyName string) (string, error) {
	c.CreateKeyCall.Receives.KeyName = keyName
	return c.CreateKeyCall.Returns.PrivateKeyPEM, c.CreateKeyCall.Returns.Error
}
func (c *AWSClient) DeleteKey(keyName string) error {
	c.DeleteKeyCall.Receives.KeyName = keyName
	return c.DeleteKeyCall.Returns.Error
}
