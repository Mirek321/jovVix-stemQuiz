package v1

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Improwised/jovvix/api/constants"
	quizUtilsHelper "github.com/Improwised/jovvix/api/helpers/utils"
	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/pkg/structs"
	"github.com/Improwised/jovvix/api/utils"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// SessionStatus returns the mode, invitation code and participant progress of
// a session. Only the session's admin can read it.
// swagger:route GET /v1/active_quizzes/{session_id}/status ActiveQuiz RequestSessionStatus
//
//	Get session mode, code and participant progress (session admin only).
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseSessionStatus
//		  403: GenericResFailNotFound
//		  404: GenericResFailNotFound
//		  500: GenericResError
func (qc *quizSocketController) SessionStatus(c *fiber.Ctx) error {
	sessionId := c.Params(constants.SessionIDParam)
	userId := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))

	if _, err := uuid.Parse(sessionId); err != nil {
		return utils.JSONFail(c, http.StatusNotFound, constants.ErrSessionNotFound)
	}

	session, err := qc.activeQuizModel.GetSession(sessionId)
	if err != nil {
		if err.Error() == constants.ErrSessionNotFound {
			return utils.JSONFail(c, http.StatusNotFound, constants.ErrSessionNotFound)
		}
		qc.logger.Error("error getting session for status", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.UnknownError)
	}

	if session.AdminID != userId {
		return utils.JSONFail(c, http.StatusForbidden, constants.ErrSessionStatusForbidden)
	}

	totalQuestions, err := qc.questionModel.GetTotalQuestionCount(session.ID.String())
	if err != nil {
		qc.logger.Error(constants.ErrInGettingTotalQuestionCount, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.UnknownError)
	}

	players, err := qc.userPlayedQuizModel.GetSessionPlayers(session.ID)
	if err != nil {
		qc.logger.Error("error getting session players for status", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.UnknownError)
	}

	var code any
	if session.InvitationCode.Valid {
		code = session.InvitationCode.Int32
	}

	return utils.JSONSuccess(c, http.StatusOK, map[string]any{
		"session_id":      session.ID.String(),
		"mode":            session.Mode,
		"code":            code,
		"title":           session.Title,
		"is_active":       session.IsActive,
		"total_questions": totalQuestions,
		"players":         players,
	})
}

// validateSelfPacedRequest checks the common query parameters of the
// self-paced player endpoints and returns the session and the played-quiz id.
// On failure the error response has already been written (ok == false).
func (qc *quizSocketController) validateSelfPacedRequest(c *fiber.Ctx) (models.ActiveQuiz, uuid.UUID, bool, error) {
	currentQuiz := c.Query(constants.CurrentUserQuiz)
	sessionId := c.Query(constants.SessionIDParam)

	if currentQuiz == "" || sessionId == "" {
		return models.ActiveQuiz{}, uuid.Nil, false, utils.JSONFail(c, http.StatusBadRequest, constants.ErrQuizNotFound)
	}

	currentQuizId, err := uuid.Parse(currentQuiz)
	if err != nil {
		return models.ActiveQuiz{}, uuid.Nil, false, utils.JSONFail(c, http.StatusBadRequest, "Neplatné UUID")
	}
	if _, err := uuid.Parse(sessionId); err != nil {
		return models.ActiveQuiz{}, uuid.Nil, false, utils.JSONFail(c, http.StatusBadRequest, "Neplatné UUID")
	}

	user, ok := quizUtilsHelper.ConvertType[models.User](c.Locals(constants.ContextUser))
	if !ok {
		qc.logger.Error("Unable to convert to user-model type from locals")
		return models.ActiveQuiz{}, uuid.Nil, false, utils.JSONFail(c, http.StatusInternalServerError, "Nepodařilo se načíst údaje uživatele")
	}

	isParticipant, err := qc.userPlayedQuizModel.IsParticipant(user.ID, currentQuizId, sessionId)
	if err != nil {
		qc.logger.Error("error validating self-paced participant", zap.Error(err))
		return models.ActiveQuiz{}, uuid.Nil, false, utils.JSONFail(c, http.StatusInternalServerError, constants.UnknownError)
	}
	if !isParticipant {
		return models.ActiveQuiz{}, uuid.Nil, false, utils.JSONFail(c, http.StatusForbidden, constants.ErrNotParticipant)
	}

	session, err := qc.activeQuizModel.GetSession(sessionId)
	if err != nil {
		if err.Error() == constants.ErrSessionNotFound {
			return models.ActiveQuiz{}, uuid.Nil, false, utils.JSONFail(c, http.StatusBadRequest, constants.ErrSessionNotFound)
		}
		qc.logger.Error("error getting self-paced session", zap.Error(err))
		return models.ActiveQuiz{}, uuid.Nil, false, utils.JSONFail(c, http.StatusInternalServerError, constants.UnknownError)
	}

	if failMessage := selfPacedSessionError(session); failMessage != "" {
		return models.ActiveQuiz{}, uuid.Nil, false, utils.JSONFail(c, http.StatusBadRequest, failMessage)
	}

	return session, currentQuizId, true, nil
}

