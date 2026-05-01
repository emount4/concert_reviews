package media_transport_http

type FileMetadataDTO struct {
	FileName string `json:"filename" validate:"required"`
	FileSize int64  `json:"file_size" validate:"required,max=52428800"` // 50MB limit
}

type BatchUploadRequest struct {
	Files []FileMetadataDTO `json:"files" validate:"required,min=1,max=10"`
}

type UploadItemDTO struct {
	FileKey    string            `json:"file_key"`
	UploadURL  string            `json:"upload_url"`
	UploadForm map[string]string `json:"upload_form"`
}

type BatchUploadResponse struct {
	Items []UploadItemDTO `json:"items"`
}
