package cmd

import (
  "github.com/spf13/cobra"
  "bdbu/app"
)

func init() {
	rootCmd.AddCommand(mysqlCmd)
}

var mysqlCmd = &cobra.Command{
  Use: "mysql",
  Short: "MySQL Option number explain",
  Run: func(cmd *cobra.Command, args []string) {
    app.MyOptionExplain()
  },
}
