-- +migrate Up
-- Session mode: 'live' (admin-driven lockstep) or 'self_paced' (players advance on their own).
ALTER TABLE active_quizzes
ADD COLUMN mode VARCHAR(16) NOT NULL DEFAULT 'live',
ADD CONSTRAINT active_quizzes_mode_check CHECK (mode IN ('live', 'self_paced'));

-- When this player first received this question (self-paced timing source of truth).
ALTER TABLE user_quiz_responses
ADD COLUMN delivered_at TIMESTAMP NULL;

-- When the player has a recorded answer for every question (self-paced).
ALTER TABLE user_played_quizzes
ADD COLUMN completed_at TIMESTAMP NULL;
