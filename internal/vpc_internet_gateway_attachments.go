package internal

import (
	"context"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
)

func NewVPCInternetGatewayAttachments(repo *repo.EC2, internetGatewayId string, app *Application) *TableView[model.EC2InternetGatewayAttachment] {
	return NewTableView(
		"VPC",
		[]string{internetGatewayId, "Attachments"},
		[]string{"VPC ID", "STATE"},
		func(ctx context.Context) ([]model.EC2InternetGatewayAttachment, error) {
			return repo.ListInternetGatewayAttachments(ctx, internetGatewayId)
		},
		func(v model.EC2InternetGatewayAttachment) []string {
			return []string{
				utils.DerefString(v.VpcId, ""),
				utils.AutoCase(string(v.State)),
			}
		},
		nil,
	)
}
