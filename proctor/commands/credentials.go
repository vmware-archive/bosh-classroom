package commands

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

func CredentialsFromEnv() (*credentials.Credentials, error) {
	id, secret := os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY")
	if id == "" || secret == "" {
		return nil, errors.New("missing AWS credentials, check your AWS_* environment variables")
	}

	return credentials.NewStaticCredentials(id, secret, ""), nil
}
