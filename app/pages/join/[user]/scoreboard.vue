<script setup>
const url = "/final_score/user";
const route = useRoute();
const username = computed(() => route.params.user);
const userPlayedQuiz = computed(() => route.query.user_played_quiz);
const isHost = computed(() => route.query.host === "1");

definePageMeta({
  layout: "empty",
  hideSidebar: true,
});

useSeoMeta({
  title: "Výsledky kvízu – STEM_HUB",
  description:
    "Po hře si prohlédni své výsledky, pořadí a analýzu jednotlivých otázek.",
});
</script>

<template>
  <FinalScoreBoard
    :is-admin="false"
    :is-host="isHost"
    :user-name="username"
    :user-u-r-l="url"
    :user-played-quiz="userPlayedQuiz"
  />
</template>