// selfPacedSessionError returns a user-facing message when the session can't
// serve self-paced requests, or an empty string when it can.
func selfPacedSessionError(session models.ActiveQuiz) string {
	if session.Mode != constants.ModeSelfPaced {
		return constants.ErrSessionNotSelfPaced
	}
	if !session.IsActive {
		return constants.ErrSelfPacedQuizEnded
	}
	return ""
}

// expireSelfPacedQuestions records a timeout for every delivered, unanswered
// question whose deadline (duration + grace) has passed and marks them as
// answered in the given slice.
func (qc *quizSocketController) expireSelfPacedQuestions(userPlayedQuizId uuid.UUID, rows []models.SelfPacedResponseRow, now time.Time) error {
	for index := range rows {
		row := &rows[index]
		if row.Answered || !row.DeliveredAt.Valid {
			continue
		}
		if !utils.IsSelfPacedTimedOut(now, utcTime(row.DeliveredAt.Time), row.DurationInSeconds) {
			continue
		}
		if _, err := qc.userQuizResponseModel.RecordSelfPacedTimeout(userPlayedQuizId, row.QuestionID, row.DurationInSeconds); err != nil {
			return err
		}
		row.Answered = true
		row.StreakCount = 0
	}
	return nil
}

// firstUnanswered returns the index of the player's current question, or -1.
func firstUnanswered(rows []models.SelfPacedResponseRow) int {
	for index, row := range rows {
		if !row.Answered {
			return index
		}
	}
	return -1
}

