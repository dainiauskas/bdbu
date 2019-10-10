package cmd

import (
  "github.com/spf13/cobra"
  "bdbu/app"
)

var (
  records     int
  bEngine     string
  bRowFormat  string
)

func init() {
  benchCmd.Flags().IntVarP(&records, "records", "r", 100000, "number of records to insert")
  benchCmd.Flags().StringVarP(&bEngine, "engine", "e", "InnoDB", "engine for create table")
  benchCmd.Flags().StringVarP(&bRowFormat, "rowformat", "f", "DYNAMIC", "row format for create table")

	rootCmd.AddCommand(benchCmd)
}

var benchCmd = &cobra.Command{
  Use: "benchmark",
  Short: "database benchmark",
  Run: func(cmd *cobra.Command, args []string) {
    app.Benchmark(records, bEngine, bRowFormat)
  },
}
