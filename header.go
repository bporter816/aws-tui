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
	stsClient   *sts.Client
	iamClient   *iam.Client
	accountInfo *tview.TextView
	keybindInfo *tview.Grid
	app         *Application
}

func NewHeader(s *sts.Client, i *iam.Client, app *Application) *Header {
	accountInfo := tview.NewTextView()
	accountInfo.SetDynamicColors(true)
	accountInfo.SetWrap(false)

	keybindInfo := tview.NewGrid()
	keybindInfo.SetRows(1, 1, 1, 1) // header is 4 rows
	keybindInfo.SetColumns(0)       // start with one column, but it will resize itself if it overflows

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

	h.keybindInfo.Clear()
	actions := h.app.GetActiveKeyActions()
	row, col := 0, 0
	for _, v := range actions {
		name := v.Key.Name()
		if strings.HasPrefix(name, "Rune[") {
			name = string(name[5])
		}
		entry := tview.NewTextView().SetDynamicColors(true).SetText(fmt.Sprintf("[pink::b]<%v>[white::-] %v", name, v.Description))
		h.keybindInfo.AddItem(entry, row, col, 1, 1, 1, 1, false)

		row++
		if row == 4 {
			row = 0
			col++
		}
	}
	// cleanup empty rows so they're the same color
	// TODO see if this can be avoided
	for row < 4 {
		h.keybindInfo.AddItem(tview.NewTextView(), row, col, 1, 1, 1, 1, false)
		row++
	}
}
