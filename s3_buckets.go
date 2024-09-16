package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type S3Buckets struct {
	*ui.Table
	view.S3
	repo *repo.S3
	app  *Application
}

func NewS3Buckets(repo *repo.S3, app *Application) *S3Buckets {
	s := &S3Buckets{
		Table: ui.NewTable([]string{
			"NAME",
			"CREATED",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	s.SetSelectedFunc(s.selectHandler)
	return s
}

func (s S3Buckets) GetLabels() []string {
	return []string{"Buckets"}
}

func (s S3Buckets) selectHandler(row, col int) {
	bucket, err := s.GetColSelection("NAME")
	if err != nil {
		return
	}
	objectsView := NewS3Objects(s.repo, bucket, s.app)
	s.app.AddAndSwitch(objectsView)
}

func (s S3Buckets) bucketPolicyHandler() {
	bucket, err := s.GetColSelection("NAME")
	if err != nil {
		return
	}
	policyView := NewS3BucketPolicy(s.repo, bucket, s.app)
	s.app.AddAndSwitch(policyView)
}

func (s S3Buckets) corsRulesHandler() {
	bucket, err := s.GetColSelection("NAME")
	if err != nil {
		return
	}
	corsRulesView := NewS3CORSRules(s.repo, bucket, s.app)
	s.app.AddAndSwitch(corsRulesView)
}

func (s S3Buckets) tagsHandler() {
	bucket, err := s.GetColSelection("NAME")
	if err != nil {
		return
	}
	tagsView := NewTags(s.repo, s.GetService(), "bucket:"+bucket, s.app)
	s.app.AddAndSwitch(tagsView)
}

func (s S3Buckets) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
			Description: "Bucket Policy",
			Action:      s.bucketPolicyHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'c', tcell.ModNone),
			Description: "CORS Rules",
			Action:      s.corsRulesHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
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
		var created string
		if v.CreationDate != nil {
			created = v.CreationDate.Format(utils.DefaultTimeFormat)
		}
		data = append(data, []string{
			utils.DerefString(v.Name, ""),
			created,
		})
	}
	s.SetData(data)
}
