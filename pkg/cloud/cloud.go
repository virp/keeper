// Package cloud provides a convenience functions for interacting with cloud services.
package cloud

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// NewS3Client constructing AWS S3 client.
func NewS3Client(key string, secret string, region string, endpoint string) *s3.Client {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(s, r string, o ...interface{}) (aws.Endpoint, error) {
		if s == s3.ServiceID && region == r && endpoint != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           endpoint,
				SigningRegion: region,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	config := aws.Config{
		Credentials:                 credentials.NewStaticCredentialsProvider(key, secret, ""),
		Region:                      region,
		EndpointResolverWithOptions: customResolver,
	}

	return s3.NewFromConfig(
		config,
		func(o *s3.Options) {
			o.UsePathStyle = true
		},
	)
}
