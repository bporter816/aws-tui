package internal

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
)

type ELBTrustStoreBundle struct {
	*ui.Text
	view.ELB
	repo          *repo.ELB
	trustStoreArn string
	app           *Application
}

func NewELBTrustStoreBundle(repo *repo.ELB, trustStoreArn string, app *Application) *ELBTrustStoreBundle {
	e := &ELBTrustStoreBundle{
		Text:          ui.NewText(false, "text"),
		repo:          repo,
		trustStoreArn: trustStoreArn,
		app:           app,
	}
	return e
}

func (e ELBTrustStoreBundle) GetLabels() []string {
	a, err := arn.Parse(e.trustStoreArn)
	if err != nil {
		panic(err)
	}
	return []string{utils.GetResourceNameFromArn(a), "Certificate Bundle"}
}

func (e ELBTrustStoreBundle) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ELBTrustStoreBundle) Render() {
	cert, err := e.repo.GetTrustStoreCACertificatesBundle(e.trustStoreArn)
	if err != nil {
		panic(err)
	}

	// TODO show the details like in ACM
	e.SetText(cert)
}
