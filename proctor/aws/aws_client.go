package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Endpoints struct {
	Route53        string
	EC2            string
	S3             string
	Cloudformation string
}

type Client struct {
	EC2            EC2Client
	S3             S3Client
	Route53        Route53Client
	Cloudformation CloudformationClient
	// HostedZoneID   string
	// HostedZoneName string
	Bucket string
}

type S3Client interface {
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
	DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
}

type EC2Client interface {
	CreateKeyPair(*ec2.CreateKeyPairInput) (*ec2.CreateKeyPairOutput, error)
	DeleteKeyPair(*ec2.DeleteKeyPairInput) (*ec2.DeleteKeyPairOutput, error)
	DescribeKeyPairs(*ec2.DescribeKeyPairsInput) (*ec2.DescribeKeyPairsOutput, error)
	DescribeInstances(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error)
}

type Route53Client interface {
	ChangeResourceRecordSets(input *route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error)
	ListResourceRecordSets(input *route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput, error)
}

type CloudformationClient interface {
	CreateStack(*cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error)
	DeleteStack(*cloudformation.DeleteStackInput) (*cloudformation.DeleteStackOutput, error)
	DescribeStacks(*cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error)
}

type AWSError struct {
	Method string
	Err    error
}

func (e *AWSError) Error() string {
	return fmt.Sprintf("%s: %s", e.Method, e.Err)
}

type Config struct {
	AccessKey  string
	SecretKey  string
	RegionName string
	// HostedZoneID      string
	// HostedZoneName    string
	Bucket            string
	EndpointOverrides *Endpoints
}

func New(config Config) *Client {
	credentials := credentials.NewStaticCredentials(config.AccessKey, config.SecretKey, "")
	sdkConfig := &aws.Config{
		Credentials: credentials,
		Region:      aws.String(config.RegionName),
	}

	endpointOverrides := config.EndpointOverrides
	if endpointOverrides == nil {
		endpointOverrides = &Endpoints{}
	}

	route53Client := route53.New(sdkConfig.Merge(&aws.Config{MaxRetries: aws.Int(7), Endpoint: aws.String(endpointOverrides.Route53)}))
	ec2Client := ec2.New(sdkConfig.Merge(&aws.Config{MaxRetries: aws.Int(7), Endpoint: aws.String(endpointOverrides.EC2)}))
	s3Client := s3.New(sdkConfig.Merge(&aws.Config{MaxRetries: aws.Int(7), Endpoint: aws.String(endpointOverrides.S3), S3ForcePathStyle: aws.Bool(true)}))
	cloudformationClient := cloudformation.New(sdkConfig.Merge(&aws.Config{MaxRetries: aws.Int(7), Endpoint: aws.String(endpointOverrides.Cloudformation)}))

	return &Client{
		EC2:            ec2Client,
		S3:             s3Client,
		Route53:        route53Client,
		Cloudformation: cloudformationClient,
		// HostedZoneID:   config.HostedZoneID,
		// HostedZoneName: config.HostedZoneName,
		Bucket: config.Bucket,
	}
}

func toStringPointers(strings ...string) []*string {
	var output []*string
	for _, s := range strings {
		output = append(output, aws.String(s))
	}
	return output
}
