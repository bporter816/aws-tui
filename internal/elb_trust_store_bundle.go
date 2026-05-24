package internal

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
)

type ELBTrustStoreBundle struct {
	*ui.Text
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

func (e ELBTrustStoreBundle) GetService() string { return "ELB" }

func (e ELBTrustStoreBundle) GetLabels() []string {
	label := e.trustStoreArn
	if a, err := arn.Parse(e.trustStoreArn); err == nil {
		label = utils.GetResourceNameFromArn(a)
	}
	return []string{label, "Certificate Bundle"}
}

func (e ELBTrustStoreBundle) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ELBTrustStoreBundle) Render() error {
	cert, err := e.repo.GetTrustStoreCACertificatesBundle(e.trustStoreArn)
	if err != nil {
		return err
	}

	// TODO show the details like in ACM
	e.SetText(cert)
	return nil
}
