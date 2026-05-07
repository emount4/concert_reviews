BEGIN;

-- ==========================================
-- 1. ARTISTS
-- ==========================================
CREATE OR REPLACE FUNCTION fn_init_artist_stats() RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO artist_stats (artist_id) 
    VALUES (NEW.artist_id) 
    ON CONFLICT (artist_id) DO NOTHING;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_init_artist_stats ON artists;
CREATE TRIGGER trg_init_artist_stats
AFTER INSERT ON artists
FOR EACH ROW EXECUTE FUNCTION fn_init_artist_stats();

-- ==========================================
-- 2. VENUES
-- ==========================================
CREATE OR REPLACE FUNCTION fn_init_venue_stats() RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO venue_stats (venue_id) 
    VALUES (NEW.venue_id) 
    ON CONFLICT (venue_id) DO NOTHING;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_init_venue_stats ON venues;
CREATE TRIGGER trg_init_venue_stats
AFTER INSERT ON venues
FOR EACH ROW EXECUTE FUNCTION fn_init_venue_stats();

-- ==========================================
-- 3. CONCERTS
-- ==========================================
CREATE OR REPLACE FUNCTION fn_init_concert_stats() RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO concert_stats (concert_id) 
    VALUES (NEW.concert_id) 
    ON CONFLICT (concert_id) DO NOTHING;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_init_concert_stats ON concerts;
CREATE TRIGGER trg_init_concert_stats
AFTER INSERT ON concerts
FOR EACH ROW EXECUTE FUNCTION fn_init_concert_stats();

-- ==========================================
-- 4. USERS
-- ==========================================
CREATE OR REPLACE FUNCTION fn_init_user_stats() RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO user_stats (user_id) 
    VALUES (NEW.user_id) 
    ON CONFLICT (user_id) DO NOTHING;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_init_user_stats ON users;
CREATE TRIGGER trg_init_user_stats
AFTER INSERT ON users
FOR EACH ROW EXECUTE FUNCTION fn_init_user_stats();



CREATE OR REPLACE FUNCTION fn_check_favorites_limit() RETURNS TRIGGER AS $$
DECLARE
    current_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO current_count 
    FROM favorites 
    WHERE user_id = NEW.user_id AND target_type = NEW.target_type;

    IF current_count >= 5 THEN
        RAISE EXCEPTION 'favorites_limit_exceeded: user can have max 5 favorites of type %', NEW.target_type
        USING ERRCODE = '23514'; -- check_violation
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_check_favorites_limit ON favorites;
CREATE TRIGGER trg_check_favorites_limit
BEFORE INSERT ON favorites
FOR EACH ROW EXECUTE FUNCTION fn_check_favorites_limit();

COMMIT;