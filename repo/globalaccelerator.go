package repo

import (
	"context"
	ga "github.com/aws/aws-sdk-go-v2/service/globalaccelerator"
	"github.com/bporter816/aws-tui/model"
)

type GlobalAccelerator struct {
	gaClient *ga.Client
}

func NewGlobalAccelerator(gaClient *ga.Client) *GlobalAccelerator {
	return &GlobalAccelerator{
		gaClient: gaClient,
	}
}

func (g GlobalAccelerator) ListAccelerators() ([]model.GlobalAcceleratorAccelerator, error) {
	pg := ga.NewListAcceleratorsPaginator(
		g.gaClient,
		&ga.ListAcceleratorsInput{},
	)
	var accelerators []model.GlobalAcceleratorAccelerator
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.GlobalAcceleratorAccelerator{}, err
		}
		for _, v := range out.Accelerators {
			accelerators = append(accelerators, model.GlobalAcceleratorAccelerator(v))
		}
	}
	return accelerators, nil
}
