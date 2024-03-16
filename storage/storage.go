package storage

type StorageProvider interface {
	GetFile(key string) ([]byte, error)
	UploadFile(key string, file []byte) error
	DeleteFile(key string) error
	GetTemporaryUrl(key string, expiry int) (string, error)
	ListKeys() ([]string, error)
}
