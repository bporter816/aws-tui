package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/bporter816/aws-tui/model"
)

type ACM struct {
	acmClient *acm.Client
}

func NewACM(acmClient *acm.Client) *ACM {
	return &ACM{
		acmClient: acmClient,
	}
}

func (a ACM) ListCertificates() ([]model.ACMCertificate, error) {
	pg := acm.NewListCertificatesPaginator(
		a.acmClient,
		&acm.ListCertificatesInput{},
	)
	var certificates []model.ACMCertificate
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ACMCertificate{}, err
		}
		for _, v := range out.CertificateSummaryList {
			certificates = append(certificates, model.ACMCertificate(v))
		}
	}
	return certificates, nil
}
