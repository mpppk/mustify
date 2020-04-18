package cmd

import (
	"bytes"
	"fmt"
	"go/format"
	"go/token"
	"io"
	"os"

	"golang.org/x/tools/imports"

	"github.com/mpppk/mustify/lib"

	"github.com/pkg/errors"

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
		Args:          cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewRootCmdConfigFromViper()
			if err != nil {
				return err
			}
			util.InitializeLog(conf.Verbose)

			filePath := "<standard input>"
			src := cmd.InOrStdin()
			if len(args) > 0 {
				filePath = args[0]
				src = nil
			}
			fset := token.NewFileSet()
			newFile, ok, err := lib.GenerateErrorWrappersFromReaderOrFile(fset, filePath, src)
			if err != nil {
				return err
			}
			if !ok {
				return nil
			}

			buf := new(bytes.Buffer)
			if err := format.Node(buf, fset, newFile); err != nil {
				return errors.Wrap(err, "failed to output")
			}
			newSrc, err := formatSrc(buf.Bytes())
			if err != nil {
				return err
			}

			if _, err := io.WriteString(cmd.OutOrStdout(), string(newSrc)); err != nil {
				return err
			}
			return nil
		},
	}

	if err := registerFlags(cmd); err != nil {
		return nil, err
	}

	return cmd, nil
}

func formatSrc(bytes []byte) ([]byte, error) {
	options := &imports.Options{
		TabWidth:  8,
		TabIndent: true,
		Comments:  true,
		Fragment:  true,
	}
	return imports.Process("<standard input>", bytes, options)
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
	rootCmd.SetOut(os.Stdout)
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
