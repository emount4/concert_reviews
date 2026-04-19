-- Задание: Откат миграции схемы БД (контент/социалка/статистика)

-- Индексы
DROP INDEX IF EXISTS idx_cities_slug;
DROP INDEX IF EXISTS idx_concert_artists_main;
DROP INDEX IF EXISTS idx_review_likes_review;
DROP INDEX IF EXISTS idx_favorites_user;
DROP INDEX IF EXISTS idx_venues_city;
DROP INDEX IF EXISTS idx_concerts_date;
DROP INDEX IF EXISTS idx_reviews_status;

-- Таблицы (в обратном порядке зависимостей)
DROP TABLE IF EXISTS moderation_logs;
DROP TABLE IF EXISTS user_stats;
DROP TABLE IF EXISTS venue_stats;
DROP TABLE IF EXISTS artist_stats;
DROP TABLE IF EXISTS concert_stats;
DROP TABLE IF EXISTS favorites;
DROP TABLE IF EXISTS review_likes;
DROP TABLE IF EXISTS review_media;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS concert_suggestions;
DROP TABLE IF EXISTS concert_artists;
DROP TABLE IF EXISTS concerts;
DROP TABLE IF EXISTS artists;
DROP TABLE IF EXISTS venues;

-- Типы
DROP TYPE IF EXISTS content_status_enum;
DROP TYPE IF EXISTS target_type_enum;
