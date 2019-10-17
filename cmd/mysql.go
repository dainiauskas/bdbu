package cmd

import (
  "github.com/spf13/cobra"
  "bdbu/app"
)

var (
  mdMysqlShow bool
)

func init() {
  mysqlCmd.Flags().BoolVarP(&mdMysqlShow, "active", "a", false, "show only active")

	rootCmd.AddCommand(mysqlCmd)
}

var mysqlCmd = &cobra.Command{
  Use: "mysql",
  Short: "MySQL Option number explain",
  Run: func(cmd *cobra.Command, args []string) {
    app.MyOptionExplain(mdMysqlShow)
  },
}
