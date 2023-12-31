package config

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/kirsle/configdir"
	"github.com/lspaccatrosi16/go-cli-tools/internal/pkgError"
	"github.com/lspaccatrosi16/go-cli-tools/storage"
)

var wrap = pkgError.WrapErrorFactory("config")

func GetConfigPath(appName string) (string, error) {
	configPath := configdir.LocalConfig(appName)
	err := configdir.MakePath(configPath)

	if err != nil {
		return "", wrap(err)
	}

	return configPath, nil
}

func GetCredentialsPath(appName string) (string, error) {
	cpath, err := GetConfigPath(appName)
	if err != nil {
		return "", wrap(err)
	}
	return filepath.Join(cpath, "credentials.json"), nil
}

func decodeConfigFile[T any](reader io.Reader) (T, error) {
	var configContents T

	decoder := json.NewDecoder(reader)

	for {
		err := decoder.Decode(&configContents)

		if err == io.EOF {
			break
		} else if err != nil {
			return *new(T), wrap(err)
		}
	}

	return configContents, nil
}

func encodeConfigFile[T any](config T) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buf)

	encoder.SetIndent("", "\t")

	err := encoder.Encode(&config)

	if err != nil {
		return nil, wrap(err)
	}

	return buf, nil
}

func ReadConfigFile[T any](path string, defaultJson []byte) (T, error) {
	var reader io.Reader

	fh, err := os.Open(path)

	if err != nil {
		if os.IsNotExist(err) {
			reader = bytes.NewReader(defaultJson)
		} else {
			return *new(T), wrap(err)
		}
	} else {
		reader = fh
	}

	defer fh.Close()

	decodeRes, err := decodeConfigFile[T](reader)

	if err != nil {
		return *new(T), wrap(err)
	}

	return decodeRes, nil
}

func ReadCloudConfigFile[T any](bucket storage.StorageProvider, key string) (T, error) {
	file, err := bucket.GetFile(key)

	if err != nil {
		return *new(T), wrap(err)
	}

	reader := bytes.NewReader(file)

	decodeRes, err := decodeConfigFile[T](reader)

	if err != nil {
		return *new(T), wrap(err)
	}

	return decodeRes, nil
}

func WriteConfigFile[T any](path string, config T) error {
	reader, err := encodeConfigFile(config)

	if err != nil {
		return wrap(err)
	}

	fh, err := os.Create(path)

	if err != nil {
		return wrap(err)
	}

	defer fh.Close()

	io.Copy(fh, reader)

	return nil
}

func WritCloudConfigFile[T any](bucket storage.StorageProvider, key string, config T) error {
	reader, err := encodeConfigFile(config)

	if err != nil {
		return wrap(err)
	}

	file := reader.Bytes()

	err = bucket.UploadFile(key, file)

	if err != nil {
		return wrap(err)
	}

	return nil
}
