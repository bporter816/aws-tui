package internal

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	rdsTypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/bporter816/aws-tui/internal/repo"
)

type rdsEndpoint struct {
	address      string
	endpointType string
}

func NewRDSEndpoints(repo *repo.RDS, app *Application, clusterId string) *TableView[rdsEndpoint] {
	return NewTableView(
		"RDS",
		[]string{clusterId, "Endpoints"},
		[]string{"ENDPOINT", "TYPE"},
		func() ([]rdsEndpoint, error) {
			model, err := repo.ListClusters([]rdsTypes.Filter{
				{
					Name:   aws.String("db-cluster-id"),
					Values: []string{clusterId},
				},
			})
			if err != nil {
				return nil, err
			}
			if len(model) != 1 {
				return nil, errors.New("expected 1 db cluster")
			}
			cluster := model[0]
			var endpoints []rdsEndpoint
			if cluster.Endpoint != nil {
				endpoints = append(endpoints, rdsEndpoint{address: *cluster.Endpoint, endpointType: "Writer"})
			}
			if cluster.ReaderEndpoint != nil {
				endpoints = append(endpoints, rdsEndpoint{address: *cluster.ReaderEndpoint, endpointType: "Reader"})
			}
			for _, v := range cluster.CustomEndpoints {
				endpoints = append(endpoints, rdsEndpoint{address: v, endpointType: "Custom"})
			}
			return endpoints, nil
		},
		func(v rdsEndpoint) []string {
			return []string{
				v.address,
				v.endpointType,
			}
		},
		nil,
	)
}
