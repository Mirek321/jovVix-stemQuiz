<script setup>
// core dependencies
import { useNuxtApp, useRouter } from "nuxt/app";
import { usePush } from "notivue";
import { useMusicStore } from "~~/store/music";
import {
  Check,
  Copy,
  Home,
  Hourglass,
  Info,
  Keyboard,
  LogOut,
  Maximize2,
  Smile,
  UserRound,
  Users,
  Volume2,
  VolumeX,
} from "lucide-vue-next";
const { baseUrl } = useRuntimeConfig().public;
const musicStore = useMusicStore();
const { getMusic, setMusic } = musicStore;

const music = computed(() => {
  return getMusic();
});

// define nuxt configs
const app = useNuxtApp();
const toast = usePush();
const router = useRouter();

import { useInvitationCodeStore } from "~/store/invitationcode.js";
import { useListUserstore } from "~/store/userlist";
import { useUsersStore } from "~~/store/users";
import { storeToRefs } from "pinia";
import usecopyToClipboard from "~~/composables/copy_to_clipboard";
import { getAvatarUrlByName } from "~~/composables/avatar";

const invitationCodeStore = useInvitationCodeStore();
const { invitationCode } = storeToRefs(invitationCodeStore);

const listUserStore = useListUserstore();
const { removeAllUsers } = listUserStore;
const { listUsers } = storeToRefs(listUserStore);

const usersStore = useUsersStore();
const { getUserData } = usersStore;

const startQuiz = ref(false);
const waitingSound = ref(null);
const qrOverlay = ref(null);
const participantAccentClasses = [
  "bg-jv-yellow",
  "bg-jv-coral text-white",
  "bg-jv-mint",
  "bg-jv-white",
];

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
  userName: {
    default: "",
    type: String,
    required: false,
  },
  quizTitleOverride: {
    default: "",
    type: String,
    required: false,
  },
});
const emits = defineEmits(["startQuiz", "terminateQuiz"]);

// custom refs
const code = ref(invitationCode.value);
const participantCount = computed(() => listUsers.value.length);
const joinUrl = computed(() => `${baseUrl}/join`);
const joinUrlWithCode = computed(() => `${joinUrl.value}?code=${code.value}`);

const getParticipantName = (user) =>
  user?.UserName || user?.username || user?.name || "Hráč";

const quizTitle = computed(
  () =>
    props.quizTitleOverride ||
    props.data?.quizTitle ||
    props.data?.title ||
    props.data?.data?.quizTitle ||
    "Živý kvíz"
);

const waitingMessage = computed(() => {
  const raw = props.data?.data;
  if (typeof raw === "string" && raw.length) return raw;
  return "Kvíz brzy začne…";
});

const toggleMusic = () => setMusic(!music.value);

const leaveLobby = () => {
  router.push("/join");
};

const displayName = computed(() => {
  const name = (props.userName || "").trim();
  return name || "Hráč";
});

const userAvatar = computed(() => {
  const stored = getUserData();
  return getAvatarUrlByName(stored?.avatar || displayName.value);
});

// watchers
watch(
  () => props.data,
  (message) => {
    if (message.status == app.$Fail) {
      toast.error(message.data);
    }
    handleEvent(message);
  },
  { deep: true, immediate: true }
);

// event handlers
function start_quiz(e) {
  e.preventDefault();
  startQuiz.value = true;
  emits("startQuiz");
}

// main function
function handleEvent(message) {
  if (message.event == app.$SentInvitaionCode) {
    code.value = invitationCode.value.toString();
  }
}

onMounted(() => {
  initializeSound();
});

const copyToClipBoard = (text) => {
  usecopyToClipboard(text);
};

const openQrFullscreen = () => qrOverlay.value?.open();

function initializeSound() {
  if (process.client) {
    waitingSound.value = new Audio("/music/waiting_area_music.mp3");
    if (music.value) {
      waitingSound.value.play();
      waitingSound.value.loop = true;
    }
  }
}

onUnmounted(() => {
  if (waitingSound.value) {
    waitingSound.value.pause();
    waitingSound.value = null;
  }
  if (!startQuiz.value) {
    emits("terminateQuiz");
  }
  if (props.isAdmin) {
    removeAllUsers();
  }
});

watch(
  music,
  (newValue) => {
    if (!newValue && waitingSound.value) {
      waitingSound.value.pause();
    } else if (newValue && waitingSound.value) {
      waitingSound.value.play();
    }
  },
  { deep: true, immediate: true }
);
</script>

