package media_service

type FileMetadata struct {
	FileName string
	FileSize int64
}

type UploadItem struct {
	FileName  string
	UploadURL string
	FileKey   string
}
