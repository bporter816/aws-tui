package internal

import (
	"context"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
)

func NewS3ObjectMetadata(repo *repo.S3, bucket string, key string, app *Application) *TableView[model.Tag] {
	return NewTableView(
		"S3",
		[]string{key, "Metadata"},
		[]string{"KEY", "VALUE"},
		func(ctx context.Context) ([]model.Tag, error) {
			tags, err := repo.GetObjectMetadata(ctx, bucket, key)
			if err != nil {
				return nil, err
			}
			return []model.Tag(tags), nil
		},
		func(v model.Tag) []string {
			return []string{
				v.Key,
				v.Value,
			}
		},
		nil,
	)
}
