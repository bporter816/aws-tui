package internal

import (
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/view"
)

type S3BucketPolicy struct {
	*ui.Text
	view.S3
	repo   *repo.S3
	bucket string
	app    *Application
}

func NewS3BucketPolicy(repo *repo.S3, bucket string, app *Application) *S3BucketPolicy {
	s := &S3BucketPolicy{
		Text:   ui.NewText(true, "json"),
		repo:   repo,
		bucket: bucket,
		app:    app,
	}
	return s
}

func (s S3BucketPolicy) GetLabels() []string {
	return []string{s.bucket, "Bucket Policy"}
}

func (s S3BucketPolicy) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s S3BucketPolicy) Render() {
	policy, err := s.repo.GetBucketPolicy(s.bucket)
	if err != nil {
		panic(err)
	}

	s.SetText(policy)
}
