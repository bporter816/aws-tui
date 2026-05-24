package internal

import (
	"strings"

	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
)

type S3Object struct {
	*ui.Text
	repo   *repo.S3
	bucket string
	key    string
	app    *Application
}

func NewS3Object(repo *repo.S3, bucket string, key string, app *Application) *S3Object {
	s := &S3Object{
		Text:   ui.NewText(false, ""),
		repo:   repo,
		bucket: bucket,
		key:    key,
		app:    app,
	}
	return s
}

func (s S3Object) GetService() string { return "S3" }

func (s S3Object) GetLabels() []string {
	return []string{s.key}
}

func (s S3Object) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s S3Object) Render() error {
	b, err := s.repo.GetObject(s.bucket, s.key)
	if err != nil {
		return err
	}

	split := strings.Split(s.key, ".")
	if len(split) > 1 {
		// TODO abstract this into the text view
		s.Text.SetDynamicColors(true)
		s.Text.HighlightSyntax = true
		s.Text.Lang = split[len(split)-1]
	}
	s.SetText(string(b))
	return nil
}
