-- Quiz questions table
CREATE TABLE IF NOT EXISTS quiz_questions (
    id SERIAL PRIMARY KEY,
    language VARCHAR(50) NOT NULL,
    category VARCHAR(50) NOT NULL,
    difficulty VARCHAR(20) NOT NULL,
    question_text TEXT NOT NULL,
    answer_options JSONB NOT NULL,
    correct_answer VARCHAR(255) NOT NULL,
    explanation TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- User quiz sessions table
CREATE TABLE IF NOT EXISTS user_quiz_sessions (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    language VARCHAR(50) NOT NULL,
    current_question_index INT NOT NULL DEFAULT 0,
    question_ids JSONB NOT NULL,
    correct_answers INT NOT NULL DEFAULT 0,
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,
    UNIQUE (user_id, language, started_at)
);

-- User quiz answers table
CREATE TABLE IF NOT EXISTS user_quiz_answers (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    session_id INT NOT NULL REFERENCES user_quiz_sessions(id),
    question_id INT NOT NULL REFERENCES quiz_questions(id),
    answer_given VARCHAR(255) NOT NULL,
    is_correct BOOLEAN NOT NULL,
    answered_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_quiz_questions_language ON quiz_questions(language);
CREATE INDEX IF NOT EXISTS idx_user_quiz_sessions_user_id ON user_quiz_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_quiz_sessions_active ON user_quiz_sessions(user_id) WHERE completed_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_user_quiz_answers_session ON user_quiz_answers(session_id);