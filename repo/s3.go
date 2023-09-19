package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bporter816/aws-tui/model"
)

type S3 struct {
	s3Client *s3.Client
}

func NewS3(s3Client *s3.Client) *S3 {
	return &S3{
		s3Client: s3Client,
	}
}

func (s S3) ListBuckets() ([]model.S3Bucket, error) {
	out, err := s.s3Client.ListBuckets(
		context.TODO(),
		&s3.ListBucketsInput{},
	)
	if err != nil {
		return []model.S3Bucket{}, err
	}
	var buckets []model.S3Bucket
	for _, v := range out.Buckets {
		buckets = append(buckets, model.S3Bucket(v))
	}
	return buckets, nil
}

func (s S3) GetBucketPolicy(bucketName string) (string, error) {
	out, err := s.s3Client.GetBucketPolicy(
		context.TODO(),
		&s3.GetBucketPolicyInput{
			Bucket: aws.String(bucketName),
		},
	)
	if err != nil || out.Policy == nil {
		return "", err
	}
	return *out.Policy, nil
}
