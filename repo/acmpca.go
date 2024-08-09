package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/acmpca"
	"github.com/bporter816/aws-tui/model"
)

type ACMPCA struct {
	acmPCAClient *acmpca.Client
}

func NewACMPCA(acmPCAClient *acmpca.Client) *ACMPCA {
	return &ACMPCA{
		acmPCAClient: acmPCAClient,
	}
}

func (a ACMPCA) ListCertificateAuthorities() ([]model.ACMPCACertificateAuthority, error) {
	pg := acmpca.NewListCertificateAuthoritiesPaginator(
		a.acmPCAClient,
		&acmpca.ListCertificateAuthoritiesInput{},
	)
	var certificateAuthorities []model.ACMPCACertificateAuthority
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ACMPCACertificateAuthority{}, err
		}
		for _, v := range out.CertificateAuthorities {
			certificateAuthorities = append(certificateAuthorities, model.ACMPCACertificateAuthority(v))
		}
	}
	return certificateAuthorities, nil
}
