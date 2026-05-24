package internal

import (
	"context"
	"strconv"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

// TODO support load balancer lists by lb arns
// TODO better manage name/arn
func NewELBListeners(repo *repo.ELB, lbArn string, lbName string, app *Application) *TableView[model.ELBListener] {
	return NewTableView(
		"ELB",
		[]string{lbName, "Listeners"},
		[]string{"PROTOCOL", "PORT", "RULES", "SSL POLICY", "DEFAULT CERTIFICATE", "MTLS MODE"},
		func(ctx context.Context) ([]model.ELBListener, error) {
			return repo.ListListeners(ctx, lbArn)
		},
		func(v model.ELBListener) []string {
			var port, defaultCertificate, mtlsMode string
			if v.Port != nil {
				port = strconv.Itoa(int(*v.Port))
			}
			for _, c := range v.Certificates {
				if c.IsDefault != nil && *c.IsDefault && c.CertificateArn != nil {
					defaultCertificate = *c.CertificateArn
					break
				}
			}
			if len(v.Certificates) == 1 && v.Certificates[0].CertificateArn != nil {
				defaultCertificate = *v.Certificates[0].CertificateArn
			}
			if v.MutualAuthentication != nil && v.MutualAuthentication.Mode != nil {
				mtlsMode = utils.AutoCase(*v.MutualAuthentication.Mode)
			}
			return []string{
				string(v.Protocol),
				port,
				strconv.Itoa(v.Rules),
				utils.DerefString(v.SslPolicy, ""),
				defaultCertificate,
				mtlsMode,
			}
		},
		func(v *TableView[model.ELBListener]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.ListenerArn == nil {
							return
						}
						// TODO display protocol and port
						app.AddAndSwitch(NewTags(repo, v.GetService(), *item.ListenerArn, app))
					},
				},
			}
		},
	)
}
