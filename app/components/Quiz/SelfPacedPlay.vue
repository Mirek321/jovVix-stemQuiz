<script setup>
// Player side of a self-paced session: fetch the current question, answer it
// (or let it time out), show feedback, repeat until the server says finished.
import { usePush } from "notivue";
import { readApiError } from "~~/composables/apiError";
import {
  QUIZ_ENDED_MESSAGE,
  buildFeedback,
  fetchNextQuestion,
  submitSelfPacedAnswer,
  toQuestionMessage,
} from "~~/composables/self_paced";

const props = defineProps({
  userPlayedQuiz: {
    type: String,
    required: true,
  },
  sessionId: {
    type: String,
    required: true,
  },
  quizTitle: {
    default: "",
    type: String,
    required: false,
  },
});
// finished(endedByHost) -> show the final scoreboard;
// unavailable(message) -> the session cannot be played any more.
const emits = defineEmits(["finished", "unavailable"]);

const toast = usePush();
const { apiUrl } = useRuntimeConfig().public;

const question = ref(null);
const questionMessage = ref({});
const feedback = ref(null);
const loading = ref(false);
const submitting = ref(false);
// Bumped on every load so QuestionSpace remounts with fresh timer state, even
// when the server hands back the same question after a failed submit.
const questionKey = ref(0);
let questionShownAt = 0;

const handleRequestError = (error, fallback) => {
  const message = readApiError(error, fallback);
  if (message === QUIZ_ENDED_MESSAGE) {
    emits("finished", true);
  } else if (!error?.response) {
    // Network hiccup: stay put, the server keeps the question's timer.
    toast.error(`${fallback} Zkontroluj připojení a zkus to znovu.`);
  } else {
    emits("unavailable", message);
  }
};

const loadNext = async () => {
  if (loading.value) return;
  loading.value = true;
  try {
    const next = await fetchNextQuestion(
      apiUrl,
      props.userPlayedQuiz,
      props.sessionId
    );
    if (next.finished) {
      emits("finished", false);
      return;
    }
    question.value = next;
    questionMessage.value = toQuestionMessage(next);
    questionKey.value += 1;
    questionShownAt = Date.now();
    feedback.value = null;
  } catch (error) {
    handleRequestError(error, "Další otázku se nepodařilo načíst.");
  } finally {
    loading.value = false;
  }
};

const submitAnswer = async (keys) => {
  if (submitting.value || !question.value || feedback.value) return;
  submitting.value = true;
  try {
    const result = await submitSelfPacedAnswer(
      apiUrl,
      props.userPlayedQuiz,
      props.sessionId,
      {
        id: question.value.id,
        keys,
        // Informational only: the server measures the time itself.
        response_time: Date.now() - questionShownAt,
      }
    );
    feedback.value = buildFeedback(result, question.value);
  } catch (error) {
    const message = readApiError(error, "Odpověď se nepodařilo odeslat.");
    if (message === QUIZ_ENDED_MESSAGE) {
      emits("finished", true);
      return;
    }
    // Typically the question is no longer current (answered in another tab,
    // or timed out on the server) — resync with whatever the server has now.
    toast.error(message);
    await loadNext();
  } finally {
    submitting.value = false;
  }
};

onMounted(loadNext);
</script>

<template>
  <QuizSelfPacedFeedback
    v-if="feedback"
    :feedback="feedback"
    :pending="loading"
    @next="loadNext"
  />
  <QuizQuestionSpace
    v-else-if="question"
    :key="questionKey"
    :data="questionMessage"
    :is-admin="false"
    :quiz-title="quizTitle"
    @send-answer="submitAnswer"
    @time-up="submitAnswer([])"
  />
  <QuizWaitingSpacePlayerSkeleton v-else />
</template>
