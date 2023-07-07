package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3BucketTags struct {
	*Table
	s3Client *s3.Client
	bucket   string
	app      *Application
}

func NewS3BucketTags(s3Client *s3.Client, bucket string, app *Application) *S3BucketTags {
	s := &S3BucketTags{
		Table: NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		s3Client: s3Client,
		bucket:   bucket,
		app:      app,
	}
	return s
}

func (s S3BucketTags) GetName() string {
	return fmt.Sprintf("S3 | Buckets | %v | Tags", s.bucket)
}

func (s S3BucketTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s S3BucketTags) Render() {
	out, err := s.s3Client.GetBucketTagging(
		context.TODO(),
		&s3.GetBucketTaggingInput{
			Bucket: aws.String(s.bucket),
		},
	)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range out.TagSet {
		data = append(data, []string{
			*v.Key,
			*v.Value,
		})
	}
	s.SetData(data)
}
