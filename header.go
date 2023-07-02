package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/rivo/tview"
	"os/exec"
	"strings"
)

type Header struct {
	*tview.Flex
	stsClient                *sts.Client
	iamClient                *iam.Client
	accountInfo, keybindInfo *tview.TextView
	app                      *Application
}

func NewHeader(s *sts.Client, i *iam.Client, app *Application) *Header {
	accountInfo := tview.NewTextView().SetDynamicColors(true)
	keybindInfo := tview.NewTextView().SetDynamicColors(true) // TODO make this a grid

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexColumn)
	flex.AddItem(accountInfo, 0, 1, true) // TODO make this fixed size
	flex.AddItem(keybindInfo, 0, 1, true)

	h := &Header{
		Flex:        flex,
		stsClient:   s,
		iamClient:   i,
		accountInfo: accountInfo,
		keybindInfo: keybindInfo,
		app:         app,
	}
	return h
}

func (h Header) Render() {
	// The AWS Go SDK doesn't provide a nice way to get the current region so get the answer from the AWS CLI
	regionCmd := exec.Command("aws", "configure", "get", "region")
	regionOutput, err := regionCmd.Output()
	if err != nil {
		panic(err)
	}

	identityOutput, err := h.stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		panic(err)
	}

	aliasesOutput, err := h.iamClient.ListAccountAliases(context.TODO(), &iam.ListAccountAliasesInput{})
	if err != nil {
		panic(err)
	}

	var aliases string
	if len(aliasesOutput.AccountAliases) > 0 {
		aliases = fmt.Sprintf(" (%v)", strings.Join(aliasesOutput.AccountAliases, ", "))
	}

	accountInfoStr := "[orange::b]Account:[white::-] " + *identityOutput.Account + aliases + "\n"
	accountInfoStr += "[orange::b]ARN:[white::-]     " + *identityOutput.Arn + "\n"
	accountInfoStr += "[orange::b]User ID:[white::-] " + *identityOutput.UserId + "\n"
	accountInfoStr += "[orange::b]Region:[white::-]  " + string(regionOutput)
	h.accountInfo.SetText(accountInfoStr)

	keybindInfoStr := ""
	actions := h.app.GetActiveKeyActions()
	for _, v := range actions {
		name := v.Key.Name()
		if strings.HasPrefix(name, "Rune[") {
			name = string(name[5])
		}
		keybindInfoStr += fmt.Sprintf("[pink::b]<%v>[white::-] %v", name, v.Description)
	}
	h.keybindInfo.SetText(keybindInfoStr)
}
