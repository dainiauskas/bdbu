package cmd

import (
  "github.com/spf13/cobra"
  "bdbu/app"
)

func init() {
	rootCmd.AddCommand(tunerCmd)
}

var tunerCmd = &cobra.Command{
  Use: "tuner",
  Short: "MariaDB Tuner",
  Run: func(cmd *cobra.Command, args []string) {
    app.Tuner()
  },
}
