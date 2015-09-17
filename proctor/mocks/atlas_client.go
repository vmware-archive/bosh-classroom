package mocks

type AtlasClient struct {
	GetLatestAMIsCall struct {
		Receives struct {
			BoxName string
		}
		Returns struct {
			AMIMap map[string]string
			Error  error
		}
	}
}

func (c *AtlasClient) GetLatestAMIs(boxName string) (map[string]string, error) {
	c.GetLatestAMIsCall.Receives.BoxName = boxName
	return c.GetLatestAMIsCall.Returns.AMIMap, c.GetLatestAMIsCall.Returns.Error
}
