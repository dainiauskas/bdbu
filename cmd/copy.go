package cmd

import (
  "github.com/spf13/cobra"
  "bdbu/app"
)

func init() {
	rootCmd.AddCommand(copyCmd)
}

var copyCmd = &cobra.Command{
  Use: "copy",
  Short: "copy database to another",
  Run: func(cmd *cobra.Command, args []string) {
    app.Copy()
    // if err := app.OpenDb(Verbose); err != nil {
    //   panic(err)
    // }
    // app.Start()
  },
}
