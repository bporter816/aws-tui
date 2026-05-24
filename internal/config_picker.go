package internal

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rivo/tview"
)

var awsRegions = []string{
	"af-south-1",
	"ap-east-1",
	"ap-northeast-1",
	"ap-northeast-2",
	"ap-northeast-3",
	"ap-south-1",
	"ap-south-2",
	"ap-southeast-1",
	"ap-southeast-2",
	"ap-southeast-3",
	"ap-southeast-4",
	"ap-southeast-5",
	"ca-central-1",
	"ca-west-1",
	"eu-central-1",
	"eu-central-2",
	"eu-north-1",
	"eu-south-1",
	"eu-south-2",
	"eu-west-1",
	"eu-west-2",
	"eu-west-3",
	"il-central-1",
	"me-central-1",
	"me-south-1",
	"sa-east-1",
	"us-east-1",
	"us-east-2",
	"us-gov-east-1",
	"us-gov-west-1",
	"us-west-1",
	"us-west-2",
}

func readAWSProfiles() []string {
	seen := map[string]bool{}
	var profiles []string

	add := func(name string) {
		if !seen[name] {
			seen[name] = true
			profiles = append(profiles, name)
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return profiles
	}

	parseAWSFile(filepath.Join(home, ".aws", "config"), func(section string) {
		if section == "default" {
			add("default")
		} else if strings.HasPrefix(section, "profile ") {
			add(strings.TrimPrefix(section, "profile "))
		}
	})

	parseAWSFile(filepath.Join(home, ".aws", "credentials"), func(section string) {
		add(section)
	})

	return profiles
}

func parseAWSFile(path string, onSection func(string)) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			onSection(line[1 : len(line)-1])
		}
	}
}

func (a *Application) showPicker(title string, items []string, selected string, onSelect func(string)) {
	list := tview.NewList()
	list.ShowSecondaryText(false)
	list.SetBorder(true).SetTitle(fmt.Sprintf(" %v ", title))
	for _, item := range items {
		item := item
		list.AddItem(item, "", 0, func() {
			a.pages.RemovePage("picker")
			a.app.SetFocus(a.pages)
			onSelect(item)
		})
	}
	for i, item := range items {
		if item == selected {
			list.SetCurrentItem(i)
			break
		}
	}

	height := min(len(items)+2, 20)
	centered := tview.NewGrid().
		SetColumns(0, 60, 0).
		SetRows(0, height, 0).
		AddItem(list, 1, 1, 1, 1, 0, 0, true)

	a.pages.AddPage("picker", centered, true, true)
	a.app.SetFocus(list)
}

func (a *Application) showRegionPicker() {
	a.showPicker("Select Region", awsRegions, a.region, func(region string) {
		a.rebuild(a.profile, region)
	})
}

func (a *Application) showProfilePicker() {
	profiles := readAWSProfiles()
	if len(profiles) == 0 {
		a.ShowError(fmt.Errorf("no AWS profiles found in ~/.aws/config or ~/.aws/credentials"))
		return
	}
	a.showPicker("Select Profile", profiles, a.profile, func(profile string) {
		a.rebuild(profile, "")
	})
}
