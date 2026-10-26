<script setup>
// core dependencies
import { useNuxtApp } from "nuxt/app";
import { usePush } from "notivue";
import { useMusicStore } from "~~/store/music";
import { useInvitationCodeStore } from "~~/store/invitationcode";
import usecopyToClipboard from "~~/composables/copy_to_clipboard";
import {
  Check,
  Copy,
  Info,
  Link2,
  SkipForward,
  Trophy,
  Volume2,
  VolumeX,
} from "lucide-vue-next";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";

const musicStore = useMusicStore();
const { getMusic, setMusic } = musicStore;

const music = computed(() => {
  return getMusic();
});
const toggleMusic = () => setMusic(!music.value);

// define nuxt configs
const app = useNuxtApp();
const toast = usePush();

// define props and emits
const props = defineProps({
  data: {
    default: () => {
      return {};
    },
    type: Object,
    required: true,
  },
  isAdmin: {
    default: false,
    type: Boolean,
    required: false,
  },
  quizTitle: {
    default: "",
    type: String,
    required: false,
  },
});
const emits = defineEmits(["sendAnswer", "askSkip", "timeUp"]);

const answerable = computed(() => !props.isAdmin);

// The host keeps the join code on screen during the game so latecomers can
// still get in.
const { baseUrl } = useRuntimeConfig().public;
const invitationCodeStore = useInvitationCodeStore();
const joinCode = computed(() =>
  invitationCodeStore.invitationCode
    ? String(invitationCodeStore.invitationCode)
    : ""
);
const joinUrlDisplay = computed(
  () => `${baseUrl.replace(/^https?:\/\//, "")}/join`
);
const joinUrlWithCode = computed(
  () => `${baseUrl}/join?code=${joinCode.value}`
);

const clockOffset = ref(0); // Calculated once, used for all timers
const isOffsetCalculated = ref(false);

// Calculate offset only ONCE when we first receive server_time
const calculateClockOffset = (serverTimeString) => {
  if (isOffsetCalculated.value || !serverTimeString) return;

  try {
    const clientReceiveTime = Date.now();
    const serverTime = new Date(serverTimeString).getTime();
    clockOffset.value = serverTime - clientReceiveTime;
    isOffsetCalculated.value = true;
  } catch (error) {
    console.error("Error calculating clock offset:", error);
  }
};

// Get current time corrected with offset
const getCorrectedNow = () => {
  return new Date(Date.now() + clockOffset.value);
};

// custom refs
const question = ref();
const selectedKey = ref(null);
const counter = ref(null);
const count = ref(0);
const timer = ref(null);
const time = ref(0);
const questionStartTime = ref(null);
const questionDuration = ref(0);
const isSubmitted = ref(false);
const skipCooling = ref(false);
let skipCoolTimer = null;
const SKIP_COOLDOWN_MS = 2000;

const remainingSeconds = computed(() =>
  Math.max(0, (questionDuration.value || 0) - time.value)
);

const questionPoints = computed(() => {
  const points = Number(question.value?.points);
  return Number.isFinite(points) && points >= 0 ? points : null;
});

const pointsLabel = computed(() => {
  if (questionPoints.value === null) return null;
  if (questionPoints.value === 0) return "Jen bonus za rychlost";
  return `Za ${czCount(questionPoints.value, "bod", "body", "bodů")}`;
});

// Determine if current question is the last one
const isLastQuestion = computed(() => {
  if (!question?.value) return false;
  return Number(question.value.no) === Number(question.value.totalQuestions);
});

const countdownDisplay = computed(() => {
  if (typeof count.value === "number" && count.value > 0)
    return String(count.value);
  if (typeof count.value === "string") return "Start!";
  return "Připrav se";
});

const isCountdownNumber = computed(
  () => typeof count.value === "number" && count.value > 0
);

// watchers
watch(
  () => props.data,
  (message) => {
    if (message.status == app.$Fail) {
      toast.error(message.data);
      return;
    }
    handleEvent(message);
  },
  { deep: true, immediate: true }
);