// SelfPacedNext hands the player their current self-paced question (or the
// finished state). The server keeps delivered_at, so reloads resume the timer.
// swagger:route GET /v1/quiz/self/next Quiz RequestSelfPacedNext
//
//	Get the player's current self-paced question.
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseSelfPacedNext
//		  400: GenericResFailNotFound
//		  403: GenericResFailNotFound
//		  500: GenericResError
func (qc *quizSocketController) SelfPacedNext(c *fiber.Ctx) error {
	session, userPlayedQuizId, ok, err := qc.validateSelfPacedRequest(c)
	if !ok {
		return err
	}

	rows, err := qc.userQuizResponseModel.GetSelfPacedResponses(userPlayedQuizId, session.ID)
	if err != nil {
		qc.logger.Error("error getting self-paced responses", zap.Error(err))
		return utils.JSONFail(c, http.StatusInternalServerError, constants.UnknownError)
	}

	now := time.Now().UTC()
	if err := qc.expireSelfPacedQuestions(userPlayedQuizId, rows, now); err != nil {
		qc.logger.Error("error recording self-paced timeout", zap.Error(err))
		return utils.JSONFail(c, http.StatusInternalServerError, constants.UnknownError)
	}

	currentIndex := firstUnanswered(rows)
	if currentIndex == -1 {
		if err := qc.userPlayedQuizModel.MarkCompleted(userPlayedQuizId, now); err != nil {
			qc.logger.Error("error marking self-paced quiz completed", zap.Error(err))
			return utils.JSONFail(c, http.StatusInternalServerError, constants.UnknownError)
		}
		progress, err := qc.userQuizResponseModel.GetSelfPacedProgress(userPlayedQuizId)
		if err != nil {
			qc.logger.Error("error getting self-paced progress", zap.Error(err))
			return utils.JSONFail(c, http.StatusInternalServerError, constants.UnknownError)
		}
		return utils.JSONSuccess(c, http.StatusOK, map[string]any{
			"finished":    true,
			"total_score": progress.TotalScore,
		})
	}

	current := rows[currentIndex]

	deliveredAt, err := qc.userQuizResponseModel.MarkSelfPacedDelivered(userPlayedQuizId, current.QuestionID, now)
	if err != nil {
		qc.logger.Error("error marking self-paced question delivered", zap.Error(err))
		return utils.JSONFail(c, http.StatusInternalServerError, constants.UnknownError)
	}

	question, err := qc.questionModel.GetCurrentQuestion(current.QuestionID)
	if err != nil {
		qc.logger.Error("error getting self-paced question", zap.Error(err))
		return utils.JSONFail(c, http.StatusInternalServerError, constants.UnknownError)
	}

	options, err := qc.userQuizResponseModel.GetOptionsForPlayer(userPlayedQuizId, question.ID, question.Options)
	if err != nil {
		// Never fall back to the unshuffled options: the answer would later be
		// translated with the stored order and point at a different option.
		qc.logger.Error("error getting player-specific option order", zap.Error(err))
		return utils.JSONFail(c, http.StatusInternalServerError, constants.UnknownError)
	}

	progress, err := qc.userQuizResponseModel.GetSelfPacedProgress(userPlayedQuizId)
	if err != nil {
		qc.logger.Error("error getting self-paced progress", zap.Error(err))
		return utils.JSONFail(c, http.StatusInternalServerError, constants.UnknownError)
	}

	// Same field names as the live send_question payload.
	return utils.JSONSuccess(c, http.StatusOK, map[string]any{
		"finished":       false,
		"id":             question.ID,
		"no":             current.OrderNo,
		"totalQuestions": len(rows),
		"duration":       question.DurationInSeconds,
		"points":         question.Points,
		"question":       question.Question,
		"options":        options,
		"question_type":  current.QuestionType,
		"question_media": question.QuestionMedia,
		"options_media":  question.OptionsMedia,
		"resource":       question.Resource.String,
		"start_time":     deliveredAt.Format(time.RFC3339Nano),
		"server_time":    time.Now().UTC().Format(time.RFC3339Nano),
		"total_score":    progress.TotalScore,
	})
}

