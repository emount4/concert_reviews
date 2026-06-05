package media_service_test

import (
	"context"
	"strings"
	"testing"
	"time"

	core_models "github.com/emount4/concert_reviews/internal/core/domain/models"
	media_service "github.com/emount4/concert_reviews/internal/features/media/service"
)

type fakeS3Provider struct {
	lastObjectName  string
	lastContentType string
	lastMinSize     int64
	lastMaxSize     int64
	lastExpires     time.Duration
	calls           int
}

func (f *fakeS3Provider) GetUploadURL(ctx context.Context, objectName string, expires time.Duration) (string, error) {
	return "put-url", nil
}

func (f *fakeS3Provider) GetUploadForm(
	ctx context.Context,
	objectName string,
	contentType string,
	minSize int64,
	maxSize int64,
	expires time.Duration,
) (string, map[string]string, error) {
	f.calls++
	f.lastObjectName = objectName
	f.lastContentType = contentType
	f.lastMinSize = minSize
	f.lastMaxSize = maxSize
	f.lastExpires = expires

	return "post-url", map[string]string{"Content-Type": contentType}, nil
}

func (f *fakeS3Provider) FileExists(ctx context.Context, objectName string) (int64, error) {
	return 1, nil
}

func (f *fakeS3Provider) GetPublicURL(objectName string) string {
	return "public/" + objectName
}

func (f *fakeS3Provider) DeleteObject(ctx context.Context, objectName string) error {
	return nil
}

func TestPrepareBatchUploadAvatarUsesUserScopedFolderAndLimit(t *testing.T) {
	s3 := &fakeS3Provider{}
	service := newMediaService(s3, 50*1024*1024)

	tickets, err := service.PrepareBatchUpload(
		context.Background(),
		"user-1",
		1,
		"avatar",
		[]core_models.MediaUploadParams{{
			FileName:    "me.jpg",
			FileSize:    1024,
			ContentType: "image/jpeg",
		}},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tickets) != 1 {
		t.Fatalf("expected one ticket, got %d", len(tickets))
	}
	if !strings.HasPrefix(tickets[0].FileKey, "avatars/user-1/") {
		t.Fatalf("expected avatar key to be user-scoped, got %q", tickets[0].FileKey)
	}
	if !strings.HasSuffix(tickets[0].FileKey, ".jpg") {
		t.Fatalf("expected original extension to be kept, got %q", tickets[0].FileKey)
	}
	if s3.lastContentType != "image/jpeg" {
		t.Fatalf("expected exact content type policy, got %q", s3.lastContentType)
	}
	if s3.lastMaxSize != 5*1024*1024 {
		t.Fatalf("expected avatar max size policy to be 5MB, got %d", s3.lastMaxSize)
	}
	if s3.lastExpires != 15*time.Minute {
		t.Fatalf("expected 15 minute policy, got %s", s3.lastExpires)
	}
}

func TestPrepareBatchUploadRejectsAdminPurposeForRegularUser(t *testing.T) {
	service := newMediaService(&fakeS3Provider{}, 50*1024*1024)

	_, err := service.PrepareBatchUpload(
		context.Background(),
		"user-1",
		1,
		"artist_photo",
		[]core_models.MediaUploadParams{{
			FileName:    "artist.png",
			FileSize:    1024,
			ContentType: "image/png",
		}},
	)
	if err == nil || !strings.Contains(err.Error(), "only for admins") {
		t.Fatalf("expected admin-only error, got %v", err)
	}
}

