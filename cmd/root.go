package cmd

import (
  "strings"
  "os"
  "path/filepath"

  "github.com/spf13/cobra"
  "github.com/spf13/viper"
  "butent/tools/log"

  "bdbu/app"
)

var (
  Verbose     bool
  configFile  string
  AppDir      = workDir()
)

var rootCmd = &cobra.Command{
	Use: strings.ToLower(Name),
	Short: Name,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
    log.SetLogTrace(Verbose || app.Config.Verbose)
    log.SetLogToConsole(app.Config.Console)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Warn("%s", err)
    return
	}
}

func workDir() string {
  ex, err := os.Executable()
  if err != nil {
    panic(err)
  }

  return filepath.Dir(ex)
}

func init() {
  if _, err := os.Stat("bdbu.yaml"); os.IsNotExist(err) {
    if err := os.Chdir(AppDir); err != nil {
      panic(err)
    }
  }

  log.Init("./log", 400, 20, 100, true)

  log.SetFilenamePrefix("", "")
  log.SetLogThrough(false)
  log.SetLogFunctionName(false)
  log.SetLogFilenameLineNum(false)

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is bdbu.yaml)")
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("bdbu")
		viper.AddConfigPath(".")
		// viper.AddConfigPath(AppDir)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
    panic(err)
	}

  if err := app.GetConfig(); err != nil {
    panic(err)
  }  
}
