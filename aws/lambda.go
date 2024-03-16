package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

type Lambda struct {
	LambdaClient *lambda.Client
}

func (l Lambda) UpdateFunctionCode(arn string, bucket string, key string) error {

	_, err := l.LambdaClient.UpdateFunctionCode(context.TODO(), &lambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(arn),
		S3Bucket:     aws.String(bucket),
		S3Key:        aws.String(key),
	})

	if err != nil {
		return wrap(err)
	}

	return nil
}

func NewLambda(sdkConfig aws.Config) Lambda {
	lambdaClient := lambda.NewFromConfig(sdkConfig)
	lambda := Lambda{LambdaClient: lambdaClient}
	return lambda
}
