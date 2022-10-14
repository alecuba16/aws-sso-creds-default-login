package list

import (
	"github.com/alecuba16/aws-sso-creds-default-login/cmd/aws-sso-creds/list/accounts"
	"github.com/alecuba16/aws-sso-creds-default-login/cmd/aws-sso-creds/list/roles"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "list commands",
		Long:  "Commands that list things",
	}

	command.AddCommand(accounts.Command())
	command.AddCommand(roles.Command())

	return command
}
