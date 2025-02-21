BEGIN;

CREATE TABLE audio_mappings (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    phrase_id VARCHAR(255) NOT NULL,
    object_name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audio_mappings_user_phrase ON audio_mappings (user_id, phrase_id);

COMMIT;