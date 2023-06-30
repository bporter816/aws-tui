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
	*tview.TextView
	stsClient *sts.Client
	iamClient *iam.Client
}

func NewHeader(s *sts.Client, i *iam.Client) *Header {
	h := &Header{
		TextView:  tview.NewTextView().SetDynamicColors(true),
		stsClient: s,
		iamClient: i,
	}
	h.Render() // TODO fix
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

	str := "[orange]Account:[white] " + *identityOutput.Account + aliases + "\n"
	str += "[orange]ARN:[white]     " + *identityOutput.Arn + "\n"
	str += "[orange]User ID:[white] " + *identityOutput.UserId + "\n"
	str += "[orange]Region:[white]  " + string(regionOutput)
	h.SetText(str)
}
