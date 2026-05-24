package internal

import (
	"strconv"
	"strings"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
)

func NewS3CORSRules(repo *repo.S3, bucket string, app *Application) *TableView[model.S3CORSRule] {
	return NewTableView(
		"S3",
		[]string{bucket, "CORS"},
		[]string{"ID", "ALLOWED METHODS", "ALLOWED ORIGINS", "ALLOWED HEADERS", "EXPOSE HEADERS", "MAX AGE"},
		func() ([]model.S3CORSRule, error) {
			return repo.GetCORSRules(bucket)
		},
		func(v model.S3CORSRule) []string {
			return []string{
				utils.DerefString(v.ID, ""),
				strings.Join(v.AllowedMethods, ", "),
				strings.Join(v.AllowedOrigins, ", "),
				strings.Join(v.AllowedHeaders, ", "),
				strings.Join(v.ExposeHeaders, ", "),
				strconv.Itoa(int(*v.MaxAgeSeconds)) + " sec",
			}
		},
		nil,
	)
}
