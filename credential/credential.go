package credential

import (
	"os"

	_ "embed"

	"github.com/lspaccatrosi16/go-cli-tools/config"
	"github.com/lspaccatrosi16/go-cli-tools/input"
	"github.com/lspaccatrosi16/go-cli-tools/logging"
)

//go:embed baseCredential.json
var baseJson []byte

func readCredentialFromFile(appName string) (credential, error) {
	credPath, err := config.GetCredentialsPath(appName)

	if err != nil {
		return *new(credential), err
	}

	cred, err := config.ReadConfigFile[credential](credPath, baseJson)

	if err != nil {
		return *new(credential), err
	}

	return cred, nil
}

func getNewCredentials(appName string) (credential, error) {
	logger := logging.GetLogger()

	logger.Log("User credentials needed")
	logger.LogDivider()

	key := input.GetInput("Access Key Id")
	secret := input.GetInput("Access Key Secret")
	logger.LogDivider()

	userCred := credential{Key: key, Secret: secret}

	credPath, err := config.GetCredentialsPath(appName)

	if err != nil {
		return *new(credential), err
	}

	config.WriteConfigFile[credential](credPath, userCred)

	return userCred, nil
}

func getEnvCredential() (bool, credential) {
	var envCredentials credential

	envKey := os.Getenv("AWS_ACCESS_KEY_ID")
	envSecret := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if envKey != "" && envSecret != "" {
		envCredentials.Key = envKey
		envCredentials.Secret = envSecret

		return true, envCredentials
	}

	return false, envCredentials

}

func GetUserAuth(appName string) (credential, error) {
	if b, c := getEnvCredential(); b {
		return c, nil
	}

	cfg, err := readCredentialFromFile(appName)

	if err != nil {
		return *new(credential), err
	}

	if cfg.Key == "" || cfg.Secret == "" {
		cfg, err = getNewCredentials(appName)
	}

	if err != nil {
		return *new(credential), err
	}

	return cfg, nil
}
