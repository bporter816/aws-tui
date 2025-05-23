package internal

import (
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/view"
)

type S3ObjectMetadata struct {
	*ui.Table
	view.S3
	repo   *repo.S3
	bucket string
	key    string
	app    *Application
}

func NewS3ObjectMetadata(repo *repo.S3, bucket string, key string, app *Application) *S3ObjectMetadata {
	s := &S3ObjectMetadata{
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

func (s S3ObjectMetadata) GetLabels() []string {
	return []string{s.key, "Metadata"}
}

func (s S3ObjectMetadata) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s S3ObjectMetadata) Render() {
	model, err := s.repo.GetObjectMetadata(s.bucket, s.key)
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
