<template>
  <main class="flex flex-col gap-6 sm:gap-8 px-4 sm:px-6 md:px-8 py-4 md:py-5">
    <LandingTopActions />
    <HeroSection />
    <JoinLiveQuiz />
    <PublicQuizSection />
  </main>
</template>

<script setup>
import HeroSection from "@/components/landing/HeroSection.vue";
import JoinLiveQuiz from "@/components/Quiz/JoinLiveQuiz.vue";
import LandingTopActions from "@/components/landing/LandingTopActions.vue";
import PublicQuizSection from "@/components/landing/PublicQuizSection.vue";

definePageMeta({
  layout: "empty",
});

const route = useRoute();
const { baseUrl } = useRuntimeConfig().public;
const canonicalUrl = new URL(route.path, baseUrl).href;
const shareImageUrl = new URL("/og-image.png", baseUrl).href;

useHead({
  link: [{ rel: "canonical", href: canonicalUrl }],
});

useSeoMeta({
  title: "STEM_HUB – kvízy na akcích",
  description:
    "Kvízový systém STEM_HUB pro naše akce. Připoj se kódem a soutěž v živém kvízu s okamžitým bodováním.",
  ogTitle: "STEM_HUB – kvízy na akcích",
  ogDescription:
    "Kvízový systém STEM_HUB pro naše akce. Připoj se kódem a soutěž v živém kvízu s okamžitým bodováním.",
  ogType: "website",
  ogUrl: canonicalUrl,
  ogImage: shareImageUrl,
});

onMounted(() => {
  if (route.hash) {
    document.querySelector(route.hash)?.scrollIntoView({ behavior: "smooth" });
  }
});
</script>
