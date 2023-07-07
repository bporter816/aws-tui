package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3ObjectTags struct {
	*Table
	s3Client *s3.Client
	bucket   string
	key      string
	app      *Application
}

func NewS3ObjectTags(s3Client *s3.Client, bucket string, key string, app *Application) *S3ObjectTags {
	s := &S3ObjectTags{
		Table: NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		s3Client: s3Client,
		bucket:   bucket,
		key:      key,
		app:      app,
	}
	return s
}

func (s S3ObjectTags) GetName() string {
	return fmt.Sprintf("S3 | Buckets | %v | %v | Tags", s.bucket, s.key)
}

func (s S3ObjectTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s S3ObjectTags) Render() {
	out, err := s.s3Client.GetObjectTagging(
		context.TODO(),
		&s3.GetObjectTaggingInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(s.key),
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
