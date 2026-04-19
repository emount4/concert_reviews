-- Задание: Миграция схемы БД (контент/социалка/статистика)

-- ==========================================
-- 0. ТИПЫ ДАННЫХ (ENUMS)
-- ==========================================

-- Типы контента для системы "Избранное" и Логов
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'target_type_enum') THEN
        CREATE TYPE target_type_enum AS ENUM ('artist', 'venue', 'concert', 'review');
    END IF;
END $$;

-- Общий статус видимости контента
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'content_status_enum') THEN
        CREATE TYPE content_status_enum AS ENUM ('active', 'hidden', 'archived');
    END IF;
END $$;

-- ==========================================
-- 3. АРТИСТЫ, ПЛОЩАДКИ, КОНЦЕРТЫ
-- ==========================================

CREATE TABLE IF NOT EXISTS venues (
    venue_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    city_id INTEGER NOT NULL REFERENCES cities(city_id) ON DELETE RESTRICT,
    name VARCHAR(255) NOT NULL CHECK (name !~ '^\s*$'),
    address TEXT CHECK (address !~ '^\s*$'),
    capacity INTEGER CHECK (capacity > 0 AND capacity <= 500000),
    social_links JSONB,
    photo_url TEXT CHECK (photo_url IS NULL OR (photo_url !~ '^\s*$' AND LENGTH(photo_url) <= 2048)),
    description TEXT CHECK (LENGTH(description) <= 2000),
    status content_status_enum DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CHECK (deleted_at IS NULL OR deleted_at > created_at)
);

CREATE TABLE IF NOT EXISTS artists (
    artist_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL CHECK (name !~ '^\s*$'),
    description TEXT CHECK (LENGTH(description) <= 2000),
    photo_url TEXT CHECK (photo_url IS NULL OR (photo_url !~ '^\s*$' AND LENGTH(photo_url) <= 2048)),
    social_links JSONB,
    status content_status_enum DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CHECK (deleted_at IS NULL OR deleted_at > created_at)
);

