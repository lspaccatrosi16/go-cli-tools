package config

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/kirsle/configdir"
	"github.com/lspaccatrosi16/go-cli-tools/storage"
)

func GetConfigPath(appName string) string {

	configPath := configdir.LocalConfig(appName)
	err := configdir.MakePath(configPath)

	if err != nil {
		log.Fatalln(err)
	}

	return configPath
}

func GetCredentialsPath(appName string) string {
	return filepath.Join(GetConfigPath(appName), "credentials.json")
}

func decodeConfigFile[T any](reader io.Reader) T {
	var configContents T

	decoder := json.NewDecoder(reader)

	for {
		err := decoder.Decode(&configContents)

		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalln(err)
		}
	}

	return configContents
}

func encodeConfigFile[T any](config T) *bytes.Buffer {
	buf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buf)

	encoder.SetIndent("", "\t")

	err := encoder.Encode(&config)

	if err != nil {
		log.Fatalln(err)
	}

	return buf
}

func ReadConfigFile[T any](path string, defaultJson []byte) T {
	var reader io.Reader

	fh, err := os.Open(path)

	if err != nil {
		if os.IsNotExist(err) {
			reader = bytes.NewReader(defaultJson)
		} else {
			log.Fatalln(err)
		}
	} else {
		reader = fh
	}

	defer fh.Close()
	return decodeConfigFile[T](reader)
}

func ReadCloudConfigFile[T any](bucket storage.StorageProvider, key string) T {
	file := bucket.GetFile(key)
	reader := bytes.NewReader(file)

	return decodeConfigFile[T](reader)
}

func WriteConfigFile[T any](path string, config T) {
	reader := encodeConfigFile[T](config)

	fh, err := os.Create(path)

	if err != nil {
		log.Fatalln(err)
	}

	defer fh.Close()

	io.Copy(fh, reader)
}

func WritCloudConfigFile[T any](bucket storage.StorageProvider, key string, config T) {
	reader := encodeConfigFile[T](config)

	file := reader.Bytes()

	bucket.UploadFile(key, file)
}
