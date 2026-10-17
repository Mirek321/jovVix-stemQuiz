-- +migrate Down
ALTER TABLE user_played_quizzes
DROP COLUMN IF EXISTS completed_at;

ALTER TABLE user_quiz_responses
DROP COLUMN IF EXISTS delivered_at;

ALTER TABLE active_quizzes
DROP CONSTRAINT IF EXISTS active_quizzes_mode_check,
DROP COLUMN IF EXISTS mode;
