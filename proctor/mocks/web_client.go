package mocks

type WebClient struct {
	GetCall struct {
		Receives struct {
			URL string
		}
		Returns struct {
			Body  []byte
			Error error
		}
	}
}

func (c *WebClient) Get(url string) ([]byte, error) {
	c.GetCall.Receives.URL = url
	return c.GetCall.Returns.Body, c.GetCall.Returns.Error
}
