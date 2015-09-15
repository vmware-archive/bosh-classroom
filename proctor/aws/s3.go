package aws

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (c *Client) StoreObject(name string, dataBytes []byte,
	downloadFileName string, contentType string) error {
	_, err := c.S3.PutObject(&s3.PutObjectInput{
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(dataBytes),
		Bucket:               aws.String(c.Bucket),
		ContentDisposition:   aws.String(fmt.Sprintf("attachment; filename=%s;", downloadFileName)),
		ContentType:          aws.String(contentType),
		Key:                  aws.String(name),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}

func (c *Client) DeleteObject(name string) error {
	_, err := c.S3.DeleteObject(&s3.DeleteObjectInput{
		Key:    aws.String(name),
		Bucket: aws.String(c.Bucket),
	})
	return err
}

func (c *Client) URLForObject(name string) string {
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", c.Bucket, name)
}
