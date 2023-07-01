package main

import (
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Services struct {
	*Table
	clients map[string]interface{}
	app     *Application
}

func NewServices(clients map[string]interface{}, app *Application) *Services {
	s := &Services{
		Table: NewTable([]string{
			"SERVICE",
			"DESCRIPTION",
		}, 1, 0),
		clients: clients,
		app:     app,
	}
	s.Render() // TODO fix
	return s
}

func (s Services) GetName() string {
	return "Services"
}

func (s Services) selectHandler() {
	r, _ := s.GetSelection()
	service := s.GetCell(r, 0).Text
	var item tview.Primitive
	switch service {
	case "KMS":
		item = NewKmsKeys(s.clients["KMS"].(*kms.Client), s.app)
	case "Route 53":
		item = NewRoute53HostedZones(s.clients["Route 53"].(*route53.Client), s.app)
	case "Secrets Manager":
		item = NewSecretsManagerSecrets(s.clients["Secrets Manager"].(*sm.Client), s.app)
	default:
		panic("unknown service")
	}
	s.app.AddAndSwitch(service, item)
}

func (s Services) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone),
			Description: "",
			Action:      s.selectHandler,
		},
	}
}

func (s Services) Render() {
	data := [][]string{
		[]string{"KMS", "Key Management Service"},
		[]string{"Route 53", "DNS"},
		[]string{"Secrets Manager", "Secrets"},
	}
	s.SetData(data)
}
