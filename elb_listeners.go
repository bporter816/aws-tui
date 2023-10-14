package main

import (
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/gdamore/tcell/v2"
	"strconv"
)

// TODO support load balancer lists by lb arns
// TODO better manage name/arn
type ELBListeners struct {
	*ui.Table
	repo   *repo.ELB
	lbArn  string
	lbName string
	app    *Application
	model  []model.ELBListener
}

func NewELBListeners(repo *repo.ELB, lbArn string, lbName string, app *Application) *ELBListeners {
	e := &ELBListeners{
		Table: ui.NewTable([]string{
			"PROTOCOL",
			"PORT",
			"RULES",
			"SSL POLICY",
			"DEFAULT CERTIFICATE",
		}, 1, 0),
		repo:   repo,
		lbArn:  lbArn,
		lbName: lbName,
		app:    app,
	}
	return e
}

func (e ELBListeners) GetService() string {
	return "ELB"
}

func (e ELBListeners) GetLabels() []string {
	return []string{e.lbName, "Listeners"}
}

func (e ELBListeners) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	protocol, err := e.GetColSelection("PROTOCOL")
	if err != nil {
		return
	}
	port, err := e.GetColSelection("PORT")
	if err != nil {
		return
	}
	if arn := e.model[row-1].ListenerArn; arn != nil {
		tagsView := NewELBTags(e.repo, ELBResourceTypeListener, *arn, protocol+":"+port, e.app)
		e.app.AddAndSwitch(tagsView)
	}
}

func (e ELBListeners) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e *ELBListeners) Render() {
	model, err := e.repo.ListListeners(e.lbArn)
	if err != nil {
		panic(err)
	}
	e.model = model

	var data [][]string
	for _, v := range model {
		var protocol, port, rules, sslPolicy, defaultCertificate string
		protocol = string(v.Protocol)
		if v.Port != nil {
			port = strconv.Itoa(int(*v.Port))
		}
		rules = strconv.Itoa(v.Rules)
		if v.SslPolicy != nil {
			sslPolicy = *v.SslPolicy
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
		data = append(data, []string{
			protocol,
			port,
			rules,
			sslPolicy,
			defaultCertificate,
		})
	}
	e.SetData(data)
}