// main function
function handleEvent(message) {
  let counterSound = null;
  // Calculate offset ONCE if server sends time
  if (message.data?.server_time) {
    calculateClockOffset(message.data.server_time);
  }

  if (music.value) {
    counterSound = new Audio("/music/clock.mp3");
  }

  if (message.event == app.$GetQuestion) {
    question.value = message.data;
    questionStartTime.value = new Date(message.data.start_time);
    questionDuration.value = Number(message.data.duration);

    time.value = 0;
    count.value = null;
    selectedKey.value = null;
    isSubmitted.value = false;
    clearTimeout(skipCoolTimer);
    skipCooling.value = false;

    handleTimer();
  } else if (message.event == app.$Counter) {
    question.value = null;
    count.value = parseInt(props.data.data.count);
    time.value = 0;
    clearTimeout(skipCoolTimer);
    skipCooling.value = false;
    handleCounter(counterSound);
    if (music.value && counterSound) {
      counterSound.play();
    }
  }
}

function handleTimer() {
  clearInterval(timer.value);

  if (!questionStartTime.value || !questionDuration.value) {
    console.error("Missing question start time or duration");
    return;
  }

  const duration = questionDuration.value;

  timer.value = setInterval(() => {
    const correctedNow = getCorrectedNow();
    const elapsedSeconds = Math.floor(
      (correctedNow.getTime() - questionStartTime.value.getTime()) / 1000
    );

    time.value = Math.min(Math.max(elapsedSeconds, 0), duration);

    if (elapsedSeconds >= duration) {
      clearInterval(timer.value);
      timer.value = null;
      time.value = duration;
      if (answerable.value && !isSubmitted.value) emits("timeUp");
    }
  }, 100);
}

function handleCounter(counterSound) {
  clearInterval(counter.value);
  counter.value = setInterval(() => {
    count.value -= 1;
    if (count.value <= 0) {
      clearInterval(counter.value);
      count.value = app.$ReadyMessage;
      counter.value = null;
      if (counterSound) {
        counterSound.pause();
      }
    }
  }, 1000);
}

function handleOptionClick(key) {
  if (!answerable.value || isSubmitted.value) return;
  if (props.data?.data?.options === undefined) {
    toast.warning(app.$NoAnswerFound);
    return;
  }
  selectedKey.value = Number(key);
  emits("sendAnswer", [Number(key)]);
  isSubmitted.value = true;
}

function handleSkip(e) {
  e.preventDefault();
  if (skipCooling.value) return;
  skipCooling.value = true;
  clearTimeout(skipCoolTimer);
  skipCoolTimer = setTimeout(() => {
    skipCooling.value = false;
  }, SKIP_COOLDOWN_MS);
  emits("askSkip");
}

// Cleanup on unmount
onUnmounted(() => {
  if (timer.value) {
    clearInterval(timer.value);
  }
  if (counter.value) {
    clearInterval(counter.value);
  }
  clearTimeout(skipCoolTimer);
});

// One palette colour per answer slot (A–E), cycling if there are more.
const OPTION_TONES = [
  {
    badge: "bg-jv-primary text-white",
    border: "border-jv-primary/25 hover:border-jv-primary",
    fill: "bg-jv-lavender",
    ring: "ring-jv-primary",
  },
  {
    badge: "bg-jv-sky text-jv-ink",
    border: "border-jv-sky/35 hover:border-jv-sky",
    fill: "bg-jv-mint",
    ring: "ring-jv-sky",
  },
  {
    badge: "bg-jv-violet text-jv-ink",
    border: "border-jv-violet/40 hover:border-jv-violet",
    fill: "bg-jv-lavender",
    ring: "ring-jv-violet",
  },
  {
    badge: "bg-jv-orange text-jv-ink",
    border: "border-jv-orange/40 hover:border-jv-orange",
    fill: "bg-jv-yellow-soft",
    ring: "ring-jv-orange",
  },
  {
    badge: "bg-jv-mauve text-jv-ink",
    border: "border-jv-mauve/45 hover:border-jv-mauve",
    fill: "bg-jv-lavender",
    ring: "ring-jv-mauve",
  },
];
const optionTone = (key) =>
  OPTION_TONES[(Number(key) - 1 + OPTION_TONES.length) % OPTION_TONES.length];
</script>

