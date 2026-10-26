<script setup>
import {
  ArrowRight,
  Check,
  Flame,
  Hourglass,
  MessageSquare,
  Trophy,
  X,
} from "lucide-vue-next";
import { czCount } from "~~/composables/plural";
import { FEEDBACK_AUTO_ADVANCE_SECONDS } from "~~/composables/self_paced";

const props = defineProps({
  // Shape produced by buildFeedback() in composables/self_paced.js.
  feedback: {
    type: Object,
    required: true,
  },
  pending: {
    default: false,
    type: Boolean,
    required: false,
  },
});
const emits = defineEmits(["next"]);

const OUTCOMES = {
  correct: {
    title: "Správně!",
    icon: Check,
    badge: "bg-jv-mint",
    band: "bg-jv-mint",
  },
  wrong: {
    title: "Špatně",
    icon: X,
    badge: "bg-jv-coral text-white",
    band: "bg-jv-lavender",
  },
  timeout: {
    title: "Čas vypršel",
    icon: Hourglass,
    badge: "bg-jv-orange",
    band: "bg-jv-yellow-soft",
  },
  survey: {
    title: "Odpověď uložena",
    icon: MessageSquare,
    badge: "bg-jv-sky text-white",
    band: "bg-jv-mint",
  },
};

const outcome = computed(
  () => OUTCOMES[props.feedback.outcome] || OUTCOMES.wrong
);

const nextLabel = computed(() =>
  props.feedback.finished ? "Zobrazit výsledky" : "Další otázka"
);

const correctLabel = computed(() =>
  props.feedback.correctOptions.length > 1
    ? "Správné odpovědi"
    : "Správná odpověď"
);

const secondsLeft = ref(FEEDBACK_AUTO_ADVANCE_SECONDS);
let countdown = null;

const goNext = () => {
  clearInterval(countdown);
  countdown = null;
  emits("next");
};

onMounted(() => {
  countdown = setInterval(() => {
    secondsLeft.value -= 1;
    if (secondsLeft.value <= 0) goNext();
  }, 1000);
});

onUnmounted(() => {
  clearInterval(countdown);
});
</script>

<template>
  <main
    class="flex min-h-[75vh] flex-col items-center justify-center bg-jv-canvas px-3 py-6 text-jv-ink sm:px-6 sm:py-10"
  >
    <article
      class="relative w-full max-w-[560px] overflow-hidden jv-border-rough bg-jv-white shadow-soft"
    >
      <header
        :class="['px-5 pb-6 pt-6 text-center sm:px-8 sm:pt-8', outcome.band]"
      >
        <p
          class="font-body text-[12px] font-black uppercase tracking-[0.18em] text-jv-muted sm:text-[13px]"
        >
          Otázka {{ feedback.answered }} z {{ feedback.totalQuestions }}
        </p>
        <span
          :class="[
            'mx-auto mt-4 grid size-16 place-items-center rounded-full border border-jv-line shadow-soft-sm sm:size-20',
            outcome.badge,
          ]"
          aria-hidden="true"
        >
          <component
            :is="outcome.icon"
            class="size-8 sm:size-10"
            :stroke-width="3"
          />
        </span>
        <h1
          class="mt-4 font-headings text-[36px] leading-none text-jv-ink sm:text-[48px]"
        >
          {{ outcome.title }}
        </h1>
        <p
          class="mt-3 font-feature text-[22px] font-black text-jv-ink sm:text-[26px]"
        >
          +{{ czCount(feedback.score, "bod", "body", "bodů") }}
        </p>
      </header>

      <div class="px-5 py-5 sm:px-8 sm:py-6">
        <div class="grid grid-cols-2 gap-3">
          <div
            class="flex items-center gap-3 rounded-2xl border border-jv-line bg-jv-white px-3 py-3 shadow-soft-sm"
          >
            <span
              class="grid size-10 shrink-0 place-items-center rounded-xl bg-jv-yellow"
            >
              <Trophy class="size-5" :stroke-width="2.4" />
            </span>
            <div class="min-w-0">
              <p
                class="font-body text-[11px] font-black uppercase tracking-[0.12em] text-jv-muted"
              >
                Celkem
              </p>
              <p class="font-feature text-[20px] font-black leading-tight">
                {{ feedback.totalScore }}
              </p>
            </div>
          </div>
          <div
            class="flex items-center gap-3 rounded-2xl border border-jv-line bg-jv-white px-3 py-3 shadow-soft-sm"
          >
            <span
              class="grid size-10 shrink-0 place-items-center rounded-xl bg-jv-orange"
            >
              <Flame class="size-5" :stroke-width="2.4" />
            </span>
            <div class="min-w-0">
              <p
                class="font-body text-[11px] font-black uppercase tracking-[0.12em] text-jv-muted"
              >
                Série
              </p>
              <p class="font-feature text-[20px] font-black leading-tight">
                {{ feedback.streak }}
              </p>
            </div>
          </div>
        </div>

        <section
          v-if="
            feedback.outcome !== 'correct' && feedback.correctOptions.length
          "
          class="mt-4 rounded-2xl bg-jv-lavender px-4 py-3"
        >
          <p
            class="font-body text-[11px] font-black uppercase tracking-[0.12em] text-jv-muted"
          >
            {{ correctLabel }}
          </p>
          <ul class="mt-2 space-y-2">
            <li
              v-for="option in feedback.correctOptions"
              :key="option.key"
              class="flex items-center gap-3"
            >
              <span
                class="grid size-8 shrink-0 place-items-center rounded-lg bg-jv-primary font-feature text-[14px] font-extrabold text-white"
              >
                {{ option.letter }}
              </span>
              <NuxtImg
                v-if="feedback.optionsMedia === 'image'"
                :src="option.value"
                :alt="`Možnost ${option.letter}`"
                class="max-h-[80px] w-auto object-contain"
              />
              <CodeBlockComponent
                v-else-if="feedback.optionsMedia === 'code'"
                :code="option.value"
                class="min-w-0 flex-1"
              />
              <span
                v-else
                class="min-w-0 break-words font-body text-[15px] font-bold leading-snug sm:text-[17px]"
              >
                {{ option.value }}
              </span>
            </li>
          </ul>
        </section>

        <button
          type="button"
          class="mt-5 flex h-[52px] w-full items-center justify-center gap-2 rounded-full border border-jv-line bg-jv-mint px-5 font-headings text-[16px] text-jv-ink shadow-soft-sm transition-transform hover:-translate-y-px disabled:cursor-not-allowed disabled:opacity-60 sm:h-14 sm:text-[18px]"
          :disabled="pending"
          @click="goNext"
        >
          <span>{{ pending ? "Načítám…" : nextLabel }}</span>
          <ArrowRight v-if="!pending" class="size-5" :stroke-width="2.6" />
        </button>
        <p
          v-if="!pending && secondsLeft > 0"
          class="mt-3 text-center font-body text-[12px] font-bold text-jv-muted sm:text-[13px]"
        >
          Pokračujeme automaticky za
          {{ czCount(secondsLeft, "sekundu", "sekundy", "sekund") }}
        </p>
      </div>
    </article>
  </main>
</template>
