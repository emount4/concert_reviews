package core_models

type MediaUploadParams struct {
	FileName string
	FileSize int64
}

type MediaUploadTicket struct {
	FileKey    string
	UploadURL  string
	UploadForm map[string]string
}
