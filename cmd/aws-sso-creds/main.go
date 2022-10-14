package main

import (
	"aws-sso-creds-default/cmd/aws-sso-creds/export"
	"aws-sso-creds-default/cmd/aws-sso-creds/get"
	"aws-sso-creds-default/cmd/aws-sso-creds/helper"
	"aws-sso-creds-default/cmd/aws-sso-creds/list"
	"aws-sso-creds-default/cmd/aws-sso-creds/set"
	"aws-sso-creds-default/cmd/aws-sso-creds/version"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	profile   string
	asdefault bool
	login     bool
)

func configureCLI() *cobra.Command {
	rootCommand := &cobra.Command{
		Use:  "aws-sso-creds",
		Long: "A helper utility to interact with AWS SSO",
	}

	rootCommand.AddCommand(get.Command())
	rootCommand.AddCommand(set.Command())
	rootCommand.AddCommand(version.Command())
	rootCommand.AddCommand(export.Command())
	rootCommand.AddCommand(list.Command())
	rootCommand.AddCommand(helper.Command())

	homeDir, err := homedir.Dir()

	if err != nil {
		panic("Cannot find home directory, fatal error")
	}

	rootCommand.PersistentFlags().StringVarP(&profile, "profile", "p", "", "the AWS profile to use")
	rootCommand.PersistentFlags().StringVarP(&homeDir, "home-directory", "H", homeDir, "specify a path to a home directory")
	rootCommand.PersistentFlags().BoolVarP(&asdefault, "asdefault", "d", false, "stores as the default profile in the credentials.")
	rootCommand.PersistentFlags().BoolVarP(&login, "login", "l", false, "login to AWS SSO before retrieving credentials.")

	viper.BindEnv("profile", "AWS_PROFILE")
	viper.BindPFlag("profile", rootCommand.PersistentFlags().Lookup("profile"))
	viper.BindPFlag("home-directory", rootCommand.PersistentFlags().Lookup("home-directory"))
	viper.BindPFlag("asdefault", rootCommand.PersistentFlags().Lookup("asdefault"))
	viper.BindPFlag("login", rootCommand.PersistentFlags().Lookup("login"))

	return rootCommand
}

func main() {
	rootCommand := configureCLI()

	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
