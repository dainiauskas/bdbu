package cmd

import (
  "github.com/spf13/cobra"
  "bdbu/app"
)

func init() {
	rootCmd.AddCommand(benchCmd)
}

var benchCmd = &cobra.Command{
  Use: "benchmark",
  Short: "database benchmark",
  Run: func(cmd *cobra.Command, args []string) {
    app.Benchmark()
  },
}
