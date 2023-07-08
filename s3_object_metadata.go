package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3ObjectMetadata struct {
	*Table
	s3Client *s3.Client
	bucket   string
	key      string
	app      *Application
}

func NewS3ObjectMetadata(s3Client *s3.Client, bucket string, key string, app *Application) *S3ObjectMetadata {
	s := &S3ObjectMetadata{
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

func (s S3ObjectMetadata) GetService() string {
	return "S3"
}

func (s S3ObjectMetadata) GetLabels() []string {
	return []string{s.key, "Metadata"}
}

func (s S3ObjectMetadata) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s S3ObjectMetadata) Render() {
	out, err := s.s3Client.GetObject(
		context.TODO(),
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(s.key),
		},
	)
	if err != nil {
		panic(err)
	}

	var data [][]string
	if out.ContentType != nil {
		data = append(data, []string{
			"Content-Type",
			*out.ContentType,
		})
	}
	for k, v := range out.Metadata {
		data = append(data, []string{
			k,
			v,
		})
	}
	s.SetData(data)
}
