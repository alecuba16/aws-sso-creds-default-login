package set

import (
	"aws-sso-creds-default-login/pkg/credentials"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"time"

	"github.com/bigkevmcd/go-configparser"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "set PROFILE",
		Short: "Create a new AWS profile with temporary credentials",
		Long:  "Create a new AWS profile with temporary credentials",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			profile := viper.GetString("profile")
			asdefault := viper.GetBool("asdefault")
			login := viper.GetBool("login")
			homeDir := viper.GetString("home-directory")
			credsPath := fmt.Sprintf("%s/.aws/credentials", homeDir)
			cfgPath := fmt.Sprintf("%s/.aws/config", homeDir)
			if profile == "" {
				profile = args[0]
			}
			if profile == "" {
				return fmt.Errorf("no profile specified")
			}
			creds, _, err := credentials.GetSSOCredentials(profile, homeDir, login)
			if err != nil {
				return err
			}

			credsFile, err := configparser.NewConfigParserFromFile(credsPath)
			if os.IsNotExist(err) {
				// Ensure the new empty credentials file is not readable by others.
				if f, err := os.OpenFile(credsPath, os.O_CREATE, 0600); err != nil {
					return err
				} else {
					f.Close()
				}
				credsFile = configparser.New()
			} else if err != nil {
				return err
			}

			configFile, err := configparser.NewConfigParserFromFile(cfgPath)
			if err != nil {
				return err
			}

			// create a new credentials section
			credsFile.AddSection(args[0])
			configFile.AddSection(fmt.Sprintf("profile %s", args[0]))

			credsFile.Set(args[0], "aws_access_key_id", *creds.RoleCredentials.AccessKeyId)
			credsFile.Set(args[0], "aws_secret_access_key", *creds.RoleCredentials.SecretAccessKey)
			credsFile.Set(args[0], "aws_session_token", *creds.RoleCredentials.SessionToken)

			if asdefault == true {
				// Set default region
				region, err := configFile.Get(fmt.Sprintf("profile %s", args[0]), "region")
				if err != nil {
					return err
				}
				configFile.AddSection("default")
				configFile.Set("default", "region", region)
				credsFile.AddSection("default")
				credsFile.Set("default", "aws_access_key_id", *creds.RoleCredentials.AccessKeyId)
				credsFile.Set("default", "aws_secret_access_key", *creds.RoleCredentials.SecretAccessKey)
				credsFile.Set("default", "aws_session_token", *creds.RoleCredentials.SessionToken)
				credsFile.Set("default", "region", region)

			}

			credsFile.SaveWithDelimiter(credsPath, "=")
			configFile.SaveWithDelimiter(cfgPath, "=")

			fmt.Println(fmt.Sprintf("credentials saved to profile: %s", args[0]))
			fmt.Println(fmt.Sprintf("these credentials will expire:  %s", time.Unix(*creds.RoleCredentials.Expiration, 0).Format(time.UnixDate)))

			return nil
		},
	}

	return command
}
