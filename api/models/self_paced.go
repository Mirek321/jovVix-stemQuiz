package models

import (
	"database/sql"
	"time"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

// SelfPacedResponseRow is one of a player's response rows joined with the
// session question order and the question timing data.
type SelfPacedResponseRow struct {
	QuestionID        uuid.UUID    `db:"question_id"`
	OrderNo           int          `db:"order_no"`
	Answered          bool         `db:"answered"`
	DeliveredAt       sql.NullTime `db:"delivered_at"`
	StreakCount       int          `db:"streak_count"`
	DurationInSeconds int          `db:"duration_in_seconds"`
	QuestionType      int          `db:"type"`
}

// GetSelfPacedResponses returns the player's response rows ordered by the
// session question order (active_quiz_questions.order_no).
func (model *UserQuizResponseModel) GetSelfPacedResponses(userPlayedQuizId uuid.UUID, sessionId uuid.UUID) ([]SelfPacedResponseRow, error) {
	rows := []SelfPacedResponseRow{}

	err := model.db.From(goqu.T(UserQuizResponsesTable).As("uqr")).
		Select(
			goqu.I("uqr.question_id"),
			goqu.I("aqq.order_no"),
			goqu.L("uqr.answers IS NOT NULL").As("answered"),
			goqu.I("uqr.delivered_at"),
			goqu.L("COALESCE(uqr.streak_count, 0)").As("streak_count"),
			goqu.I("q.duration_in_seconds"),
			goqu.I("q.type"),
		).
		Join(goqu.T(ActiveQuizQuestionsTable).As("aqq"), goqu.On(goqu.Ex{
			"aqq.question_id":    goqu.I("uqr.question_id"),
			"aqq.active_quiz_id": sessionId,
		})).
		Join(goqu.T(QuestionTable).As("q"), goqu.On(goqu.I("q.id").Eq(goqu.I("uqr.question_id")))).
		Where(goqu.I("uqr.user_played_quiz_id").Eq(userPlayedQuizId)).
		Order(goqu.I("aqq.order_no").Asc()).
		ScanStructs(&rows)

	if err != nil {
		return nil, err
	}

	return rows, nil
}

// RecordSelfPacedTimeout stores a timed-out (unanswered) self-paced question.
// It only touches rows that are still unanswered, so it can't overwrite a
// concurrently submitted answer.
func (model *UserQuizResponseModel) RecordSelfPacedTimeout(userPlayedQuizId, questionId uuid.UUID, durationInSeconds int) (bool, error) {
	result, err := model.db.Update(UserQuizResponsesTable).Set(
		goqu.Record{
			"answers":           "[]",
			"calculated_points": 0,
			"calculated_score":  0,
			"streak_count":      0,
			"is_attend":         false,
			"response_time":     durationInSeconds * 1000,
			"updated_at":        goqu.L("now()"),
		},
	).Where(
		goqu.I("user_played_quiz_id").Eq(userPlayedQuizId),
		goqu.I("question_id").Eq(questionId),
		goqu.I("answers").IsNull(),
	).Executor().Exec()

	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

// MarkSelfPacedDelivered sets delivered_at the first time a question is handed
// to the player and returns the stored value. Reloading never resets the timer.
// The timestamp is passed from Go (UTC) so all deadline math uses one clock.
func (model *UserQuizResponseModel) MarkSelfPacedDelivered(userPlayedQuizId, questionId uuid.UUID, now time.Time) (time.Time, error) {
	var deliveredAt time.Time

	found, err := model.db.Update(UserQuizResponsesTable).Set(
		goqu.Record{
			"delivered_at": goqu.L("COALESCE(delivered_at, ?::timestamp)", now.UTC().Format("2006-01-02 15:04:05.999999")),
		},
	).Where(
		goqu.I("user_played_quiz_id").Eq(userPlayedQuizId),
		goqu.I("question_id").Eq(questionId),
	).Returning(goqu.I("delivered_at")).Executor().ScanVal(&deliveredAt)

	if err != nil {
		return time.Time{}, err
	}

	if !found {
		return time.Time{}, sql.ErrNoRows
	}

	return asUTC(deliveredAt), nil
}

// SelfPacedProgress summarises a player's progress in a session.
type SelfPacedProgress struct {
	Answered       int `db:"answered"`
	TotalQuestions int `db:"total_questions"`
	TotalScore     int `db:"total_score"`
}

// GetSelfPacedProgress returns how many questions the player has a recorded
// answer for, the number of questions and the total score so far.
func (model *UserQuizResponseModel) GetSelfPacedProgress(userPlayedQuizId uuid.UUID) (SelfPacedProgress, error) {
	var progress SelfPacedProgress

	_, err := model.db.From(UserQuizResponsesTable).
		Select(
			goqu.L("COUNT(*) FILTER (WHERE answers IS NOT NULL)").As("answered"),
			goqu.L("COUNT(*)").As("total_questions"),
			goqu.L("COALESCE(SUM(calculated_score), 0)").As("total_score"),
		).
		Where(goqu.I("user_played_quiz_id").Eq(userPlayedQuizId)).
		ScanStruct(&progress)

	return progress, err
}

// MarkCompleted sets completed_at once for a player that answered everything.
func (model *UserPlayedQuizModel) MarkCompleted(userPlayedQuizId uuid.UUID, now time.Time) error {
	_, err := model.db.Update(UserPlayedQuizTable).Set(
		goqu.Record{
			"completed_at": goqu.L("?::timestamp", now.UTC().Format("2006-01-02 15:04:05.999999")),
			"updated_at":   goqu.L("now()"),
		},
	).Where(
		goqu.I("id").Eq(userPlayedQuizId),
		goqu.I("completed_at").IsNull(),
	).Executor().Exec()

	return err
}

// SessionPlayer is a participant row of the admin session status.
type SessionPlayer struct {
	UserPlayedQuizID string  `json:"user_played_quiz_id" db:"user_played_quiz_id"`
	Username         string  `json:"username" db:"username"`
	FirstName        string  `json:"first_name" db:"first_name"`
	ImageKey         *string `json:"img_key" db:"img_key"`
	Answered         int     `json:"answered" db:"answered"`
	Score            int     `json:"score" db:"score"`
	Completed        bool    `json:"completed" db:"completed"`
}

// GetSessionPlayers lists the session participants (hosts excluded) with
// their progress, sorted by score desc, then username.
func (model *UserPlayedQuizModel) GetSessionPlayers(sessionId uuid.UUID) ([]SessionPlayer, error) {
	players := []SessionPlayer{}

	err := model.db.From(goqu.T(UserPlayedQuizTable).As("upq")).
		Select(
			goqu.L("upq.id::text").As("user_played_quiz_id"),
			goqu.L("COALESCE(u.username, '')").As("username"),
			goqu.L("COALESCE(u.first_name, '')").As("first_name"),
			goqu.L("NULLIF(u.img_key, '')").As("img_key"),
			goqu.L("COUNT(uqr.id) FILTER (WHERE uqr.answers IS NOT NULL)").As("answered"),
			goqu.L("COALESCE(SUM(uqr.calculated_score), 0)").As("score"),
			goqu.L("upq.completed_at IS NOT NULL").As("completed"),
		).
		Join(goqu.T(UserTable).As("u"), goqu.On(goqu.I("u.id").Eq(goqu.I("upq.user_id")))).
		LeftJoin(goqu.T(UserQuizResponsesTable).As("uqr"), goqu.On(goqu.I("uqr.user_played_quiz_id").Eq(goqu.I("upq.id")))).
		Where(goqu.Ex{
			"upq.active_quiz_id": sessionId,
			"upq.is_host":        false,
		}).
		GroupBy(goqu.I("upq.id"), goqu.I("u.id")).
		Order(goqu.L("score").Desc(), goqu.L("username").Asc()).
		ScanStructs(&players)

	if err != nil {
		return nil, err
	}

	return players, nil
}

// asUTC reinterprets a TIMESTAMP (without time zone) value as UTC wall clock.
func asUTC(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
}

// CountQuizQuestions returns how many questions are linked to a quiz.
func (model *QuestionModel) CountQuizQuestions(quizId string) (int64, error) {
	return model.db.From(constants.QuizQuestionsTable).Where(goqu.Ex{"quiz_id": quizId}).Count()
}
