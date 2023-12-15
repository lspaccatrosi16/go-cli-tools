package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/lspaccatrosi16/go-cli-tools/internal/pkgError"

	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/aws/aws-sdk-go-v2/config"
)

var wrap = pkgError.WrapErrorFactory("aws")

func GetCredentialFromValue(key string, secret string) credentials.StaticCredentialsProvider {
	provider := credentials.NewStaticCredentialsProvider(key, secret, "")
	return provider
}

func GetConfigWithCredential(cred credentials.StaticCredentialsProvider, region string) (*aws.Config, error) {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region), config.WithCredentialsProvider(cred))

	if err != nil {
		return nil, wrap(err)
	}

	return &sdkConfig, nil

}

func GetConfig(region string) (*aws.Config, error) {

	sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))

	if err != nil {
		return nil, wrap(err)
	}

	return &sdkConfig, nil
}
