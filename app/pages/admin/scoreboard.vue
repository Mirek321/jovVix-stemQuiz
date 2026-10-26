<script setup>
const url = "/final_score/admin";

definePageMeta({
  layout: "empty",
  hideSidebar: true,
});

useSeoMeta({
  title: "Žebříček kvízu – STEM_HUB",
  description:
    "Výsledky a vítězové kvízu v žebříčku. Přepínej mezi tabulkou a stupni vítězů.",
  robots: "noindex, nofollow",
});
</script>

<template>
  <FinalScoreBoard :is-admin="true" :user-u-r-l="url" />
</template>
