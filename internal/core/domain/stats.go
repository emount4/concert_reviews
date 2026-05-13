package domain

import "time"

// ContentStats — агрегированная статистика для артистов и площадок.
type ContentStats struct {
	ReviewsCount   int
	SumRatingTotal int64
	ConcertsCount  int
	FavoritesCount int
	UpdatedAt      time.Time
}
