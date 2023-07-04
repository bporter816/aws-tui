package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3BucketPolicy struct {
	*Text
	s3Client *s3.Client
	bucket   string
}

func NewS3BucketPolicy(s3Client *s3.Client, bucket string) *S3BucketPolicy {
	s := &S3BucketPolicy{
		Text:     NewText(true, "json"),
		s3Client: s3Client,
		bucket:   bucket,
	}
	return s
}

func (s S3BucketPolicy) GetName() string {
	return fmt.Sprintf("S3 | %v | Bucket Policy", s.bucket)
}

func (s S3BucketPolicy) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s S3BucketPolicy) Render() {
	out, err := s.s3Client.GetBucketPolicy(
		context.TODO(),
		&s3.GetBucketPolicyInput{
			Bucket: aws.String(s.bucket),
		},
	)
	if err != nil {
		panic(err)
	}
	var policy string
	if out.Policy != nil {
		policy = *out.Policy
	}
	s.SetText(policy)
}
