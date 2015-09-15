package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Endpoints struct {
	Route53 string
	EC2     string
	S3      string
}

type Client struct {
	EC2            EC2Client
	S3             S3Client
	Route53        Route53Client
	HostedZoneID   string
	HostedZoneName string
	Bucket         string
}

type S3Client interface {
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
}

type EC2Client interface {
	CreateKeyPair(*ec2.CreateKeyPairInput) (*ec2.CreateKeyPairOutput, error)
	DeleteKeyPair(*ec2.DeleteKeyPairInput) (*ec2.DeleteKeyPairOutput, error)
}

type Route53Client interface {
	ChangeResourceRecordSets(input *route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error)
	ListResourceRecordSets(input *route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput, error)
}

type AWSError struct {
	Method string
	Err    error
}

func (e *AWSError) Error() string {
	return fmt.Sprintf("%s: %s", e.Method, e.Err)
}

type Config struct {
	AccessKey         string
	SecretKey         string
	RegionName        string
	HostedZoneID      string
	HostedZoneName    string
	RecoveryBucket    string
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

	return &Client{
		EC2:            ec2Client,
		S3:             s3Client,
		Route53:        route53Client,
		HostedZoneID:   config.HostedZoneID,
		HostedZoneName: config.HostedZoneName,
		Bucket:         config.RecoveryBucket,
	}
}

func toStringPointers(strings ...string) []*string {
	var output []*string
	for _, s := range strings {
		output = append(output, aws.String(s))
	}
	return output
}
