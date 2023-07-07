package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gdamore/tcell/v2"
)

type S3Buckets struct {
	*Table
	s3Client *s3.Client
	app      *Application
}

func NewS3Buckets(s3Client *s3.Client, app *Application) *S3Buckets {
	s := &S3Buckets{
		Table: NewTable([]string{
			"NAME",
			"CREATED",
		}, 1, 0),
		s3Client: s3Client,
		app:      app,
	}
	s.SetSelectedFunc(s.selectHandler)
	return s
}

func (s S3Buckets) GetName() string {
	return "S3 | Buckets"
}

func (s S3Buckets) selectHandler(row, col int) {
	bucket, err := s.GetColSelection("NAME")
	if err != nil {
		return
	}
	objectsView := NewS3Objects(s.s3Client, bucket, s.app)
	s.app.AddAndSwitch(objectsView)
}

func (s S3Buckets) bucketPolicyHandler() {
	bucket, err := s.GetColSelection("NAME")
	if err != nil {
		return
	}
	policyView := NewS3BucketPolicy(s.s3Client, bucket)
	s.app.AddAndSwitch(policyView)
}

func (s S3Buckets) tagsHandler() {
	bucket, err := s.GetColSelection("NAME")
	if err != nil {
		return
	}
	tagsView := NewS3BucketTags(s.s3Client, bucket, s.app)
	s.app.AddAndSwitch(tagsView)
}

func (s S3Buckets) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
			Description: "Bucket Policy",
			Action:      s.bucketPolicyHandler,
		},
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      s.tagsHandler,
		},
	}
}

func (s S3Buckets) Render() {
	var buckets []s3Types.Bucket
	out, err := s.s3Client.ListBuckets(
		context.TODO(),
		&s3.ListBucketsInput{},
	)
	if err != nil {
		panic(err)
	}
	buckets = out.Buckets

	var data [][]string
	for _, v := range buckets {
		var name, created string
		if v.Name != nil {
			name = *v.Name
		}
		if v.CreationDate != nil {
			created = v.CreationDate.Format("2006-01-02 15:04:05")
		}
		data = append(data, []string{
			name,
			created,
		})
	}
	s.SetData(data)
}
