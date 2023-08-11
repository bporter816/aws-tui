package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bporter816/aws-tui/ui"
)

type S3BucketPolicy struct {
	*ui.Text
	s3Client *s3.Client
	bucket   string
	app      *Application
}

func NewS3BucketPolicy(s3Client *s3.Client, bucket string, app *Application) *S3BucketPolicy {
	s := &S3BucketPolicy{
		Text:     ui.NewText(true, "json"),
		s3Client: s3Client,
		bucket:   bucket,
		app:      app,
	}
	return s
}

func (s S3BucketPolicy) GetService() string {
	return "S3"
}

func (s S3BucketPolicy) GetLabels() []string {
	return []string{s.bucket, "Bucket Policy"}
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
