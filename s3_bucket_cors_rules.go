package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"strconv"
	"strings"
)

type S3CORSRules struct {
	*ui.Table
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

func (s S3CORSRules) GetService() string {
	return "S3"
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
		var id, allowedMethods, allowedOrigins, allowedHeaders, exposeHeaders, maxAge string
		if v.ID != nil {
			id = *v.ID
		}
		allowedMethods = strings.Join(v.AllowedMethods, ", ")
		allowedOrigins = strings.Join(v.AllowedOrigins, ", ")
		allowedHeaders = strings.Join(v.AllowedHeaders, ", ")
		exposeHeaders = strings.Join(v.ExposeHeaders, ", ")
		maxAge = strconv.Itoa(int(*v.MaxAgeSeconds)) + " sec"
		data = append(data, []string{
			id,
			allowedMethods,
			allowedOrigins,
			allowedHeaders,
			exposeHeaders,
			maxAge,
		})
	}
	s.SetData(data)
}
