package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewS3Buckets(repo *repo.S3, app *Application) *TableView[model.S3Bucket] {
	v := NewTableView(
		"S3",
		[]string{"Buckets"},
		[]string{"NAME", "CREATED"},
		repo.ListBuckets,
		func(b model.S3Bucket) []string {
			var created string
			if b.CreationDate != nil {
				created = b.CreationDate.Format(utils.DefaultTimeFormat)
			}
			return []string{
				utils.DerefString(b.Name, ""),
				created,
			}
		},
		func(v *TableView[model.S3Bucket]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
					Description: "Bucket Policy",
					Action: func() {
						bucket, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewS3BucketPolicy(repo, bucket, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'c', tcell.ModNone),
					Description: "CORS Rules",
					Action: func() {
						bucket, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewS3CORSRules(repo, bucket, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						bucket, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), "bucket:"+bucket, app))
					},
				},
			}
		},
	)
	v.SetSelectedFunc(func(row, col int) {
		bucket, err := v.GetColSelection("NAME")
		if err != nil {
			return
		}
		app.AddAndSwitch(NewS3Objects(repo, bucket, app))
	})
	return v
}

func init() {
	register("S3.Buckets", func(r Repos, app *Application) Component {
		return NewS3Buckets(r.S3, app)
	})
}
