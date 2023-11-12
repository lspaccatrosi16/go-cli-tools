package credential

import (
	_ "embed"
	"fmt"

	"github.com/lspaccatrosi16/go-cli-tools/pkgError"
)

var wrap = pkgError.WrapErrorFactory("credential")

var chosenCredentials *Credential

func GetUserAuth(appName string) (Credential, error) {
	if chosenCredentials != nil {
		return *chosenCredentials, nil
	}

	return GetUserAuthFresh(appName)
}

func GetUserAuthFresh(appName string) (Credential, error) {
	var credential Credential
	manager, err := loadManager()
	if err != nil {
		return credential, err
	}

	cred, err := manager.pick(appName, false)
	if err != nil {
		return credential, err
	}

	err = saveManager(manager)
	if err != nil {
		return credential, err
	}

	if cred == nil {
		return credential, wrap(fmt.Errorf("did not select a credential"))
	}

	chosenCredentials = cred

	return *cred, nil
}

func GetDefaultUserAuth(appName string) (Credential, error) {
	var credential Credential
	cred, err := getLegacyCredentials(appName)

	if err != nil {
		return credential, nil
	}

	if cred.Key == "" || cred.Secret == "" {
		return credential, fmt.Errorf("credential has empty fields")
	}

	return *cred, nil
}

func StandaloneManager() error {
	manager, err := loadManager()

	if err != nil {
		return err
	}

	_, err = manager.pick("", true)
	if err != nil {
		return err
	}

	err = saveManager(manager)
	return err
}
