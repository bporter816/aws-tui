package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
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

func (g GlobalAccelerator) ListListeners(acceleratorArn string) ([]model.GlobalAcceleratorListener, error) {
	pg := ga.NewListListenersPaginator(
		g.gaClient,
		&ga.ListListenersInput{
			AcceleratorArn: aws.String(acceleratorArn),
		},
	)
	var listeners []model.GlobalAcceleratorListener
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.GlobalAcceleratorListener{}, err
		}
		for _, v := range out.Listeners {
			listeners = append(listeners, model.GlobalAcceleratorListener(v))
		}
	}
	return listeners, nil
}

func (g GlobalAccelerator) ListTags(resourceId string) (model.Tags, error) {
	out, err := g.gaClient.ListTagsForResource(
		context.TODO(),
		&ga.ListTagsForResourceInput{
			ResourceArn: aws.String(resourceId),
		},
	)
	if err != nil {
		return model.Tags{}, err
	}
	var tags model.Tags
	for _, v := range out.Tags {
		tags = append(tags, model.Tag{Key: *v.Key, Value: *v.Value})
	}
	return tags, nil
}
