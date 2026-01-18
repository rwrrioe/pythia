CREATE TABLE IF NOT EXISTS languages
(
    id int,
    language VARCHAR(50)
    )

CREATE TABLE IF NOT EXISTS levels
(
    id int,
    level VARCHAR(2)
    )

CREATE TABLE IF NOT EXISTS users
(
    id int,
    email varchar(100),
    name string NOT NULL,
    level_id int NOT NULL,
    lang_id int NOT NULL,
    words_per_day INT NOT NULL,

    CONSTRAINT pk_users PRIMARY KEY(id),
    CONSTRAINT fk_users_levels FOREIGN KEY(level_id) REFERENCES levels(id),
    CONSTRAINT fk_users_languages FOREIGN KEY(lang_id) REFERENCES languages(id)
    );
CREATE INDEX IF NOT EXISTS idx_id ON users(id)

CREATE TABLE IF NOT EXISTS sessions
(
    id int,
    user_id int,
    started_at timestamp,
    ended_at timestamp,
    accuracy float
)
CREATE INDEX IF NOT EXISTS idx_session_id ON decks(id)

CREATE TABLE IF NOT EXISTS words
(
    id int
    user_id int,
    word varchar(100),
    transl varchar(100),
    description varchar(100),
    lang_id int,

    CONSTRAINT pk_words PRIMARY KEY,
    CONSTRAINT  words_users FOREIGN KEY(user_id) REFERENCES users(id),
    CONSTRAINT fk_words_languages FOREIGN KEY(lang_id) REFERENCES languages(id)
    );
CREATE INDEX IF NOT EXISTS idx_word_id ON words(id)

CREATE TABLE IF NOT EXISTS flashcards
(
    id int
    user_id int,
    word varchar(100),
    transl varchar(100),
    description varchar(100),
    lang_id int,

    CONSTRAINT pk_flashcards PRIMARY KEY (id),
    CONSTRAINT  flashcards_users FOREIGN KEY(user_id) REFERENCES users(id),
    CONSTRAINT fk_words_languages FOREIGN KEY(lang_id) REFERENCES languages(id)
    );
CREATE INDEX IF NOT EXISTS idx_flashcard_id ON flaschards(id)

CREATE TABLE IF NOT EXISTS decks
(
    id int
    user_id int,
    started_at timestamp,
    finished_at timestamp

    CONSTRAINT pk_decks PRIMARY KEY (id),
    CONSTRAINT decks_users FOREIGN KEY(user_id) REFERENCES users(id),
    );
CREATE INDEX IF NOT EXISTS idx_deck_id ON decks(id)


CREATE TABLE IF NOT EXISTS decks_flashcards
(
    deck_id int,
    flashcard_id int,

    CONSTRAINT decks_flashcards_decks FOREIGN KEY(deck_id) REFERENCES decks(id),
    CONSTRAINT decks_flashcards_flashcards FOREIGN KEY(flashcard_id) REFERENCES flashcards(id),

    PRIMARY KEY (deck_id, flashcard_id)
    );

CREATE TABLE IF NOT EXISTS sessions_decks
(
    session_id int,
    deck_id int

    CONSTRAINT sessions_decks_sessions FOREIGN KEY(session_id) REFERENCES sessions(id),
    CONSTRAINT sessions_decks_decks FOREIGN KEY(deck_id) REFERENCES decks(id),
    PRIMARY KEY (session_id, deck_id)
    )

