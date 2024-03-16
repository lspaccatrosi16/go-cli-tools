package credential

import (
	_ "embed"

	"github.com/lspaccatrosi16/go-cli-tools/config"
)

//go:embed baseCredential.json
var baseJson []byte

func readCredentialFromLegacyFile(path string) (*Credential, error) {
	cred, err := config.ReadConfigFile[Credential](path, baseJson)
	if err != nil {
		return new(Credential), wrap(err)
	}
	return &cred, nil
}

func writeCredentialToLegacyFile(path string, cred Credential) error {
	err := config.WriteConfigFile(path, cred)
	return err
}

func getLegacyCredentials(appName string) (*Credential, error) {
	credPath, err := config.GetCredentialsPath(appName)
	if err != nil {
		return nil, wrap(err)
	}

	return readCredentialFromLegacyFile(credPath)
}

func setLegacyCredentials(appName string, cred *Credential) error {
	credPath, err := config.GetCredentialsPath(appName)
	if err != nil {
		return wrap(err)
	}

	return writeCredentialToLegacyFile(credPath, *cred)
}
