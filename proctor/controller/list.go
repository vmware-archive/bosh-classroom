package controller

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (c *Controller) ListClassrooms(format string) (string, error) {
	keys, err := c.AWSClient.ListKeys("classroom-")
	if err != nil {
		return "", err
	}
	for i, k := range keys {
		keys[i] = strings.TrimPrefix(k, "classroom-")
	}

	if format == "json" {
		jsonBytes, err := json.MarshalIndent(keys, "", "    ")
		return string(jsonBytes), err
	}
	if format == "plain" {
		return strings.Join(keys, "\n"), nil
	}
	return "", fmt.Errorf("expected format to be either 'json' or 'plain'")
}
