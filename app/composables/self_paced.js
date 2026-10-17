// Self-paced quiz mode: the player walks through the questions on their own,
// driven by REST calls instead of the live websocket. The server owns the
// timing, so a reload resumes the current question where it left off.

import constants from "~~/config/constants";

export const QUIZ_MODES = {
  live: "live",
  selfPaced: "self_paced",
};

export const SURVEY_QUESTION_TYPE = 2;

// Returned by the API once the host has ended the session.
export const QUIZ_ENDED_MESSAGE = "Kvíz už skončil.";

export const FEEDBACK_AUTO_ADVANCE_SECONDS = 5;

// QuestionSpace consumes the live `send_question` socket message, so wrap the
// next-question payload in the same envelope.
export const toQuestionMessage = (payload = {}) => {
  const question = { ...payload };
  delete question.finished;
  delete question.total_score;
  return {
    status: "success",
    event: constants.GetQuestion,
    component: constants.Question,
    data: question,
  };
};

const optionLetter = (key) => String.fromCharCode(64 + Number(key));

// Turn the answer result into what the feedback card shows. `question` is the
// payload the player answered, used to look up option texts by key.
export const buildFeedback = (result = {}, question = {}) => {
  const isSurvey = Number(question.question_type) === SURVEY_QUESTION_TYPE;
  const options = question.options || {};

  let outcome = "wrong";
  if (result.timed_out) outcome = "timeout";
  else if (isSurvey) outcome = "survey";
  else if (result.is_correct) outcome = "correct";

  const correctOptions = isSurvey
    ? []
    : (result.correct_keys || [])
        .map(Number)
        .filter((key) => options[key] !== undefined)
        .map((key) => ({
          key,
          letter: optionLetter(key),
          value: options[key],
        }));

  return {
    outcome,
    score: Number(result.score) || 0,
    streak: Number(result.streak) || 0,
    totalScore: Number(result.total_score) || 0,
    answered: Number(result.answered) || 0,
    totalQuestions:
      Number(result.total_questions) || Number(question.totalQuestions) || 0,
    finished: Boolean(result.finished),
    correctOptions,
    optionsMedia: question.options_media || "text",
  };
};

const playParams = (userPlayedQuiz, sessionId) =>
  `user_played_quiz=${encodeURIComponent(
    userPlayedQuiz
  )}&session_id=${encodeURIComponent(sessionId)}`;

export const fetchNextQuestion = async (apiUrl, userPlayedQuiz, sessionId) => {
  const response = await $fetch(
    `${apiUrl}/quiz/self/next?${playParams(userPlayedQuiz, sessionId)}`,
    {
      method: "GET",
      credentials: "include",
      headers: { Accept: "application/json" },
    }
  );
  return response?.data || {};
};

export const submitSelfPacedAnswer = async (
  apiUrl,
  userPlayedQuiz,
  sessionId,
  answer
) => {
  const response = await $fetch(
    `${apiUrl}/quiz/answer?${playParams(userPlayedQuiz, sessionId)}`,
    {
      method: "POST",
      credentials: "include",
      headers: { Accept: "application/json" },
      body: answer,
    }
  );
  return response?.data || {};
};

// Host view of a session: mode, join code and per-player progress.
export const fetchSessionStatus = async (apiUrl, sessionId) => {
  const response = await $fetch(
    `${apiUrl}/active_quizzes/${encodeURIComponent(sessionId)}/status`,
    {
      method: "GET",
      credentials: "include",
      headers: { Accept: "application/json" },
    }
  );
  return response?.data || null;
};
