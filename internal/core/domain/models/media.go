package core_models

type MediaUploadParams struct {
	FileName    string
	FileSize    int64
	ContentType string
}

type MediaUploadTicket struct {
	FileKey    string
	UploadURL  string
	UploadForm map[string]string
}
