package internal

import (
	"fmt"
	"strings"

	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/rivo/tview"
)

type Header struct {
	*tview.Flex
	stsRepo     *repo.STS
	iamRepo     *repo.IAM
	app         *Application
	accountInfo *tview.TextView
	keybindInfo *tview.Grid
}

func NewHeader(stsRepo *repo.STS, iamRepo *repo.IAM, app *Application) *Header {
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
		stsRepo:     stsRepo,
		iamRepo:     iamRepo,
		accountInfo: accountInfo,
		keybindInfo: keybindInfo,
		app:         app,
	}
	return h
}

func (h *Header) UpdateRepos(stsRepo *repo.STS, iamRepo *repo.IAM) {
	h.stsRepo = stsRepo
	h.iamRepo = iamRepo
}

func (h Header) Render() {
	region := h.app.GetRegion()
	if region == "" {
		region = "unknown"
	}

	var account, arn, userId string
	if identityModel, err := h.stsRepo.GetCallerIdentity(); err == nil {
		if identityModel.Account != nil {
			account = *identityModel.Account
		}
		if identityModel.Arn != nil {
			arn = *identityModel.Arn
		}
		if identityModel.UserId != nil {
			userId = *identityModel.UserId
		}
	}

	aliases, _ := h.iamRepo.ListAccountAliases()
	var aliasesStr string
	if len(aliases) > 0 {
		aliasesStr = fmt.Sprintf(" (%v)", strings.Join(aliases, ", "))
	}

	accountInfoStr := fmt.Sprintf("[orange::b]Account:[white::-] %v%v\n", account, aliasesStr)
	accountInfoStr += fmt.Sprintf("[orange::b]ARN:[white::-]     %v\n", arn)
	accountInfoStr += fmt.Sprintf("[orange::b]User ID:[white::-] %v\n", userId)
	accountInfoStr += fmt.Sprintf("[orange::b]Region:[white::-]  %v", region)
	h.accountInfo.SetText(accountInfoStr)

	h.Box = tview.NewBox() // this is needed because the areas not covered by items are considered transparent and will linger otherwise
	h.keybindInfo.Clear()
	actions := h.app.GetActiveKeyActions()
	row, col := 0, 0
	for _, v := range actions {
		entry := tview.NewTextView().SetDynamicColors(true).SetText(fmt.Sprintf("[pink::b]<%v>[white::-] %v", v.String(), v.Description))
		h.keybindInfo.AddItem(entry, row, col, 1, 1, 1, 1, false)

		row++
		if row == 4 {
			row = 0
			col++
		}
	}
}
