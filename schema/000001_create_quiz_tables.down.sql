-- Drop indexes
DROP INDEX IF EXISTS idx_quiz_questions_language;
DROP INDEX IF EXISTS idx_user_quiz_sessions_user_id;
DROP INDEX IF EXISTS idx_user_quiz_sessions_active;
DROP INDEX IF EXISTS idx_user_quiz_answers_session;

-- Drop tables in the correct order to respect foreign key constraints
DROP TABLE IF EXISTS user_quiz_answers;
DROP TABLE IF EXISTS user_quiz_sessions;
DROP TABLE IF EXISTS quiz_questions;
