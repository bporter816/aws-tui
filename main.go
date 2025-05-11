package main

import (

	// "os/exec"
	// "sort"

	"github.com/bporter816/aws-tui/cmd"
)

func main() {
	// playing around with region selection panel, TODO move to dialog
	/*
		regionCmd := exec.Command("aws", "configure", "get", "region")
		regionOutput, err := regionCmd.Output()
		if err != nil {
			panic(err)
		}
		region := strings.TrimSpace(string(regionOutput))
		fmt.Printf("region: %v\n", region)

		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			panic(err)
		}

		ec2Client := ec2.NewFromConfig(cfg)
		regions, err := ec2Client.DescribeRegions(context.TODO(), &ec2.DescribeRegionsInput{})
		if err != nil {
			panic(err)
		}

		var regionsArr []string
		for _, r := range regions.Regions {
			if r.RegionName != nil {
				regionsArr = append(regionsArr, *r.RegionName)
			}
		}
		sort.Strings(regionsArr)

		l := tview.NewList()
		l.ShowSecondaryText(false)
		l.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Rune() == 'k' {
				return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
			} else if event.Rune() == 'j' {
				return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
			}
			return event
		})

		for i, v := range regionsArr {
			l.AddItem(v, "", 0, nil)
			if v == region {
				l.SetCurrentItem(i)
			}
		}
		l.SetOffset(0, 0)
	*/

	cmd.Execute()
}
