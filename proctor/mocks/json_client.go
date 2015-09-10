package mocks

import "encoding/json"

type JSONClient struct {
	GetCall struct {
		Args struct {
			Route        string
			ResponseData interface{}
		}
		Return struct {
			Error error
		}

		ResponseJSON string
	}
}

func (c *JSONClient) Get(route string, responseData interface{}) error {
	c.GetCall.Args.Route = route
	c.GetCall.Args.ResponseData = responseData

	json.Unmarshal([]byte(c.GetCall.ResponseJSON), responseData)

	return c.GetCall.Return.Error
}
