BEGIN;

CREATE TABLE IF NOT EXISTS languages (
    id INT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
    );

CREATE TABLE IF NOT EXISTS levels (
    id INT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
    );



CREATE TABLE IF NOT EXISTS users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    level_id INT NOT NULL,
    lang_id INT NOT NULL,
    words_per_day INT NOT NULL,

    CONSTRAINT fk_users_levels
    FOREIGN KEY (level_id) REFERENCES levels(id),

    CONSTRAINT fk_users_languages
    FOREIGN KEY (lang_id) REFERENCES languages(id)
    );



CREATE TABLE IF NOT EXISTS sessions (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    user_id BIGINT NOT NULL,
    status VARCHAR(50) NOT NULL,
    lang_id INT NOT NULL,
    started_at TIMESTAMP NOT NULL,
    ended_at TIMESTAMP,
    accuracy DOUBLE PRECISION,

    CONSTRAINT fk_sessions_users
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    CONSTRAINT fk_sessions_languages
    FOREIGN KEY (lang_id) REFERENCES languages(id)
    );

CREATE TABLE IF NOT EXISTS flashcards (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    word VARCHAR(100) NOT NULL,
    transl VARCHAR(100) NOT NULL,
    lang_id INT NOT NULL,

    CONSTRAINT uq_flashcards_user_word_lang
    UNIQUE (user_id, word, lang_id),

    CONSTRAINT fk_flashcards_users
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    CONSTRAINT fk_flashcards_languages
    FOREIGN KEY (lang_id) REFERENCES languages(id)
    );

CREATE TABLE IF NOT EXISTS decks (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    session_id BIGINT NOT NULL,
    CONSTRAINT uq_decks_user_session
    UNIQUE (user_id, session_id),

    CONSTRAINT fk_decks_users
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    CONSTRAINT fk_decks_sessions
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
    );


CREATE TABLE IF NOT EXISTS decks_flashcards (
    deck_id BIGINT NOT NULL,
    flashcard_id BIGINT NOT NULL,

    PRIMARY KEY (deck_id, flashcard_id),
    CONSTRAINT fk_df_decks
    FOREIGN KEY (deck_id) REFERENCES decks(id) ON DELETE CASCADE,

    CONSTRAINT fk_df_flashcards
    FOREIGN KEY (flashcard_id) REFERENCES flashcards(id) ON DELETE CASCADE
    );

INSERT INTO languages (id, name)
VALUES
    (1, 'english'),
    (2, 'deutsch')
    ON CONFLICT DO NOTHING;

INSERT INTO levels (id, name)
VALUES
        (1,'A1'),
        (2, 'A2')
    ON CONFLICT DO NOTHING;

COMMIT;