// setSelfPacedAnswer handles /quiz/answer for self-paced sessions. Timing is
// computed on the server from delivered_at; the client response_time is ignored.
func (qc *quizSocketController) setSelfPacedAnswer(c *fiber.Ctx, session models.ActiveQuiz, userPlayedQuizId uuid.UUID) error {
	if failMessage := selfPacedSessionError(session); failMessage != "" {
		return utils.JSONFail(c, http.StatusBadRequest, failMessage)
	}

	var answer structs.ReqAnswerSubmit
	if err := json.Unmarshal(c.Body(), &answer); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrAnswerSubmit)
	}
	// Relaxed validation: empty keys mean "no answer / timed out" and
	// response_time is computed server-side.
	if answer.QuestionId == uuid.Nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrAnswerSubmit)
	}
	if answer.AnswerKeys == nil {
		answer.AnswerKeys = []int{}
	}
	displayedKeys := answer.AnswerKeys

	rows, err := qc.userQuizResponseModel.GetSelfPacedResponses(userPlayedQuizId, session.ID)
	if err != nil {
		qc.logger.Error("error getting self-paced responses", zap.Error(err))
		return utils.JSONFail(c, http.StatusInternalServerError, constants.UnknownError)
	}

	currentIndex := firstUnanswered(rows)
	if currentIndex == -1 || rows[currentIndex].QuestionID != answer.QuestionId || !rows[currentIndex].DeliveredAt.Valid {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrSelfPacedQuestionNotActive)
	}
	current := rows[currentIndex]
	deliveredAt := utcTime(current.DeliveredAt.Time)

	answers, answerPoints, answerDurationInSeconds, questionType, err := qc.questionModel.GetAnswersPointsDurationType(answer.QuestionId.String())
	if err != nil {
		qc.logger.Error("error while get answer, points, duration and type", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, "Nepodařilo se načíst údaje otázky")
	}

	now := time.Now().UTC()
	// Empty keys come from the client's auto-submit when its timer runs out,
	// usually still inside the grace window, so they count as a timeout too.
	timedOut := len(answer.AnswerKeys) == 0 || utils.IsSelfPacedTimedOut(now, deliveredAt, answerDurationInSeconds)

	var (
		points       sql.NullInt16
		finalScore   int
		newStreak    int
		isCorrect    bool
		originalKeys = []int{}
	)

	if timedOut {
		recorded, err := qc.userQuizResponseModel.RecordSelfPacedTimeout(userPlayedQuizId, answer.QuestionId, answerDurationInSeconds)
		if err != nil {
			qc.logger.Error("error recording self-paced timeout", zap.Error(err))
			return utils.JSONFail(c, http.StatusInternalServerError, constants.UnknownError)
		}
		if !recorded {
			return utils.JSONFail(c, http.StatusBadRequest, constants.ErrAnswerAlreadySubmitted)
		}
	} else {
		originalKeys, err = qc.userQuizResponseModel.TranslateAnswerKeys(userPlayedQuizId, answer.QuestionId, displayedKeys)
		if err != nil {
			qc.logger.Error("error translating player answer keys", zap.Error(err))
			return utils.JSONFail(c, http.StatusBadRequest, "Neplatná možnost odpovědi")
		}

		answer.AnswerKeys = originalKeys
		answer.ResponseTime = utils.SelfPacedResponseTimeMs(now, deliveredAt, answerDurationInSeconds)

		var score int
		points, score = utils.CalculatePointsAndScore(answer, answers, answerPoints, answerDurationInSeconds, questionType)

		// The previous question in the player's own order carries the streak.
		previousStreak := 0
		if currentIndex > 0 {
			previousStreak = rows[currentIndex-1].StreakCount
		}
		finalScore, newStreak = utils.CalculateStreakScore(previousStreak, score)

		if err := qc.userQuizResponseModel.SubmitAnswer(userPlayedQuizId, answer, points, finalScore, newStreak); err != nil {
			if err == sql.ErrNoRows {
				return utils.JSONFail(c, http.StatusBadRequest, constants.ErrAnswerAlreadySubmitted)
			}
			qc.logger.Error("error during self-paced answer submit", zap.Error(err))
			return utils.JSONFail(c, http.StatusInternalServerError, constants.UnknownError)
		}

		isCorrect = utils.IsAnswerCorrect(originalKeys, answers, questionType)
	}

	// Map keys back to the numbering the player saw.
	correctKeys, selectedKeys := answers, originalKeys
	question, err := qc.questionModel.GetCurrentQuestion(answer.QuestionId)
	if err != nil {
		qc.logger.Error("error getting self-paced question options", zap.Error(err))
	} else {
		_, correctKeys, selectedKeys, err = qc.userQuizResponseModel.ApplyOptionOrder(userPlayedQuizId, answer.QuestionId, question.Options, answers, originalKeys)
		if err != nil {
			qc.logger.Error("error applying player option order to answer result", zap.Error(err))
			correctKeys, selectedKeys = answers, originalKeys
		}
	}
	if correctKeys == nil {
		correctKeys = []int{}
	}
	if selectedKeys == nil {
		selectedKeys = []int{}
	}

	progress, err := qc.userQuizResponseModel.GetSelfPacedProgress(userPlayedQuizId)
	if err != nil {
		qc.logger.Error("error getting self-paced progress", zap.Error(err))
		return utils.JSONFail(c, http.StatusInternalServerError, constants.UnknownError)
	}

	finished := progress.TotalQuestions > 0 && progress.Answered >= progress.TotalQuestions
	if finished {
		if err := qc.userPlayedQuizModel.MarkCompleted(userPlayedQuizId, now); err != nil {
			qc.logger.Error("error marking self-paced quiz completed", zap.Error(err))
		}
	}

	pointsValue := 0
	if points.Valid {
		pointsValue = int(points.Int16)
	}

	return utils.JSONSuccess(c, http.StatusAccepted, map[string]any{
		"points":          pointsValue,
		"score":           finalScore,
		"streak":          newStreak,
		"is_correct":      isCorrect,
		"timed_out":       timedOut,
		"correct_keys":    correctKeys,
		"selected_keys":   selectedKeys,
		"total_score":     progress.TotalScore,
		"answered":        progress.Answered,
		"total_questions": progress.TotalQuestions,
		"finished":        finished,
	})
}

// utcTime reinterprets a TIMESTAMP (without time zone) value as UTC wall clock.
func utcTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
}
