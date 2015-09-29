package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (c *Controller) DescribeClassroom(name, format string) (string, error) {
	prefixedName := prefix(name)

	status, stackID, parameters, err := c.AWSClient.DescribeStack(prefixedName)
	if err != nil {
		return "", err
	}

	keyURL := c.AWSClient.URLForObject("keys/" + prefixedName)

	var description struct {
		Status string            `json:"status"`
		Number int               `json:"number"`
		SSHKey string            `json:"ssh_key"`
		Hosts  map[string]string `json:"hosts"`
	}
	description.Status = status
	description.SSHKey = keyURL
	description.Number, err = strconv.Atoi(parameters["InstanceCount"])
	if err != nil {
		return "", errors.New("malformed CloudFormation stack: missing or invalid parameter 'InstanceCount'")
	}
	description.Hosts, err = c.AWSClient.GetHostsFromStackID(stackID)
	if err != nil {
		return "", fmt.Errorf("error fetching hosts for stack: %s", err)
	}

	if format == "json" {
		descriptionBytes, err := json.MarshalIndent(description, "", "    ")
		return string(descriptionBytes), err
	}
	if format == "plain" {
		hosts := []string{}
		for k, v := range description.Hosts {
			hosts = append(hosts, fmt.Sprintf("%s\t%s", k, v))
		}
		return fmt.Sprintf("%s: %s\n%s: %d\n%s: %s\n%s:\n%s",
			"status", description.Status,
			"number", description.Number,
			"ssh_key", description.SSHKey,
			"hosts", strings.Join(hosts, "\n"),
		), nil
	}
	return "", fmt.Errorf("expected format to be either 'json' or 'plain'")
}