<template>
  <!-- Countdown view -->
  <main
    v-if="question == null"
    class="flex min-h-[75vh] flex-col items-center justify-center bg-jv-canvas px-4 py-10 text-jv-ink sm:py-14"
  >
    <article
      class="relative w-full max-w-[460px] jv-border-rough bg-jv-white px-6 py-9 shadow-soft sm:px-10 sm:py-12"
    >
      <p
        class="text-center font-body text-[12px] font-black uppercase tracking-[0.22em] text-jv-muted sm:text-[13px]"
      >
        Připrav se
      </p>

      <transition name="pop" mode="out-in">
        <div
          v-if="isCountdownNumber"
          :key="String(count)"
          class="mx-auto mt-5 grid size-[180px] place-items-center rounded-full border border-dashed border-jv-line bg-jv-yellow font-headings text-[110px] leading-none text-jv-ink sm:size-[220px] sm:text-[140px]"
        >
          {{ countdownDisplay }}
        </div>
        <h1
          v-else
          :key="'msg'"
          class="mt-6 text-center font-headings text-[48px] leading-none text-jv-ink sm:text-[64px]"
        >
          {{ countdownDisplay }}
        </h1>
      </transition>

      <div class="mt-6 flex justify-center gap-2">
        <span class="size-2 rounded-full bg-jv-coral"></span>
        <span class="size-2 rounded-full bg-jv-yellow"></span>
        <span class="size-2 rounded-full bg-jv-mint"></span>
      </div>
    </article>
  </main>

  <!-- Question view -->
  <main
    v-else
    :class="[
      'flex flex-col bg-jv-canvas px-3 py-4 text-jv-ink sm:px-6 sm:py-6 md:px-10 md:py-8',
      isAdmin ? '' : 'min-h-screen',
    ]"
  >
    <section class="mx-auto w-full max-w-[1180px]">
      <article
        class="relative jv-border-rough bg-jv-white p-5 shadow-soft sm:p-7 md:p-9"
      >
        <!-- Header: question count + timer -->
        <header class="flex items-start justify-between gap-4">
          <div class="min-w-0 flex-1">
            <p
              v-if="quizTitle"
              class="mb-2 inline-flex max-w-full items-center gap-2 rounded-full border border-jv-line bg-jv-yellow px-3 py-1 font-body text-[11px] font-black uppercase tracking-[0.12em] text-jv-ink shadow-soft-sm sm:text-[12px]"
              :title="quizTitle"
            >
              <span
                class="size-1.5 rounded-full bg-jv-coral"
                aria-hidden="true"
              ></span>
              <span class="truncate">{{ quizTitle }}</span>
            </p>
            <h2
              class="font-headings text-[24px] leading-none text-jv-ink min-[420px]:text-[30px] sm:text-[36px] md:text-[40px]"
            >
              Otázka {{ question.no }}
              <span class="text-jv-muted">/ {{ question.totalQuestions }}</span>
            </h2>
            <p
              v-if="!pointsLabel"
              class="mt-2 font-body text-[12px] font-bold text-jv-muted sm:text-[14px]"
            >
              Hrajeme!
            </p>
            <Popover v-else>
              <PopoverTrigger as-child>
                <button
                  type="button"
                  class="mt-2 inline-flex max-w-full items-center gap-1.5 rounded-full border border-jv-line bg-jv-mint px-3 py-1 font-body text-[11px] font-black uppercase tracking-[0.12em] text-jv-ink shadow-soft-sm transition-transform sm:text-[12px]"
                  aria-label="Jak funguje bodování"
                >
                  <span class="truncate">{{ pointsLabel }}</span>
                  <Info class="size-3 shrink-0" :stroke-width="2.6" />
                </button>
              </PopoverTrigger>
              <PopoverContent
                side="bottom"
                align="start"
                :side-offset="8"
                class="w-72 jv-border-rough bg-jv-white p-4 font-body text-[13px] text-jv-ink shadow-soft"
              >
                <p
                  class="font-headings text-[15px] uppercase tracking-[0.06em]"
                >
                  Jak se počítá skóre
                </p>
                <ul class="mt-3 space-y-1.5">
                  <li>Správná odpověď: <b>500</b> základních bodů</li>
                  <li>
                    Bonus za rychlost: až <b>400</b> podle toho, kolik času
                    zbývá
                  </li>
                  <li>
                    Body za otázku: <b>{{ questionPoints }} &times; 100</b>
                  </li>
                  <li>
                    Bonus za sérii: od 2. správné odpovědi v řadě
                    <b>+100</b> a dál roste
                  </li>
                </ul>
                <p class="mt-3 text-[12px] font-bold text-jv-muted">
                  Špatná nebo chybějící odpověď dává 0 bodů a nuluje sérii.
                </p>
              </PopoverContent>
            </Popover>
          </div>

          <div class="flex shrink-0 items-center gap-2 sm:gap-3">
            <button
              type="button"
              class="grid size-10 place-items-center rounded-[8px] border border-jv-line bg-jv-white text-jv-ink shadow-soft-sm transition-transform sm:size-11"
              :aria-label="music ? 'Vypnout hudbu' : 'Zapnout hudbu'"
              @click="toggleMusic"
            >
              <Volume2 v-if="music" class="size-4" :stroke-width="2.4" />
              <VolumeX v-else class="size-4" :stroke-width="2.4" />
            </button>

            <div
              class="grid size-12 place-items-center rounded-full border border-dashed border-jv-line bg-jv-yellow font-feature text-[20px] font-black text-jv-ink sm:size-[60px] sm:text-[24px]"
              :aria-label="`Zbývá ${czCount(
                remainingSeconds,
                'sekunda',
                'sekundy',
                'sekund'
              )}`"
            >
              {{ remainingSeconds }}
            </div>
          </div>
        </header>

        <div
          v-if="isAdmin && joinCode"
          class="mt-4 flex flex-wrap items-center gap-2 rounded-2xl bg-jv-lavender px-3 py-2 text-[13px] text-jv-ink sm:gap-3 sm:px-4 sm:text-[14px]"
        >
          <span class="text-jv-muted">
            Připoj se na
            <strong class="text-jv-ink">{{ joinUrlDisplay }}</strong>
          </span>
          <span class="text-jv-muted">· Kód</span>
          <strong
            class="font-feature text-[18px] font-extrabold tracking-[0.12em] sm:text-[20px]"
            >{{ joinCode }}</strong
          >
          <span class="ml-auto flex gap-2">
            <button
              type="button"
              class="inline-flex h-8 items-center gap-1.5 rounded-full bg-jv-white px-3 font-semibold shadow-soft-sm transition hover:shadow-soft"
              @click="usecopyToClipboard(joinCode)"
            >
              <Copy class="size-3.5" :stroke-width="2.4" />
              Kopírovat kód
            </button>
            <button
              type="button"
              class="inline-flex h-8 items-center gap-1.5 rounded-full bg-jv-white px-3 font-semibold shadow-soft-sm transition hover:shadow-soft"
              @click="usecopyToClipboard(joinUrlWithCode)"
            >
              <Link2 class="size-3.5" :stroke-width="2.4" />
              Kopírovat odkaz
            </button>
          </span>
        </div>

        <div class="my-5 border-t border-dashed border-jv-line sm:my-6"></div>

        <!-- Question text -->
        <h3
          class="font-headings text-[20px] leading-snug text-jv-ink min-[420px]:text-[24px] sm:text-[30px] md:text-[36px]"
        >
          {{ question.question }}
        </h3>

        <div
          v-if="question.question_media === 'image'"
          class="mt-5 flex justify-center"
        >
          <NuxtImg
            :src="question.resource"
            :alt="question.question"
            class="max-h-[260px] w-auto border border-jv-line bg-jv-white object-contain shadow-soft-sm"
          />
        </div>
        <div v-else-if="question.question_media === 'code'" class="mt-5">
          <CodeBlockComponent :code="question.resource" />
        </div>

        <!-- Options grid -->
        <div class="mt-6 grid gap-3 sm:mt-7 sm:gap-4 md:grid-cols-2 md:gap-5">
          <button
            v-for="(value, key) in question.options"
            :key="key"
            type="button"
            :disabled="!answerable || isSubmitted"
            :class="[
              'group relative flex w-full items-center gap-3 overflow-hidden rounded-2xl border-2 bg-jv-white px-3 py-3 text-left transition-all duration-150 sm:gap-4 sm:px-4 sm:py-4',
              optionTone(key).border,
              answerable && !isSubmitted
                ? 'cursor-pointer shadow-soft-sm hover:-translate-y-0.5 hover:shadow-soft active:scale-[0.99]'
                : 'shadow-soft-sm',
              isSubmitted && selectedKey === Number(key)
                ? [
                    optionTone(key).fill,
                    'ring-4 ring-offset-2',
                    optionTone(key).ring,
                  ]
                : '',
              isSubmitted && selectedKey !== Number(key) ? 'opacity-45' : '',
              !answerable ? 'cursor-default' : '',
            ]"
            @click="handleOptionClick(key)"
          >
            <span
              :class="[
                'grid size-10 shrink-0 place-items-center rounded-xl font-feature text-[15px] font-extrabold sm:size-12 sm:text-[17px]',
                optionTone(key).badge,
              ]"
            >
              {{ String.fromCharCode(64 + Number(key)) }}
            </span>
            <div class="min-w-0 flex-1">
              <div
                v-if="question.options_media === 'image'"
                class="flex items-center justify-center"
              >
                <NuxtImg
                  :src="value"
                  :alt="`Možnost ${key}`"
                  class="max-h-[120px] w-auto object-contain"
                />
              </div>
              <CodeBlockComponent
                v-else-if="question.options_media === 'code'"
                :code="value"
              />
              <span
                v-else
                class="block break-words font-body text-[15px] font-bold leading-snug text-jv-ink sm:text-[18px] md:text-[20px]"
              >
                {{ value }}
              </span>
            </div>
            <Check
              v-if="isSubmitted && selectedKey === Number(key)"
              class="size-6 shrink-0 text-jv-ink"
              :stroke-width="3"
            />
          </button>
        </div>

        <!-- Admin footer: waiting status + skip -->
        <template v-if="isAdmin">
          <div class="mt-6 border-t border-dashed border-jv-line sm:mt-7"></div>
          <footer
            class="mt-4 flex flex-col items-stretch justify-between gap-3 sm:flex-row sm:items-center sm:gap-4"
          >
            <div class="flex items-center gap-3">
              <div class="flex gap-1">
                <span
                  class="size-2 rounded-full bg-jv-coral sm:size-2.5"
                ></span>
                <span
                  class="size-2 rounded-full bg-jv-yellow sm:size-2.5"
                ></span>
                <span class="size-2 rounded-full bg-jv-mint sm:size-2.5"></span>
              </div>
              <p
                class="font-body text-[13px] font-bold text-jv-muted sm:text-[14px]"
              >
                Čekáme, až hráči odpoví
              </p>
            </div>
            <button
              type="button"
              class="inline-flex h-10 items-center justify-center gap-2 self-end rounded-full border border-jv-line bg-jv-white px-5 font-body text-[14px] font-black text-jv-ink shadow-soft-sm transition-transform disabled:cursor-not-allowed disabled:opacity-50 sm:h-11 sm:self-auto sm:px-6 sm:text-[15px]"
              :disabled="skipCooling"
              @click="handleSkip"
            >
              <Trophy
                v-if="isLastQuestion"
                class="size-4"
                :stroke-width="2.4"
              />
              <SkipForward v-else class="size-4" :stroke-width="2.4" />
              <span>{{ isLastQuestion ? "Ukončit kvíz" : "Přeskočit" }}</span>
            </button>
          </footer>
        </template>

        <!-- User submission acknowledgement -->
        <div
          v-if="answerable && isSubmitted"
          class="mt-6 flex items-center justify-center gap-2 rounded-full border border-dashed border-jv-line px-4 py-2 font-body text-[13px] font-bold text-jv-muted sm:text-[14px]"
        >
          <Check class="size-4" :stroke-width="2.6" />
          Odpověď odeslána – čekej na další otázku
        </div>
      </article>
    </section>
  </main>
</template>

<style scoped>
.pop-enter-from {
  opacity: 0;
  transform: scale(0.4) rotate(-12deg);
}
.pop-leave-to {
  opacity: 0;
  transform: scale(1.3) rotate(8deg);
}
.pop-enter-active,
.pop-leave-active {
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
</style>
