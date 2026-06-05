package media_service

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_models "github.com/emount4/concert_reviews/internal/core/domain/models"
	core_ports "github.com/emount4/concert_reviews/internal/core/domain/ports"
	"github.com/google/uuid"
)

const (
	PurposeAvatar        = "avatar"
	PurposeBanner        = "banner"
	PurposeReviewMedia   = "review_media"
	PurposeArtistPhoto   = "artist_photo"
	PurposeVenuePhoto    = "venue_photo"
	PurposeConcertPoster = "concert_poster"
)

type uploadPurposeRule struct {
	folder       string
	maxSizeBytes int64
	contentTypes []string
	extensions   map[string]string
	adminOnly    bool
	userScoped   bool
}

type MediaService struct {
	s3            core_ports.S3Provider
	allowedExt    map[string]bool
	minUploadSize int64
	maxUploadSize int64
}

func NewMediaService(s3 core_ports.S3Provider, allowedExtensions map[string]bool, minUploadSize, maxUploadSize int64) *MediaService {
	return &MediaService{
		s3:            s3,
		allowedExt:    allowedExtensions,
		minUploadSize: minUploadSize,
		maxUploadSize: maxUploadSize,
	}
}

func (s *MediaService) PrepareBatchUpload(
	ctx context.Context,
	userID string,
	roleID int,
	purpose string,
	files []core_models.MediaUploadParams,
) ([]core_models.MediaUploadTicket, error) {
	rule, err := s.uploadRule(purpose)
	if err != nil {
		return nil, err
	}
	if rule.adminOnly && roleID < domain.RoleAdminID {
		return nil, fmt.Errorf("purpose %s is available only for admins", purpose)
	}
	if rule.userScoped && strings.TrimSpace(userID) == "" {
		return nil, fmt.Errorf("user id is required for purpose %s", purpose)
	}

	tickets := make([]core_models.MediaUploadTicket, 0, len(files))
	for _, f := range files {
		ext := strings.ToLower(filepath.Ext(f.FileName))
		if ext == "" {
			return nil, fmt.Errorf("file extension is required")
		}

		expectedContentType, ok := rule.extensions[ext]
		if !ok || !s.allowedExt[ext] {
			return nil, fmt.Errorf("file type %s is not allowed for purpose %s", ext, purpose)
		}

		contentType := strings.ToLower(strings.TrimSpace(f.ContentType))
		if contentType == "" {
			return nil, fmt.Errorf("content_type is required")
		}
		if contentType != expectedContentType || !slices.Contains(rule.contentTypes, contentType) {
			return nil, fmt.Errorf("content_type %s is not allowed for extension %s and purpose %s", contentType, ext, purpose)
		}

		maxSize := rule.maxSizeBytes
		if s.maxUploadSize > 0 && s.maxUploadSize < maxSize {
			maxSize = s.maxUploadSize
		}
		if f.FileSize < s.minUploadSize || f.FileSize > maxSize {
			return nil, fmt.Errorf("file size %d is out of allowed range for purpose %s", f.FileSize, purpose)
		}

		fileKey := s.buildFileKey(rule, userID, ext)
		uploadURL, formData, err := s.s3.GetUploadForm(ctx, fileKey, contentType, s.minUploadSize, maxSize, 15*time.Minute)
		if err != nil {
			return nil, fmt.Errorf("failed to get presigned post policy for %s: %w", f.FileName, err)
		}

		tickets = append(tickets, core_models.MediaUploadTicket{
			FileKey:    fileKey,
			UploadURL:  uploadURL,
			UploadForm: formData,
		})
	}

	return tickets, nil
}

func (s *MediaService) uploadRule(purpose string) (uploadPurposeRule, error) {
	switch strings.TrimSpace(purpose) {
	case PurposeAvatar:
		return uploadPurposeRule{
			folder:       "avatars",
			maxSizeBytes: 5 * 1024 * 1024,
			contentTypes: imageContentTypes(),
			extensions:   imageExtensions(),
			userScoped:   true,
		}, nil
	case PurposeBanner:
		return uploadPurposeRule{
			folder:       "banners",
			maxSizeBytes: 10 * 1024 * 1024,
			contentTypes: imageContentTypes(),
			extensions:   imageExtensions(),
			userScoped:   true,
		}, nil
	case PurposeReviewMedia:
		return uploadPurposeRule{
			folder:       "reviews",
			maxSizeBytes: 50 * 1024 * 1024,
			contentTypes: []string{"image/jpeg", "image/png", "image/webp", "image/gif", "video/mp4"},
			extensions: map[string]string{
				".jpg":  "image/jpeg",
				".jpeg": "image/jpeg",
				".png":  "image/png",
				".webp": "image/webp",
				".gif":  "image/gif",
				".mp4":  "video/mp4",
			},
			userScoped: true,
		}, nil
	case PurposeArtistPhoto:
		return adminImageRule("artists"), nil
	case PurposeVenuePhoto:
		return adminImageRule("venues"), nil
	case PurposeConcertPoster:
		return adminImageRule("concerts"), nil
	default:
		return uploadPurposeRule{}, fmt.Errorf("unknown upload purpose %s", purpose)
	}
}

func (s *MediaService) buildFileKey(rule uploadPurposeRule, userID string, ext string) string {
	if rule.userScoped {
		return fmt.Sprintf("%s/%s/%s%s", rule.folder, userID, uuid.New().String(), ext)
	}
	return fmt.Sprintf("%s/%s%s", rule.folder, uuid.New().String(), ext)
}

func adminImageRule(folder string) uploadPurposeRule {
	return uploadPurposeRule{
		folder:       folder,
		maxSizeBytes: 10 * 1024 * 1024,
		contentTypes: imageContentTypes(),
		extensions:   imageExtensions(),
		adminOnly:    true,
	}
}

func imageContentTypes() []string {
	return []string{"image/jpeg", "image/png", "image/webp"}
}

func imageExtensions() map[string]string {
	return map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".webp": "image/webp",
	}
}
