package main

import (
	"strconv"

	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

// TODO support load balancer lists by lb arns
// TODO better manage name/arn
type ELBListeners struct {
	*ui.Table
	view.ELB
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
			"MTLS MODE",
		}, 1, 0),
		repo:   repo,
		lbArn:  lbArn,
		lbName: lbName,
		app:    app,
	}
	return e
}

func (e ELBListeners) GetLabels() []string {
	return []string{e.lbName, "Listeners"}
}

func (e ELBListeners) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	/*
		protocol, err := e.GetColSelection("PROTOCOL")
		if err != nil {
			return
		}
		port, err := e.GetColSelection("PORT")
		if err != nil {
			return
		}
	*/
	if arn := e.model[row-1].ListenerArn; arn != nil {
		// TODO display protocol and port
		tagsView := NewTags(e.repo, e.GetService(), *arn, e.app)
		e.app.AddAndSwitch(tagsView)
	}
}

func (e ELBListeners) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
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
		var protocol, port, rules, sslPolicy, defaultCertificate, mtlsMode string
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
		if v.MutualAuthentication != nil && v.MutualAuthentication.Mode != nil {
			mtlsMode = utils.AutoCase(*v.MutualAuthentication.Mode)
		}
		data = append(data, []string{
			protocol,
			port,
			rules,
			sslPolicy,
			defaultCertificate,
			mtlsMode,
		})
	}
	e.SetData(data)
}
