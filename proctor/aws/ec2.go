package aws

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (c *Client) CreateKey(name string) (string, error) {
	out, err := c.EC2.CreateKeyPair(&ec2.CreateKeyPairInput{KeyName: aws.String(name)})
	if err != nil {
		return "", err
	}
	if out.KeyName == nil {
		return "", errors.New("CreateKeyPair returned invalid data")
	}

	if *out.KeyName != name {
		return "", fmt.Errorf("tried to create key named '%s' but generated key was called '%s'",
			name, *out.KeyName)
	}

	if out.KeyMaterial == nil || *out.KeyMaterial == "" {
		return "", fmt.Errorf("CreateKeyPair returned an empty key")
	}

	return *out.KeyMaterial, nil
}

func (c *Client) DeleteKey(name string) error {
	_, err := c.EC2.DeleteKeyPair(&ec2.DeleteKeyPairInput{KeyName: aws.String(name)})
	return err
}

func (c *Client) ListKeys(prefix string) ([]string, error) {
	out, err := c.EC2.DescribeKeyPairs(&ec2.DescribeKeyPairsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("key-name"),
				Values: []*string{
					aws.String(prefix + "*"),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	ret := []string{}
	for _, kp := range out.KeyPairs {
		ret = append(ret, *kp.KeyName)
	}
	return ret, nil
}