func TestPrepareBatchUploadAllowsAdminPurposeForAdmin(t *testing.T) {
	s3 := &fakeS3Provider{}
	service := newMediaService(s3, 50*1024*1024)

	tickets, err := service.PrepareBatchUpload(
		context.Background(),
		"admin-1",
		2,
		"concert_poster",
		[]core_models.MediaUploadParams{{
			FileName:    "poster.webp",
			FileSize:    1024,
			ContentType: "image/webp",
		}},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(tickets[0].FileKey, "concerts/") {
		t.Fatalf("expected concert poster folder, got %q", tickets[0].FileKey)
	}
	if strings.Contains(tickets[0].FileKey, "admin-1") {
		t.Fatalf("expected admin media key not to be user-scoped, got %q", tickets[0].FileKey)
	}
	if s3.lastMaxSize != 10*1024*1024 {
		t.Fatalf("expected admin image max size policy to be 10MB, got %d", s3.lastMaxSize)
	}
}

func TestPrepareBatchUploadAllowsReviewVideo(t *testing.T) {
	s3 := &fakeS3Provider{}
	service := newMediaService(s3, 50*1024*1024)

	tickets, err := service.PrepareBatchUpload(
		context.Background(),
		"user-1",
		1,
		"review_media",
		[]core_models.MediaUploadParams{{
			FileName:    "clip.mp4",
			FileSize:    10 * 1024 * 1024,
			ContentType: "video/mp4",
		}},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(tickets[0].FileKey, "reviews/user-1/") {
		t.Fatalf("expected review media key to be user-scoped, got %q", tickets[0].FileKey)
	}
	if s3.lastContentType != "video/mp4" {
		t.Fatalf("expected video content type, got %q", s3.lastContentType)
	}
}

func TestPrepareBatchUploadRejectsUnsupportedPurposeType(t *testing.T) {
	service := newMediaService(&fakeS3Provider{}, 50*1024*1024)

	_, err := service.PrepareBatchUpload(
		context.Background(),
		"user-1",
		1,
		"avatar",
		[]core_models.MediaUploadParams{{
			FileName:    "avatar.gif",
			FileSize:    1024,
			ContentType: "image/gif",
		}},
	)
	if err == nil || !strings.Contains(err.Error(), "not allowed for purpose avatar") {
		t.Fatalf("expected avatar gif rejection, got %v", err)
	}
}

func TestPrepareBatchUploadRejectsContentTypeMismatch(t *testing.T) {
	service := newMediaService(&fakeS3Provider{}, 50*1024*1024)

	_, err := service.PrepareBatchUpload(
		context.Background(),
		"user-1",
		1,
		"review_media",
		[]core_models.MediaUploadParams{{
			FileName:    "photo.png",
			FileSize:    1024,
			ContentType: "image/jpeg",
		}},
	)
	if err == nil || !strings.Contains(err.Error(), "content_type image/jpeg is not allowed") {
		t.Fatalf("expected content type mismatch, got %v", err)
	}
}

func TestPrepareBatchUploadRejectsPurposeSizeLimit(t *testing.T) {
	service := newMediaService(&fakeS3Provider{}, 50*1024*1024)

	_, err := service.PrepareBatchUpload(
		context.Background(),
		"user-1",
		1,
		"banner",
		[]core_models.MediaUploadParams{{
			FileName:    "banner.jpg",
			FileSize:    11 * 1024 * 1024,
			ContentType: "image/jpeg",
		}},
	)
	if err == nil || !strings.Contains(err.Error(), "out of allowed range") {
		t.Fatalf("expected banner size rejection, got %v", err)
	}
}

func TestPrepareBatchUploadUsesGlobalMaxWhenItIsStricter(t *testing.T) {
	s3 := &fakeS3Provider{}
	service := newMediaService(s3, 3*1024*1024)

	_, err := service.PrepareBatchUpload(
		context.Background(),
		"user-1",
		1,
		"avatar",
		[]core_models.MediaUploadParams{{
			FileName:    "avatar.jpg",
			FileSize:    2 * 1024 * 1024,
			ContentType: "image/jpeg",
		}},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s3.lastMaxSize != 3*1024*1024 {
		t.Fatalf("expected stricter global max size policy, got %d", s3.lastMaxSize)
	}
}

func newMediaService(s3 *fakeS3Provider, maxUploadSize int64) *media_service.MediaService {
	return media_service.NewMediaService(
		s3,
		map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".webp": true,
			".gif":  true,
			".mp4":  true,
		},
		0,
		maxUploadSize,
	)
}
