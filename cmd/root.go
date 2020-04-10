package cmd

import (
	"fmt"
	"go/format"
	"go/token"
	"os"

	"github.com/pkg/errors"

	"github.com/mpppk/mustify/lib"

	"github.com/mpppk/mustify/util"

	"github.com/mpppk/mustify/cmd/option"

	"github.com/spf13/afero"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// NewRootCmd generate root cmd
func NewRootCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:           "mustify",
		Short:         "generate MustXXX methods from go source",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewRootCmdConfigFromViper()
			if err != nil {
				return err
			}
			util.InitializeLog(conf.Verbose)

			filePath := args[0]

			fileMap, err := lib.GenerateErrorWrappersFromPackage(filePath, "main", "must-")
			if err != nil {
				panic(err)
			}

			for _, file := range fileMap {
				if err := format.Node(os.Stdout, token.NewFileSet(), file); err != nil {
					return errors.Wrap(err, "failed to write ast file to stdout")
				}
			}

			return nil
		},
	}

	if err := registerFlags(cmd); err != nil {
		return nil, err
	}

	return cmd, nil
}

func registerFlags(cmd *cobra.Command) error {
	flags := []option.Flag{}
	return option.RegisterFlags(cmd, flags)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd, err := NewRootCmd(afero.NewOsFs())
	if err != nil {
		panic(err)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Print(util.PrettyPrintError(err))
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
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

		// Search config in home directory with name ".mustify" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".mustify")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
