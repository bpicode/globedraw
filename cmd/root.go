package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any sub-commands.
var RootCmd = &cobra.Command{
	Use: "globedraw",
}

func init() {
	cobra.OnInitialize()
}
