package main

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/gdamore/tcell/v2"
)

type S3Buckets struct {
	*ui.Table
	repo     *repo.S3
	s3Client *s3.Client
	app      *Application
}

func NewS3Buckets(repo *repo.S3, s3Client *s3.Client, app *Application) *S3Buckets {
	s := &S3Buckets{
		Table: ui.NewTable([]string{
			"NAME",
			"CREATED",
		}, 1, 0),
		repo:     repo,
		s3Client: s3Client,
		app:      app,
	}
	s.SetSelectedFunc(s.selectHandler)
	return s
}

func (s S3Buckets) GetService() string {
	return "S3"
}

func (s S3Buckets) GetLabels() []string {
	return []string{"Buckets"}
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
	policyView := NewS3BucketPolicy(s.s3Client, bucket, s.app)
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
	model, err := s.repo.ListBuckets()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var name, created string
		if v.Name != nil {
			name = *v.Name
		}
		if v.CreationDate != nil {
			created = v.CreationDate.Format(utils.DefaultTimeFormat)
		}
		data = append(data, []string{
			name,
			created,
		})
	}
	s.SetData(data)
}
