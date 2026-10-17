import { describe, it, expect } from "vitest";
import {
  buildFeedback,
  toQuestionMessage,
  SURVEY_QUESTION_TYPE,
} from "~/composables/self_paced";

const question = {
  finished: false,
  id: "q-1",
  no: 2,
  totalQuestions: 5,
  duration: 30,
  points: 1,
  question: "Kolik je 2 + 2?",
  options: { 1: "3", 2: "4", 3: "5" },
  question_type: 1,
  question_media: "text",
  options_media: "text",
  resource: "",
  start_time: "2026-10-17T10:00:00Z",
  server_time: "2026-10-17T10:00:03Z",
  total_score: 1240,
};

describe("toQuestionMessage", () => {
  it("wraps the payload like the live send_question event", () => {
    const message = toQuestionMessage(question);

    expect(message).toMatchObject({
      status: "success",
      event: "send_question",
      component: "Question",
    });
    expect(message.data).toMatchObject({
      id: "q-1",
      no: 2,
      totalQuestions: 5,
      start_time: question.start_time,
      server_time: question.server_time,
      options: question.options,
    });
    expect(message.data).not.toHaveProperty("finished");
    expect(message.data).not.toHaveProperty("total_score");
  });
});

describe("buildFeedback", () => {
  const result = {
    points: 1,
    score: 640,
    streak: 2,
    is_correct: true,
    timed_out: false,
    correct_keys: [2],
    selected_keys: [2],
    total_score: 1880,
    answered: 2,
    total_questions: 5,
    finished: false,
  };

  it("marks a correct answer and resolves the correct option text", () => {
    expect(buildFeedback(result, question)).toEqual({
      outcome: "correct",
      score: 640,
      streak: 2,
      totalScore: 1880,
      answered: 2,
      totalQuestions: 5,
      finished: false,
      correctOptions: [{ key: 2, letter: "B", value: "4" }],
      optionsMedia: "text",
    });
  });

  it("marks a wrong answer", () => {
    const feedback = buildFeedback(
      { ...result, is_correct: false, score: 0, streak: 0, selected_keys: [1] },
      question
    );
    expect(feedback.outcome).toBe("wrong");
    expect(feedback.score).toBe(0);
    expect(feedback.correctOptions).toHaveLength(1);
  });

  it("prefers the timeout outcome over correctness", () => {
    const feedback = buildFeedback(
      { ...result, is_correct: false, timed_out: true, selected_keys: [] },
      question
    );
    expect(feedback.outcome).toBe("timeout");
  });

  it("treats survey questions as saved answers without a correct option", () => {
    const feedback = buildFeedback(result, {
      ...question,
      question_type: SURVEY_QUESTION_TYPE,
    });
    expect(feedback.outcome).toBe("survey");
    expect(feedback.correctOptions).toEqual([]);
  });

  it("passes the last-question flag through and tolerates missing fields", () => {
    const feedback = buildFeedback({ finished: true }, question);
    expect(feedback).toMatchObject({
      outcome: "wrong",
      score: 0,
      totalQuestions: 5,
      finished: true,
      correctOptions: [],
    });
  });
});
