package main

import (
	"strconv"
	"strings"

	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
)

type S3CORSRules struct {
	*ui.Table
	view.S3
	repo   *repo.S3
	bucket string
	app    *Application
}

func NewS3CORSRules(repo *repo.S3, bucket string, app *Application) *S3CORSRules {
	s := &S3CORSRules{
		Table: ui.NewTable([]string{
			"ID",
			"ALLOWED METHODS",
			"ALLOWED ORIGINS",
			"ALLOWED HEADERS",
			"EXPOSE HEADERS",
			"MAX AGE",
		}, 1, 0),
		repo:   repo,
		bucket: bucket,
		app:    app,
	}
	return s
}

func (s S3CORSRules) GetLabels() []string {
	return []string{s.bucket, "CORS"}
}

func (s S3CORSRules) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s S3CORSRules) Render() {
	model, err := s.repo.GetCORSRules(s.bucket)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			utils.DerefString(v.ID, ""),
			strings.Join(v.AllowedMethods, ", "),
			strings.Join(v.AllowedOrigins, ", "),
			strings.Join(v.AllowedHeaders, ", "),
			strings.Join(v.ExposeHeaders, ", "),
			strconv.Itoa(int(*v.MaxAgeSeconds)) + " sec",
		})
	}
	s.SetData(data)
}
