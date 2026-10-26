<script setup>
import { Home, Library, AlertCircle } from "lucide-vue-next";

definePageMeta({
  layout: "empty",
  hideSidebar: true,
});

const route = useRoute();

// Admins land here from /admin/arrange/[session_id] when the websocket reports
// that the session was completed, the invitation code was missing, or session
// validation failed. The redirect tacks on ?status=...&error=... so we can show
// a context-specific message rather than a generic "temporary redirection".
const presentation = computed(() => {
  const err = String(route.query.error || "").toLowerCase();

  if (err === "session was completed") {
    return {
      title: "Kvíz skončil",
      message:
        "Tohle sezení kvízu už skončilo. Po skončení kvízu se k ovládání moderátora vrátit nedá – na úvodní stránce můžeš spustit nový kvíz nebo si projít výsledky v přehledu.",
    };
  }
  if (err === "quiz-session-validation-failed") {
    return {
      title: "Sezení kvízu už není aktivní",
      message:
        "Tohle sezení kvízu už nejde moderovat. Mohlo skončit, být ukončeno nebo vypršet. Pro nové hraní spusť nové sezení.",
    };
  }
  if (err === "invitation code not found") {
    return {
      title: "Kód kvízu už neplatí",
      message:
        "Kód kvízu pro toto sezení už není aktivní. Sezení nejspíš skončilo.",
    };
  }
  return {
    title: "Žádné aktivní sezení kvízu",
    message:
      "Teď tu není co moderovat. Vrať se na úvodní stránku nebo otevři své kvízy a spusť nové sezení.",
  };
});

useHead({
  title: "Sezení kvízu skončilo - Admin",
  meta: [
    {
      name: "description",
      content:
        "Sezení kvízu, které jsi moderoval/a, skončilo. Vrať se na úvodní stránku nebo do přehledu a spusť nový kvíz.",
    },
  ],
});

useSeoMeta({
  title: "Příprava kvízu – STEM_HUB",
  description:
    "Připrav si sezení kvízu, než ho v systému STEM_HUB spustíš naživo.",
  robots: "noindex, nofollow",
});
</script>

<template>
  <div
    class="relative min-h-screen overflow-hidden bg-jv-canvas px-4 py-10 sm:px-6 sm:py-14"
  >
    <div
      class="pointer-events-none absolute inset-0 jv-grid"
      aria-hidden="true"
    ></div>

    <div
      class="relative mx-auto flex w-full max-w-[640px] flex-col items-center"
    >
      <div
        class="relative w-full jv-border-rough bg-jv-white p-6 shadow-soft-lg sm:p-9"
      >
        <div class="flex items-start gap-4">
          <span
            class="grid size-12 shrink-0 place-items-center rounded-[10px] border border-jv-line bg-jv-yellow text-jv-ink shadow-soft-sm sm:size-14"
            aria-hidden="true"
          >
            <AlertCircle class="size-6 sm:size-7" :stroke-width="2.4" />
          </span>
          <div class="min-w-0 flex-1">
            <p
              class="font-body text-[11px] font-black uppercase tracking-[0.22em] text-jv-muted sm:text-[12px]"
            >
              Admin
            </p>
            <h1
              class="mt-1 font-headings text-[26px] leading-tight text-jv-ink sm:text-[32px]"
            >
              {{ presentation.title }}
            </h1>
          </div>
        </div>

        <hr class="my-5 border-t border-dashed border-jv-line" />

        <p
          class="font-body text-[14px] font-bold leading-relaxed text-jv-muted sm:text-[15px]"
        >
          {{ presentation.message }}
        </p>

        <div
          class="mt-6 flex flex-col items-stretch gap-3 sm:flex-row sm:justify-center"
          role="group"
          aria-label="Navigace zpět"
        >
          <NuxtLink
            to="/"
            class="inline-flex h-11 items-center justify-center gap-2 rounded-full border border-jv-line bg-jv-coral px-6 font-body text-[14px] font-black text-white shadow-soft-sm transition-transform sm:h-12 sm:text-[15px]"
            aria-label="Zpět na úvodní stránku"
          >
            <Home class="size-4" :stroke-width="2.4" />
            Zpět domů
          </NuxtLink>
          <NuxtLink
            to="/admin/quiz/list-quiz"
            class="inline-flex h-11 items-center justify-center gap-2 rounded-full border border-jv-line bg-jv-white px-6 font-body text-[14px] font-black text-jv-ink shadow-soft-sm transition-transform sm:h-12 sm:text-[15px]"
            aria-label="Přejít na moje kvízy"
          >
            <Library class="size-4" :stroke-width="2.4" />
            Moje kvízy
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.jv-grid {
  background-image: radial-gradient(var(--jv-charcoal) 1px, transparent 1px);
  background-size: 20px 20px;
  opacity: 0.08;
}
</style>
