package cmd

import (
	"bdbu/app"

	"github.com/spf13/cobra"
)

var (
	tableName  string
	dropTables bool
)

func init() {
	copyCmd.Flags().StringVarP(&tableName, "table", "t", "", "table name to copy")
	copyCmd.Flags().BoolVarP(&dropTables, "with-drop", "D", false, "drop table if exists")

	rootCmd.AddCommand(copyCmd)
}

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "copy database to another",
	Run: func(cmd *cobra.Command, args []string) {
		app.Copy(tableName, dropTables)
	},
}
