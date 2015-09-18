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
	ListKeysCall struct {
		Receives struct {
			Prefix string
		}
		Returns struct {
			Keys  []string
			Error error
		}
	}

	StoreObjectCall struct {
		Receives struct {
			Name             string
			Bytes            []byte
			DownloadFileName string
			ContentType      string
		}
		Returns struct {
			Error error
		}
	}
	DeleteObjectCall struct {
		Receives struct {
			Name string
		}
		Returns struct {
			Error error
		}
	}

	URLForObjectCall struct {
		Receives struct {
			Name string
		}

		Returns struct {
			URL string
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
func (c *AWSClient) ListKeys(prefix string) ([]string, error) {
	c.ListKeysCall.Receives.Prefix = prefix
	return c.ListKeysCall.Returns.Keys, c.ListKeysCall.Returns.Error
}

func (c *AWSClient) StoreObject(name string, bytes []byte,
	downloadFileName string, contentType string) error {
	c.StoreObjectCall.Receives.Name = name
	c.StoreObjectCall.Receives.Bytes = bytes
	c.StoreObjectCall.Receives.DownloadFileName = downloadFileName
	c.StoreObjectCall.Receives.ContentType = contentType
	return c.StoreObjectCall.Returns.Error
}

func (c *AWSClient) DeleteObject(name string) error {
	c.DeleteObjectCall.Receives.Name = name
	return c.DeleteObjectCall.Returns.Error
}

func (c *AWSClient) URLForObject(name string) string {
	c.URLForObjectCall.Receives.Name = name
	return c.URLForObjectCall.Returns.URL
}
