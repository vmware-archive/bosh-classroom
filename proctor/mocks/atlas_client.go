package mocks

type AtlasClient struct {
	GetLatestAMICall struct {
		Receives struct {
			BoxName string
		}
		Returns struct {
			AMI   string
			Error error
		}
	}
}

func (c *AtlasClient) GetLatestAMI(boxName string) (string, error) {
	c.GetLatestAMICall.Receives.BoxName = boxName
	return c.GetLatestAMICall.Returns.AMI, c.GetLatestAMICall.Returns.Error
}
