BEGIN;

DROP TRIGGER IF EXISTS trg_init_artist_stats ON artists;
DROP TRIGGER IF EXISTS trg_init_venue_stats ON venues;
DROP TRIGGER IF EXISTS trg_init_concert_stats ON concerts;
DROP TRIGGER IF EXISTS trg_init_user_stats ON users;

DROP FUNCTION IF EXISTS fn_init_artist_stats();
DROP FUNCTION IF EXISTS fn_init_venue_stats();
DROP FUNCTION IF EXISTS fn_init_concert_stats();
DROP FUNCTION IF EXISTS fn_init_user_stats();

DROP TRIGGER IF EXISTS trg_check_favorites_limit ON favorites;
DROP FUNCTION IF EXISTS fn_check_favorites_limit();
COMMIT;