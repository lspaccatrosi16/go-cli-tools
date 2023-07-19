package storage

type StorageProvider interface {
	GetFile(key string) []byte
	UploadFile(key string, file []byte)
	GetTemporaryUrl(key string, expiry int) string
	ListKeys() []string
}
