package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/aws/aws-sdk-go-v2/config"
)

func GetCredentialFromValue(key string, secret string) credentials.StaticCredentialsProvider {
	provider := credentials.NewStaticCredentialsProvider(key, secret, "")
	return provider
}

func GetConfigWithCredential(cred credentials.StaticCredentialsProvider, region string) aws.Config {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region), config.WithCredentialsProvider(cred))

	if err != nil {
		log.Fatalln(err)
	}

	return sdkConfig

}

func GetConfig(region string) aws.Config {

	sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))

	if err != nil {
		log.Fatalln(err)
	}

	return sdkConfig
}