CREATE TABLE IF NOT EXISTS concerts (
    concert_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id UUID NOT NULL REFERENCES venues(venue_id) ON DELETE RESTRICT,
    title VARCHAR(255) NOT NULL CHECK (title !~ '^\s*$'),
    date TIMESTAMP NOT NULL,
    poster_url TEXT CHECK (poster_url IS NULL OR (poster_url !~ '^\s*$' AND LENGTH(poster_url) <= 2048)),
    is_verified BOOLEAN DEFAULT FALSE,
    created_by_user_id UUID REFERENCES users(user_id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CHECK (deleted_at IS NULL OR deleted_at > created_at)
);

-- Многие-ко-многим: Артисты на Концерте
CREATE TABLE IF NOT EXISTS concert_artists (
    concert_id UUID REFERENCES concerts(concert_id) ON DELETE CASCADE,
    artist_id UUID REFERENCES artists(artist_id) ON DELETE CASCADE,
    is_main BOOLEAN DEFAULT TRUE,
    PRIMARY KEY (concert_id, artist_id)
);

-- Предложения концертов от юзеров
CREATE TABLE IF NOT EXISTS concert_suggestions (
    suggestion_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(user_id),
    raw_artist_name VARCHAR(255),
    raw_venue_name VARCHAR(255),
    concert_date TIMESTAMP,
    info TEXT,
    status moderation_status_enum DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW()
);

-- ==========================================
-- 4. РЕЦЕНЗИИ И МЕДИА
-- ==========================================

CREATE TABLE IF NOT EXISTS reviews (
    review_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(user_id) ON DELETE SET NULL,
    concert_id UUID NOT NULL REFERENCES concerts(concert_id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL CHECK (title !~ '^\s*$'),
    text TEXT NOT NULL CHECK (LENGTH(text) BETWEEN 100 AND 8000),
    original_text TEXT CHECK (original_text IS NULL OR LENGTH(original_text) <= 8000),

    p1 SMALLINT NOT NULL CHECK (p1 BETWEEN 1 AND 10),
    p2 SMALLINT NOT NULL CHECK (p2 BETWEEN 1 AND 10),
    p3 SMALLINT NOT NULL CHECK (p3 BETWEEN 1 AND 10),
    p4 SMALLINT NOT NULL CHECK (p4 BETWEEN 1 AND 10),
    p5 SMALLINT NOT NULL CHECK (p5 BETWEEN 1 AND 10),

    rating_total INTEGER CHECK (rating_total <= 90),

    status moderation_status_enum DEFAULT 'pending',
    rejection_reason TEXT CHECK (LENGTH(rejection_reason) <= 500),
    moderated_by_user_id UUID REFERENCES users(user_id) ON DELETE SET NULL,
    is_visible BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    
    CHECK (
        (status IN ('approved', 'rejected') AND moderated_by_user_id IS NOT NULL) OR
        (status = 'pending' AND moderated_by_user_id IS NULL)
    ),
    CHECK (deleted_at IS NULL OR deleted_at > created_at)
);

-- Медиа рецензий с индивидуальной модерацией
CREATE TABLE IF NOT EXISTS review_media (
    media_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    review_id UUID NOT NULL REFERENCES reviews(review_id) ON DELETE CASCADE,
    media_url TEXT NOT NULL CHECK (media_url !~ '^\s*$' AND LENGTH(media_url) <= 2048),
    media_type VARCHAR(20) NOT NULL CHECK (media_type IN ('image', 'video')),
    file_size INTEGER CHECK (file_size > 0 AND file_size <= 104857600),
    status moderation_status_enum DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW()
);

-- ==========================================
-- 5. СОЦИАЛЬНОЕ ВЗАИМОДЕЙСТВИЕ
-- ==========================================

CREATE TABLE IF NOT EXISTS review_likes (
    like_id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    review_id UUID NOT NULL REFERENCES reviews(review_id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT unq_review_likes_user_review UNIQUE (user_id, review_id)
);

CREATE TABLE IF NOT EXISTS favorites (
    favorite_id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    target_id UUID NOT NULL,
    target_type target_type_enum NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT unq_favorites_user_target UNIQUE (user_id, target_id, target_type)
);

-- ==========================================
-- 6. КЭШ СТАТИСТИКИ (Статистические таблицы)
-- ==========================================

CREATE TABLE IF NOT EXISTS concert_stats (
    concert_id UUID PRIMARY KEY REFERENCES concerts(concert_id) ON DELETE CASCADE,
    sum_p1 INTEGER DEFAULT 0 CHECK (sum_p1 >= 0),
    sum_p2 INTEGER DEFAULT 0 CHECK (sum_p2 >= 0),
    sum_p3 INTEGER DEFAULT 0 CHECK (sum_p3 >= 0),
    sum_p4 INTEGER DEFAULT 0 CHECK (sum_p4 >= 0),
    sum_p5 INTEGER DEFAULT 0 CHECK (sum_p5 >= 0),
    sum_rating_total BIGINT DEFAULT 0 CHECK (sum_rating_total >= 0),
    reviews_count INTEGER DEFAULT 0 CHECK (reviews_count >= 0),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS artist_stats (
    artist_id UUID PRIMARY KEY REFERENCES artists(artist_id) ON DELETE CASCADE,
    sum_rating_total BIGINT DEFAULT 0 CHECK (sum_rating_total >= 0),
    reviews_count INTEGER DEFAULT 0 CHECK (reviews_count >= 0),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS venue_stats (
    venue_id UUID PRIMARY KEY REFERENCES venues(venue_id) ON DELETE CASCADE,
    sum_rating_total BIGINT DEFAULT 0 CHECK (sum_rating_total >= 0),
    reviews_count INTEGER DEFAULT 0 CHECK (reviews_count >= 0),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_stats (
    user_id UUID PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE,
    reviews_count INTEGER DEFAULT 0 CHECK (reviews_count >= 0),
    likes_given_count INTEGER DEFAULT 0 CHECK (likes_given_count >= 0),
    likes_received_count INTEGER DEFAULT 0 CHECK (likes_received_count >= 0),
    updated_at TIMESTAMP DEFAULT NOW()
);


-- ==========================================
-- 7. ЛОГИРОВАНИЕ И ИНДЕКСЫ
-- ==========================================

CREATE TABLE IF NOT EXISTS moderation_logs (
    log_id SERIAL PRIMARY KEY,
    moderator_user_id UUID REFERENCES users(user_id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL CHECK (action !~ '^\s*$'),
    target_id UUID,
    target_type target_type_enum,
    details JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Основные индексы для скорости
CREATE INDEX IF NOT EXISTS idx_reviews_status ON reviews(status) WHERE status = 'approved';
CREATE INDEX IF NOT EXISTS idx_concerts_date ON concerts(date DESC);
CREATE INDEX IF NOT EXISTS idx_venues_city ON venues(city_id);
CREATE INDEX IF NOT EXISTS idx_favorites_user ON favorites(user_id);
CREATE INDEX IF NOT EXISTS idx_review_likes_review ON review_likes(review_id);
CREATE INDEX IF NOT EXISTS idx_concert_artists_main ON concert_artists(concert_id) WHERE is_main = TRUE;
CREATE INDEX IF NOT EXISTS idx_cities_slug ON cities(slug);
