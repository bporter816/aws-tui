package internal

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	rdsTypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/view"
)

type RDSEndpoints struct {
	*ui.Table
	view.RDS
	repo      *repo.RDS
	app       *Application
	clusterId string
}

func NewRDSEndpoints(repo *repo.RDS, app *Application, clusterId string) *RDSEndpoints {
	r := &RDSEndpoints{
		Table: ui.NewTable([]string{
			"ENDPOINT",
			"TYPE",
		}, 1, 0),
		repo:      repo,
		app:       app,
		clusterId: clusterId,
	}
	return r
}

func (r RDSEndpoints) GetLabels() []string {
	return []string{r.clusterId, "Endpoints"}
}

func (r RDSEndpoints) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (r RDSEndpoints) Render() {
	model, err := r.repo.ListClusters([]rdsTypes.Filter{
		{
			Name:   aws.String("db-cluster-id"),
			Values: []string{r.clusterId},
		},
	})
	if err != nil {
		panic(err)
	}
	if len(model) != 1 {
		panic(errors.New("expected 1 db cluster"))
	}
	cluster := model[0]
	var data [][]string
	if cluster.Endpoint != nil {
		data = append(data, []string{
			*cluster.Endpoint,
			"Writer",
		})
	}
	if cluster.ReaderEndpoint != nil {
		data = append(data, []string{
			*cluster.ReaderEndpoint,
			"Reader",
		})
	}
	for _, v := range cluster.CustomEndpoints {
		data = append(data, []string{
			v,
			"Custom",
		})
	}
	r.SetData(data)
}
