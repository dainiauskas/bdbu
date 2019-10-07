package cmd

import (
  "github.com/spf13/cobra"
  "bdbu/app"
)

var (
  tableName string
)

func init() {
  copyCmd.Flags().StringVarP(&tableName, "table", "t", "", "table name to copy")

	rootCmd.AddCommand(copyCmd)
}

var copyCmd = &cobra.Command{
  Use: "copy",
  Short: "copy database to another",
  Run: func(cmd *cobra.Command, args []string) {
    app.Copy(tableName)
  },
}
