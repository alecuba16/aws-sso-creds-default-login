package credentials

import (
	"aws-sso-creds-default/pkg/config"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sso"
)

func GetSSOCredentials(profile string, homedir string, withLogin bool) (*sso.GetRoleCredentialsOutput, string, error) {
	if withLogin {
		err := LoginSSO(profile, homedir)
		if err != nil {
			return nil, "", err
		}
	}
	ssoConfig, err := config.GetSSOConfig(profile, homedir)
	if err != nil {
		return nil, "", fmt.Errorf("error retrieving SSO config: %w", err)
	}

	cacheFiles, err := ioutil.ReadDir(fmt.Sprintf("%s/.aws/sso/cache", homedir))
	if err != nil {
		return nil, "", fmt.Errorf("error retrieving cache files - perhaps you need to login?: %w", err)
	}

	token, err := config.GetSSOToken(cacheFiles, *ssoConfig, homedir)
	if err != nil {
		return nil, "", fmt.Errorf("error retrieving SSO token from cache files: %w", err)
	}

	sess := session.Must(session.NewSession())
	svc := sso.New(sess, aws.NewConfig().WithRegion(ssoConfig.Region))

	creds, err := svc.GetRoleCredentials(&sso.GetRoleCredentialsInput{
		AccessToken: &token,
		AccountId:   &ssoConfig.AccountID,
		RoleName:    &ssoConfig.RoleName,
	})

	if err != nil {
		return nil, "", fmt.Errorf("error retrieving credentials from AWS: %w", err)
	}

	return creds, ssoConfig.AccountID, nil

}

func LoginSSO(profile string, homedir string) error {
	app := "aws"
	pathCredentials := fmt.Sprintf("%s/.aws/credentials", homedir)
	pathConfig := fmt.Sprintf("%s/.aws/config", homedir)
	arg0 := "sso"
	arg1 := "login"
	arg2 := "--profile"
	arg3 := profile

	cmd := exec.Command(app, arg0, arg1, arg2, arg3)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("AWS_SHARED_CREDENTIALS_FILE=%s", pathCredentials))
	cmd.Env = append(cmd.Env, fmt.Sprintf("AWS_CONFIG_FILE=%s", pathConfig))
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
