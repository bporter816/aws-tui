package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
)

type S3ObjectTags struct {
	*ui.Table
	repo   *repo.S3
	bucket string
	key    string
	app    *Application
}

func NewS3ObjectTags(repo *repo.S3, bucket string, key string, app *Application) *S3ObjectTags {
	s := &S3ObjectTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo:   repo,
		bucket: bucket,
		key:    key,
		app:    app,
	}
	return s
}

func (s S3ObjectTags) GetService() string {
	return "S3"
}

func (s S3ObjectTags) GetLabels() []string {
	return []string{s.key, "Tags"}
}

func (s S3ObjectTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s S3ObjectTags) Render() {
	model, err := s.repo.ListObjectTags(s.bucket, s.key)
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
