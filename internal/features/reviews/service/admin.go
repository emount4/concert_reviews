package review_service

import (
	"context"
	"fmt"

	"github.com/emount4/concert_reviews/internal/core/domain"
	core_errors "github.com/emount4/concert_reviews/internal/core/errors"
	"github.com/google/uuid"
)

func (s *ReviewService) ApproveReview(
	ctx context.Context,
	id uuid.UUID,
	moderatorID uuid.UUID,
	finalTitle string,
	finalText string,
	allowedMediaIDs []uuid.UUID,
) (domain.Review, error) {

	// 1. Получаем текущую рецензию, чтобы знать баллы и ID объектов
	review, err := s.reviewRepository.GetReviewByID(ctx, id)
	if err != nil {
		return domain.Review{}, fmt.Errorf("review not found: %w", err)
	}

	// 2. Проверяем, не одобрена ли она уже (защита от двойного пересчета)
	if review.Status == domain.StatusApproved {
		return domain.Review{}, fmt.Errorf("review is already approved: %w", core_errors.ErrConflict)
	}

	// 3. Вызываем «Тяжелую транзакцию» в репозитории
	err = s.reviewRepository.ApproveReview(ctx, id, moderatorID, finalTitle, finalText, allowedMediaIDs, review)
	if err != nil {
		return domain.Review{}, fmt.Errorf("failed to execute approval transaction: %w", err)
	}

	// 4. Сбрасываем кэш глобальной статистики в Redis (так как число рецензий изменилось)
	_ = s.statsRedis.InvalidateGlobalStats(ctx)

	// Возвращаем обновленную рецензию
	return s.reviewRepository.GetReviewByID(ctx, id)
}
