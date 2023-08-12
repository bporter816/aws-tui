package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bporter816/aws-tui/ui"
	"io"
	"strings"
)

type S3Object struct {
	*ui.Text
	s3Client *s3.Client
	bucket   string
	key      string
	app      *Application
}

func NewS3Object(s3Client *s3.Client, bucket string, key string, app *Application) *S3Object {
	s := &S3Object{
		Text:     ui.NewText(false, ""),
		s3Client: s3Client,
		bucket:   bucket,
		key:      key,
		app:      app,
	}
	return s
}

func (s S3Object) GetService() string {
	return "S3"
}

func (s S3Object) GetLabels() []string {
	return []string{s.key}
}

func (s S3Object) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s S3Object) Render() {
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
	defer out.Body.Close()
	b := make([]byte, out.ContentLength)
	n, err := out.Body.Read(b)
	if err != nil && err != io.EOF {
		panic(err)
	}
	split := strings.Split(s.key, ".")
	if len(split) > 1 {
		// TODO abstract this into the text view
		s.Text.SetDynamicColors(true)
		s.Text.HighlightSyntax = true
		s.Text.Lang = split[len(split)-1]
	}
	s.SetText(string(b[0:n]))
}
