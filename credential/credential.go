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

func readCredentialFromFile(appName string) credential {
	credPath := config.GetCredentialsPath(appName)
	return config.ReadConfigFile[credential](credPath, baseJson)
}

func getNewCredentials(appName string) credential {
	logger := logging.GetLogger()

	logger.Log("User credentials needed")
	logger.LogDivider()

	key := input.GetInput("Access Key Id")
	secret := input.GetInput("Access Key Secret")
	logger.LogDivider()

	userCred := credential{Key: key, Secret: secret}

	credPath := config.GetCredentialsPath(appName)
	config.WriteConfigFile[credential](credPath, userCred)

	return userCred
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

func GetUserAuth(appName string) credential {
	if b, c := getEnvCredential(); b {
		return c
	}

	cfg := readCredentialFromFile(appName)

	if cfg.Key == "" || cfg.Secret == "" {
		cfg = getNewCredentials(appName)
	}

	return cfg
}
