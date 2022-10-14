package export

import (
	"aws-sso-creds-default/pkg/credentials"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:          "export",
		Short:        "Generates a set of shell commands to export AWS temporary creds to your environment",
		Long:         "Generates a set of shell commands to export AWS temporary creds to your environment",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			cmd.SilenceUsage = true

			profile := viper.GetString("profile")
			homeDir := viper.GetString("home-directory")
			login := viper.GetBool("login")

			if profile == "" {
				return fmt.Errorf("no profile specified")
			}

			creds, _, err := credentials.GetSSOCredentials(profile, homeDir, login)

			if err != nil {
				return err
			}

			fmt.Printf("export AWS_ACCESS_KEY_ID=%s\n", *creds.RoleCredentials.AccessKeyId)
			fmt.Printf("export AWS_SECRET_ACCESS_KEY=%s\n", *creds.RoleCredentials.SecretAccessKey)
			fmt.Printf("export AWS_SESSION_TOKEN=%s\n", *creds.RoleCredentials.SessionToken)

			return nil
		},
	}

	return command
}
