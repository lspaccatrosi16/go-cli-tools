package storage

type StorageProvider interface {
	GetFile(key string) []byte
	UploadFile(key string, file []byte)
}
