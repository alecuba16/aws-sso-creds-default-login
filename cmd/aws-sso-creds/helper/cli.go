package helper

import (
	"aws-sso-creds-default/pkg/credentials"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

type CredentialsProcessOutput struct {
	Version         int    `json:"page"`
	AccessKeyId     string `json:"AccessKeyId"`
	SecretAccessKey string `json:"SecretAccessKey"`
	SessionToken    string `json:"SessionToken"`
	Expiration      string `json:"Expiration"`
}

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:          "helper",
		Short:        "Generate credential helper compatible credentials",
		Long:         "Retrieve AWS temporary credentials and output in a format suitable for the credential_process field in an AWS profile",
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

			rawCreds := CredentialsProcessOutput{
				Version:         1,
				AccessKeyId:     *creds.RoleCredentials.AccessKeyId,
				SecretAccessKey: *creds.RoleCredentials.SecretAccessKey,
				SessionToken:    *creds.RoleCredentials.SessionToken,
				Expiration:      time.Unix(*creds.RoleCredentials.Expiration/1000, 0).Format(time.RFC3339),
			}

			output, err := json.Marshal(rawCreds)

			if err != nil {
				return err
			}

			fmt.Println(string(output))

			return nil
		},
	}

	return command
}
