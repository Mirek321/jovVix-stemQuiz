<script setup>
definePageMeta({
  layout: "empty",
});

useSeoMeta({
  title: "Stránka nenalezena – STEM_HUB",
  description:
    "Stránku, kterou hledáš, se nepodařilo najít. Vrať se do systému STEM_HUB a spusť kvíz nebo se k nějakému připoj.",
  robots: "noindex, nofollow",
});

const route = useRoute();

const status = ref(
  route.query.status || "Něco se pokazilo. Zkus to prosím později."
);
const message = ref(route.query.error);
</script>

<template>
  <div class="bg-jv-canvas">
    <Frame :page-title="status?.toUpperCase()" :page-message="message" />
  </div>
</template>
