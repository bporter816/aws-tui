package cmd

import (
	"fmt"
	"os"

	"github.com/bporter816/aws-tui/internal"
	"github.com/bporter816/aws-tui/internal/template"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "aws-tui",
	Run: func(cmd *cobra.Command, args []string) {
		template.Init()

		app := internal.NewApplication()
		if err := app.Run(); err != nil {
			panic(err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
