package credential

import (
	"os"

	_ "embed"

	"github.com/lspaccatrosi16/go-cli-tools/config"
	"github.com/lspaccatrosi16/go-cli-tools/input"
	"github.com/lspaccatrosi16/go-cli-tools/logging"
	"github.com/lspaccatrosi16/go-cli-tools/pkgError"
)

var wrap = pkgError.WrapErrorFactory("credential")

//go:embed baseCredential.json
var baseJson []byte

func readCredentialFromFile(appName string) (Credential, error) {
	credPath, err := config.GetCredentialsPath(appName)

	if err != nil {
		return *new(Credential), wrap(err)
	}

	cred, err := config.ReadConfigFile[Credential](credPath, baseJson)

	if err != nil {
		return *new(Credential), wrap(err)
	}

	return cred, nil
}

func getNewCredentials(appName string) (Credential, error) {
	logger := logging.GetLogger()

	logger.Log("User credentials needed")
	logger.LogDivider()

	key := input.GetInput("Access Key Id")
	secret := input.GetInput("Access Key Secret")
	logger.LogDivider()

	userCred := Credential{Key: key, Secret: secret}

	credPath, err := config.GetCredentialsPath(appName)

	if err != nil {
		return *new(Credential), wrap(err)
	}

	config.WriteConfigFile(credPath, userCred)

	return userCred, nil
}

func getEnvCredential() (bool, Credential) {
	var envCredentials Credential

	envKey := os.Getenv("AWS_ACCESS_KEY_ID")
	envSecret := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if envKey != "" && envSecret != "" {
		envCredentials.Key = envKey
		envCredentials.Secret = envSecret

		return true, envCredentials
	}

	return false, envCredentials

}

func GetUserAuth(appName string) (Credential, error) {
	if b, c := getEnvCredential(); b {
		return c, nil
	}

	cfg, err := readCredentialFromFile(appName)

	if err != nil {
		return *new(Credential), wrap(err)
	}

	if cfg.Key == "" || cfg.Secret == "" {
		cfg, err = getNewCredentials(appName)
	}

	if err != nil {
		return *new(Credential), wrap(err)
	}

	return cfg, nil
}

func RefreshUserCredentials(appName string) (Credential, error) {
	return getNewCredentials(appName)
}
