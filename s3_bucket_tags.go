package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
)

type S3BucketTags struct {
	*ui.Table
	view.S3
	repo   *repo.S3
	bucket string
	app    *Application
}

func NewS3BucketTags(repo *repo.S3, bucket string, app *Application) *S3BucketTags {
	s := &S3BucketTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo:   repo,
		bucket: bucket,
		app:    app,
	}
	return s
}

func (s S3BucketTags) GetLabels() []string {
	return []string{s.bucket, "Tags"}
}

func (s S3BucketTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s S3BucketTags) Render() {
	model, err := s.repo.ListBucketTags(s.bucket)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			v.Key,
			v.Value,
		})
	}
	s.SetData(data)
}
