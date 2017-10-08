package bowielib

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	userName    string
	projectName string
)

// RootCmd is the main command executed when bowie is run
var RootCmd = &cobra.Command{
	Use:   "bowie",
	Short: "Bowie is a pretty changelog generator",
	Long: `A pretty changelog generator,
built with love by mattstratton in Go.
	
Complete documentation is available at https://github.com/mattstratton/bowie`,

	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bowie.yaml)")
	RootCmd.PersistentFlags().StringVarP(&userName, "username", "u", "", "github user/org name, i.e., mattstratton")
	RootCmd.PersistentFlags().StringVarP(&projectName, "projectname", "p", "", "github project name, i.e., bowie")
	viper.BindPFlag("username", RootCmd.PersistentFlags().Lookup("username"))       // nolint: errcheck
	viper.BindPFlag("projectname", RootCmd.PersistentFlags().Lookup("projectname")) // nolint: errcheck
}

// Execute runs the main bowie command
func Execute() {
	RootCmd.Execute() // nolint: errcheck
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".bowie" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".bowie")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
