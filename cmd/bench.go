package cmd

import (
	"bdbu/app"

	"github.com/spf13/cobra"
)

var (
	records    int
	bEngine    string
	bRowFormat string
	diskWrite  bool
	diskAll    bool
	diskSize   int
)

func init() {
	benchDisk.Flags().BoolVarP(&diskWrite, "write", "w", false, "write to disk test")
	benchDisk.Flags().BoolVarP(&diskAll, "all", "a", true, "all tests")
	benchDisk.Flags().IntVarP(&diskSize, "size", "s", 8, "file size in GiB")
	benchCmd.AddCommand(benchDisk)

	benchCmd.Flags().IntVarP(&records, "records", "r", 100000, "number of records to insert")
	benchCmd.Flags().StringVarP(&bEngine, "engine", "e", "InnoDB", "engine for create table")
	benchCmd.Flags().StringVarP(&bRowFormat, "rowformat", "f", "DYNAMIC", "row format for create table")

	rootCmd.AddCommand(benchCmd)
}

var benchCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "database benchmark",
	Run: func(cmd *cobra.Command, args []string) {
		app.Benchmark(records, bEngine, bRowFormat)
	},
}

var benchDisk = &cobra.Command{
	Use:   "disk",
	Short: "disk benchmark",
	Run: func(cmd *cobra.Command, args []string) {
		app.BenchDiskWrite()
	},
}
