package internal

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"strings"
)

type Taggable interface {
	ListTags(string) (model.Tags, error)
}

type Tags struct {
	*ui.Table
	repo       Taggable
	service    string
	resourceId string
	app        *Application
}

func NewTags(repo Taggable, service string, resourceId string, app *Application) *Tags {
	t := &Tags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo:       repo,
		service:    service,
		resourceId: resourceId,
		app:        app,
	}
	return t
}

func (t Tags) GetService() string {
	return t.service
}

func (t Tags) GetLabels() []string {
	var name string
	if arn, err := arn.Parse(t.resourceId); err == nil {
		// first try to interpret the resource id as an arn
		name = utils.GetResourceNameFromArn(arn)
	} else if strings.HasPrefix(t.resourceId, "https") {
		// this case is for SQS queues, which are defined by a URL
		parts := strings.Split(t.resourceId, "/")
		name = parts[len(parts)-1]
	} else {
		// use the id as is, removing prefixes added for IAM, Route 53, and S3
		parts := strings.Split(t.resourceId, ":")
		name = parts[len(parts)-1]
	}
	return []string{name, "Tags"}
}

func (t Tags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (t Tags) Render() {
	model, err := t.repo.ListTags(t.resourceId)
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
	t.SetData(data)
}
