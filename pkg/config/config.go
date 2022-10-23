package config

import (
	"fmt"
)

// GetSSOConfig retrieves the SSO configuration for a given AWS profile
func GetSSOConfig(profile string, homedir string) (*SSOConfig, error) {

	// parse the configuration file
	p, err := configparser.NewConfigParserFromFile(fmt.Sprintf("%s/.aws/config", homedir))

	if err != nil {
		return nil, err
	}

	// build a section name
	section := fmt.Sprintf("profile %s", profile)

	// FIXME: make this better
	if p.HasSection(section) {
		ssoStartUrl, err := p.Get(section, "sso_start_url")
		if err != nil {
			return nil, fmt.Errorf("no SSO url in profile: %s", profile)
		}
		ssoRegion, err := p.Get(section, "sso_region")
		if err != nil {
			return nil, fmt.Errorf("no SSO region in profile: %s", profile)
		}
		ssoAccountId, err := p.Get(section, "sso_account_id")
		if err != nil {
			return nil, fmt.Errorf("no SSO account id in profile: %s", profile)
		}
		ssoRoleName, err := p.Get(section, "sso_role_name")
		if err != nil {
			return nil, fmt.Errorf("no SSO role name in profile: %s", profile)
		}

		return &SSOConfig{
			StartUrl:  ssoStartUrl,
			Region:    ssoRegion,
			AccountID: ssoAccountId,
			RoleName:  ssoRoleName,
		}, nil

	} else {
		return nil, fmt.Errorf("unable to find profile %s", profile)
	}
}
