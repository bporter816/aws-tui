package repo

import (
	"context"
	sq "github.com/aws/aws-sdk-go-v2/service/servicequotas"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/bporter816/aws-tui/model"
)

type ServiceQuotas struct {
	sqClient *sq.Client
}

func NewServiceQuotas(sqClient *sq.Client) *ServiceQuotas {
	return &ServiceQuotas{
		sqClient: sqClient,
	}
}

func (s ServiceQuotas) ListServices() ([]model.ServiceQuotasService, error) {
	pg := sq.NewListServicesPaginator(
		s.sqClient,
		&sq.ListServicesInput{},
	)
	var services []model.ServiceQuotasService
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ServiceQuotasService{}, err
		}
		for _, v := range out.Services {
			services = append(services, model.ServiceQuotasService(v))
		}
	}
	return services, nil
}

func (s ServiceQuotas) ListQuotas(serviceCode string) ([]model.ServiceQuotasQuota, error) {
	defaultsPg := sq.NewListAWSDefaultServiceQuotasPaginator(
		s.sqClient,
		&sq.ListAWSDefaultServiceQuotasInput{
			ServiceCode: aws.String(serviceCode),
		},
	)
	var quotas []model.ServiceQuotasQuota
	for defaultsPg.HasMorePages() {
		out, err := defaultsPg.NextPage(context.TODO())
		if err != nil {
			return []model.ServiceQuotasQuota{}, err
		}
		for _, v := range out.Quotas {
			m := model.ServiceQuotasQuota{
				ServiceQuota: v,
				DefaultValue: v.Value,
			}
			quotas = append(quotas, m)
		}
	}

	appliedPg := sq.NewListServiceQuotasPaginator(
		s.sqClient,
		&sq.ListServiceQuotasInput{
			ServiceCode: aws.String(serviceCode),
		},
	)
	for appliedPg.HasMorePages() {
		out, err := appliedPg.NextPage(context.TODO())
		if err != nil {
			return []model.ServiceQuotasQuota{}, err
		}
		for _, q := range out.Quotas {
			if q.QuotaName != nil {
				for i, v := range quotas {
					if v.QuotaName != nil && *q.QuotaName == *v.QuotaName {
						quotas[i].AppliedValue = q.Value
						break
					}
				}
			}
		}
	}
	return quotas, nil
}