<template>
  <main
    v-if="isAdmin"
    class="flex min-h-screen flex-col gap-8 bg-jv-canvas px-3 py-5 text-jv-ink sm:gap-10 sm:px-6 md:px-8 md:py-6"
  >
    <section
      class="relative mx-auto flex w-full max-w-[1220px] flex-1 flex-col"
    >
      <header class="relative mb-7 text-center sm:mb-9">
        <div
          class="mb-3 flex items-center justify-between gap-2 sm:mb-2 sm:grid sm:grid-cols-[1fr_auto_1fr]"
        >
          <div class="hidden sm:block" aria-hidden="true"></div>

          <p
            class="font-body text-[11px] font-black uppercase tracking-[0.18em] text-jv-muted sm:text-[13px] sm:justify-self-center"
          >
            Čekárna
          </p>

          <div class="flex shrink-0 justify-end sm:justify-self-end">
            <NuxtLink
              to="/"
              class="inline-flex h-9 items-center justify-center gap-1.5 rounded-[8px] border border-jv-line bg-jv-white px-3 font-body text-[13px] font-bold text-jv-ink shadow-soft-sm transition-transform sm:h-12 sm:gap-2 sm:px-5 sm:text-[18px]"
              aria-label="Zpět domů"
            >
              <Home class="size-4 sm:size-5" :stroke-width="2.4" />
              <span>Domů</span>
            </NuxtLink>
          </div>
        </div>

        <h1
          class="relative mx-auto mt-2 inline-block max-w-[940px] break-words font-headings text-[30px] leading-[1.05] text-jv-ink min-[420px]:text-[38px] sm:text-[48px] md:text-[56px]"
        >
          <span class="relative z-10">{{ quizTitle }}</span>
          <span
            class="absolute bottom-[4px] left-0 z-0 h-[10px] w-full max-w-[260px] bg-jv-yellow/70"
            aria-hidden="true"
          ></span>
        </h1>
        <p
          class="mx-auto mt-3 max-w-[680px] text-[15px] leading-[1.5] text-jv-muted sm:text-[18px] md:text-[20px]"
        >
          Pošli kód hráčům. Kvíz spustíš, až budeš připravený!
        </p>
      </header>

      <div
        class="grid flex-1 gap-6 md:gap-7 xl:grid-cols-[minmax(0,1fr)_minmax(340px,0.95fr)] xl:items-stretch"
      >
        <form
          class="relative flex min-h-0 flex-col overflow-hidden bg-jv-white shadow-soft-sm jv-border-rough sm:shadow-soft-lg xl:min-h-[640px]"
          @submit="start_quiz"
        >
          <span
            class="absolute left-1/2 top-[-2px] z-10 h-4 w-12 -translate-x-1/2 bg-jv-coral"
            aria-hidden="true"
          ></span>

          <div class="bg-jv-yellow px-4 pb-7 pt-6 sm:px-8 sm:pb-9 md:px-10">
            <h2
              class="join-page-title font-headings text-[30px] leading-tight text-jv-ink min-[420px]:text-[34px] sm:text-[42px] lg:text-[48px]"
            >
              Připravit, pozor, start!
            </h2>

            <div
              class="mt-4 flex flex-col gap-3 jv-border-rough bg-jv-white p-3 shadow-soft-sm sm:flex-row sm:items-center sm:justify-between sm:p-4"
            >
              <div
                class="flex min-w-0 flex-wrap items-baseline gap-x-2 gap-y-1"
              >
                <p
                  class="shrink-0 font-body text-[12px] font-black uppercase tracking-[0.08em] text-jv-muted"
                >
                  Připoj se na
                </p>
                <p
                  class="min-w-0 break-all font-body text-[16px] font-extrabold leading-[1.35] text-jv-ink min-[420px]:text-[17px] sm:text-[20px] md:text-[22px]"
                >
                  {{ joinUrl.replace(/^https?:\/\//, "") }}
                </p>
              </div>
              <button
                id="URL-input-container"
                type="button"
                class="inline-flex h-11 w-full shrink-0 items-center justify-center gap-2 rounded-[999px] border border-jv-line bg-jv-coral px-5 font-body text-[16px] font-bold text-white shadow-soft-sm transition-transform sm:h-12 sm:w-fit"
                @click="copyToClipBoard(joinUrlWithCode)"
              >
                <Copy class="size-4" :stroke-width="2.5" />
                <span>Kopírovat</span>
              </button>
            </div>
          </div>

          <div class="flex flex-1 flex-col px-4 py-5 sm:px-8 sm:py-6 md:px-10">
            <div
              class="-mt-11 flex flex-col gap-3 jv-border-rough bg-jv-white p-3 shadow-soft-sm sm:-mt-14 sm:flex-row sm:items-center sm:justify-between sm:p-4 sm:shadow-soft"
            >
              <div class="flex min-w-0 items-center gap-3 sm:gap-4">
                <span
                  class="grid size-10 shrink-0 place-items-center border border-jv-line bg-jv-mint text-jv-ink sm:size-11"
                >
                  <Keyboard class="size-5" :stroke-width="2.2" />
                </span>
                <div
                  class="flex min-w-0 flex-wrap items-baseline gap-x-2 gap-y-1"
                >
                  <p
                    class="shrink-0 font-body text-[12px] font-black uppercase tracking-[0.08em] text-jv-muted"
                  >
                    Kód kvízu
                  </p>
                  <h3
                    class="code min-w-0 break-all font-feature text-[22px] font-black leading-none text-jv-ink min-[420px]:text-[24px] sm:text-[26px]"
                  >
                    {{ code }}
                  </h3>
                </div>
              </div>
              <button
                id="OTP-input-container"
                type="button"
                class="inline-flex h-11 w-full shrink-0 items-center justify-center gap-2 rounded-[999px] border border-jv-line bg-jv-coral px-5 font-body text-[16px] font-bold text-white shadow-soft-sm transition-transform sm:h-12 sm:w-fit"
                @click="copyToClipBoard(code)"
              >
                <Copy class="size-4" :stroke-width="2.5" />
                <span>Kopírovat</span>
              </button>
            </div>

            <div
              class="my-5 flex items-center justify-center gap-3 text-center font-body text-[11px] font-black uppercase tracking-[0.08em] text-jv-muted min-[420px]:text-[12px] sm:my-6 sm:gap-4 sm:text-[13px]"
            >
              <span
                class="h-px min-w-8 flex-1 border-t border-dashed border-jv-line sm:max-w-24"
              ></span>
              <span class="shrink-0">Nebo naskenuj QR kód</span>
              <span
                class="h-px min-w-8 flex-1 border-t border-dashed border-jv-line sm:max-w-24"
              ></span>
            </div>

            <div class="flex flex-col items-center gap-3">
              <div
                class="qr-card relative grid size-[176px] place-items-center bg-jv-white p-3 shadow-soft-sm jv-border-rough min-[420px]:size-[196px] sm:size-[220px] sm:p-4 md:size-[240px]"
              >
                <QrCode :scan-u-r-l="joinUrl" :quiz-code="code" :size="200" />
              </div>
              <button
                type="button"
                class="inline-flex h-11 items-center justify-center gap-2 rounded-[999px] border border-jv-line bg-jv-mint px-5 font-body text-[14px] font-black text-jv-ink shadow-soft-sm transition-transform sm:h-12 sm:text-[15px]"
                @click="openQrFullscreen"
              >
                <Maximize2 class="size-4" :stroke-width="2.4" />
                <span>Zvětšit QR kód</span>
              </button>
            </div>

            <div class="mt-auto pt-5 sm:pt-6">
              <button
                type="submit"
                class="mx-auto flex h-[52px] w-full max-w-[390px] items-center justify-center rounded-[999px] border border-jv-line bg-jv-mint px-5 font-headings text-[16px] text-jv-ink shadow-soft-sm transition-transform sm:h-16 sm:px-6 sm:text-[18px]"
              >
                Spustit kvíz
              </button>
              <p
                class="mt-4 flex items-center justify-center gap-2 text-center font-body text-[12px] leading-[1.4] text-jv-muted sm:text-[13px]"
              >
                <Info class="size-4" :stroke-width="2.3" />
                <span>Kvíz můžeš spustit kdykoli</span>
              </p>
            </div>
          </div>
        </form>

        <aside
          class="relative flex min-h-0 flex-col bg-jv-white px-4 py-6 shadow-soft-sm jv-border-rough sm:px-8 sm:py-8 sm:shadow-soft-lg md:px-10 xl:min-h-[640px]"
        >
          <div
            class="flex items-start justify-between gap-3 sm:items-center sm:gap-4"
          >
            <div class="flex min-w-0 items-center gap-3 sm:gap-4">
              <span
                class="grid size-10 shrink-0 place-items-center jv-border-rough bg-jv-mint sm:size-12"
              >
                <Users class="size-5 sm:size-6" :stroke-width="2.4" />
              </span>
              <h2
                class="min-w-0 break-words font-headings text-[28px] leading-tight text-jv-ink min-[420px]:text-[32px] sm:text-[42px] xl:text-[48px]"
              >
                Hráči
              </h2>
            </div>
            <span
              class="grid size-11 shrink-0 place-items-center border border-jv-line bg-jv-yellow font-feature text-[24px] font-black shadow-soft-sm sm:size-14 sm:text-[30px]"
            >
              {{ participantCount }}
            </span>
          </div>

          <div class="my-5 border-t border-dashed border-jv-line sm:my-6"></div>

          <div
            v-if="participantCount"
            class="flex flex-col gap-3 overflow-y-auto pr-1 sm:gap-4 xl:max-h-[430px]"
          >
            <div
              v-for="(user, index) in listUsers"
              :key="user.UserId || user.UserName || index"
              class="flex min-h-14 items-center justify-between gap-3 jv-border-rough bg-jv-white px-3 py-2 shadow-soft-sm sm:min-h-16"
            >
              <div class="flex min-w-0 items-center gap-3">
                <span
                  class="grid size-9 shrink-0 place-items-center border border-jv-line sm:size-10"
                  :class="
                    participantAccentClasses[
                      index % participantAccentClasses.length
                    ]
                  "
                >
                  <Smile
                    v-if="index % 3 === 0"
                    class="size-5"
                    :stroke-width="2.2"
                  />
                  <UserRound v-else class="size-5" :stroke-width="2.2" />
                </span>
                <span
                  class="truncate font-body text-[16px] font-semibold text-jv-ink sm:text-[20px]"
                >
                  {{ getParticipantName(user) }}
                </span>
              </div>
              <span
                class="grid size-8 shrink-0 place-items-center border border-jv-line bg-jv-mint sm:size-9"
                aria-label="Připojen"
              >
                <Check class="size-5" :stroke-width="2.4" />
              </span>
            </div>
          </div>

          <div
            v-else
            class="grid min-h-[220px] flex-1 place-items-center border border-dashed border-jv-line bg-jv-canvas/50 p-6 text-center sm:min-h-[300px] sm:p-8"
          >
            <div>
              <div
                class="mx-auto grid size-14 place-items-center rounded-full border border-jv-line bg-jv-mint shadow-soft-sm sm:size-16"
              >
                <Users class="size-7 sm:size-8" :stroke-width="2.4" />
              </div>
              <p
                class="mt-5 font-body text-[18px] font-bold text-jv-muted sm:text-[22px]"
              >
                Čekáme na hráče…
              </p>
            </div>
          </div>

          <div class="mt-auto pt-6 sm:pt-8">
            <div class="border-t border-jv-line pt-5 text-center sm:pt-6">
              <div class="mb-3 flex justify-center gap-2 sm:mb-4">
                <span
                  class="size-2.5 rounded-full bg-jv-coral sm:size-3"
                ></span>
                <span
                  class="size-2.5 rounded-full bg-jv-yellow sm:size-3"
                ></span>
                <span class="size-2.5 rounded-full bg-jv-mint sm:size-3"></span>
              </div>
              <p
                class="font-body text-[18px] text-jv-muted sm:text-[22px] md:text-[24px]"
              >
                Čekáme na další hráče…
              </p>
            </div>
          </div>
        </aside>
      </div>
    </section>
    <QrFullscreenOverlay ref="qrOverlay" :join-url="joinUrl" :code="code" />
  </main>
  <main v-else class="flex min-h-screen flex-col bg-jv-canvas text-jv-ink">
    <header
      class="flex flex-wrap items-center justify-between gap-3 px-4 py-4 sm:gap-4 sm:px-6 sm:py-5 md:px-10"
    >
      <div class="flex min-w-0 items-center gap-3 sm:gap-4">
        <!-- <span
          class="grid size-11 shrink-0 place-items-center rounded-full border border-jv-line bg-jv-coral shadow-soft-sm sm:size-12"
        >
          <Zap
            class="size-5 text-white sm:size-6"
            :stroke-width="2.6"
            fill="white"
          />
        </span>
        <h1 class="font-headings text-[24px] tracking-tight sm:text-[30px]">
          JovVix
        </h1>
        <span
          class="hidden h-12 w-1 rotate-[10deg] rounded-full bg-jv-ink/20 sm:block"
          aria-hidden="true"
        ></span> -->
        <div class="flex min-w-0 flex-col">
          <p
            class="font-body text-[10px] font-bold uppercase tracking-[0.12em] text-jv-muted sm:text-[12px]"
          >
            Živý kvíz
          </p>
          <h2
            class="relative truncate font-headings text-[22px] leading-tight sm:text-[28px] md:text-[32px]"
          >
            <span class="relative z-10">{{ quizTitle }}</span>
            <span
              class="absolute bottom-[2px] left-0 z-0 h-[6px] w-[60%] max-w-[120px] bg-jv-yellow/70"
              aria-hidden="true"
            ></span>
          </h2>
        </div>
      </div>

      <div class="flex shrink-0 items-center gap-2 sm:gap-3">
        <button
          type="button"
          class="grid size-11 place-items-center rounded-[8px] border border-jv-line bg-jv-white text-jv-ink shadow-soft-sm transition-transform sm:size-12"
          :aria-label="music ? 'Vypnout hudbu' : 'Zapnout hudbu'"
          @click="toggleMusic"
        >
          <Volume2 v-if="music" class="size-5" :stroke-width="2.4" />
          <VolumeX v-else class="size-5" :stroke-width="2.4" />
        </button>
        <button
          type="button"
          class="inline-flex h-11 items-center justify-center gap-2 rounded-[8px] border border-jv-line bg-jv-white px-4 font-body text-[15px] font-bold text-jv-ink shadow-soft-sm transition-transform sm:h-12 sm:px-5 sm:text-[18px]"
          @click="leaveLobby"
        >
          <LogOut class="size-4 sm:size-5" :stroke-width="2.4" />
          <span>Opustit čekárnu</span>
        </button>
      </div>
    </header>

    <div class="flex justify-center px-4 sm:px-6 md:px-10">
      <div
        class="inline-flex max-w-full items-center gap-3 jv-border-rough bg-jv-white px-3 py-2 shadow-soft-sm sm:gap-4 sm:px-5 sm:py-2.5"
      >
        <span
          class="grid size-10 shrink-0 place-items-center overflow-hidden rounded-full border border-jv-line bg-jv-mint sm:size-12"
        >
          <NuxtImg
            v-if="userAvatar"
            :src="userAvatar"
            :alt="`Avatar hráče ${displayName}`"
            width="48"
            height="48"
            class="size-full object-cover"
          />
          <UserRound v-else class="size-5 sm:size-6" :stroke-width="2.4" />
        </span>
        <div class="flex min-w-0 flex-col">
          <p
            class="font-body text-[10px] font-bold uppercase tracking-[0.12em] text-jv-muted sm:text-[12px]"
          >
            Hraješ jako
          </p>
          <p
            class="min-w-0 truncate font-headings text-[18px] leading-tight text-jv-ink sm:text-[22px]"
          >
            {{ displayName }}
          </p>
        </div>
      </div>
    </div>

    <section
      class="flex flex-1 items-center justify-center px-4 py-6 sm:px-6 sm:py-10 md:px-10 md:py-12"
    >
      <div class="relative w-full max-w-[980px]">
        <div
          class="jv-border-rough bg-jv-white p-6 shadow-soft-lg sm:p-10 md:p-12"
        >
          <div class="border-b border-dashed border-jv-line pb-6 sm:pb-7">
            <h2
              class="font-headings text-[28px] leading-tight sm:text-[36px] md:text-[48px]"
            >
              Připravit, pozor, start!
            </h2>
            <p
              class="mt-1 font-body text-[14px] font-bold text-jv-muted sm:text-[18px]"
            >
              Kvíz brzy začne
            </p>
          </div>

          <div
            class="flex flex-wrap items-center justify-center gap-4 px-2 py-12 text-center sm:gap-5 sm:py-20 md:py-24"
          >
            <h3
              class="font-headings text-[30px] leading-tight sm:text-[44px] md:text-[60px]"
            >
              {{ waitingMessage }}
            </h3>
            <span
              class="relative grid size-14 shrink-0 place-items-center rounded-full border border-jv-line bg-jv-mint shadow-soft sm:size-16"
            >
              <Hourglass
                class="size-6 animate-hourglass-flip sm:size-7"
                :stroke-width="2.4"
              />
              <span
                class="absolute inset-1 rounded-full border border-dashed border-jv-line"
                aria-hidden="true"
              ></span>
            </span>
          </div>
        </div>
      </div>
    </section>
  </main>
</template>
<style scoped>
.code {
  letter-spacing: 0.4rem;
}

@media (max-width: 768px) {
  .code {
    letter-spacing: 0.22rem;
  }
}

.qr-card :deep(svg) {
  width: 100%;
  height: 100%;
}
</style>
