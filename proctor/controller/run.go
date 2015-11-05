package controller

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/shell"
)

func (c *Controller) RunOnVMs(name, command string) error {
	if err := validName(name); err != nil {
		return err
	}

	prefixedName := prefix(name)

	stackStatus, stackID, _, err := c.AWSClient.DescribeStack(prefixedName)
	if err != nil {
		return fmt.Errorf("error querying status: %s", err)
	}

	if stackStatus != "CREATE_COMPLETE" {
		return fmt.Errorf("classroom is not operational (status '%s'), aborting", stackStatus)
	}

	s3URL := c.AWSClient.URLForObject("keys/" + prefixedName)
	pemBytes, err := c.WebClient.Get(s3URL)
	if err != nil {
		return fmt.Errorf("getting SSH key: %s", err)
	}

	hosts, err := c.AWSClient.GetHostsFromStackID(stackID)
	if err != nil {
		return fmt.Errorf("error fetching hosts: %s", err)
	}

	targets := []string{}
	for host, status := range hosts {
		if status == "running" {
			targets = append(targets, host)
		}
	}

	options := &shell.ConnectionOptions{
		Port:          c.SSHPort,
		Username:      c.SSHUser,
		PrivateKeyPEM: pemBytes,
	}
	results := c.ParallelRunner.ConnectAndRun(targets, command, options)

	for host, r := range results {
		c.Log.Println(0, "%s", c.Log.Green("%s", host))
		if len(strings.TrimSpace(r.Stdout)) > 0 {
			c.Log.Println(1, "%s", r.Stdout)
		}
		if r.Error != nil {
			c.Log.Println(1, "%s", c.Log.Red("%s", r.Error))
			err = multierror.Append(err, fmt.Errorf("on %s: %s", host, r.Error))
		}
	}

	return err
}
