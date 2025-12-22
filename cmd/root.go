package cmd

import (
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use: "topk",
}

func init() {
	rootCommand.AddCommand(
		startRunner(),
	)
}

func Execute() error {
	return rootCommand.Execute()
}
