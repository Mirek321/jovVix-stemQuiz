<script setup>
// Host screen for a self-paced session: no lobby or question loop, just the
// join details and a live table of how far each player got.
import {
  Ban,
  Check,
  Copy,
  Footprints,
  Home,
  Keyboard,
  Link2,
  Maximize2,
  UserRound,
  Users,
} from "lucide-vue-next";
import usecopyToClipboard from "~~/composables/copy_to_clipboard";
import { getAvatarUrlByName } from "~~/composables/avatar";
import { czCount } from "~~/composables/plural";
import { fetchSessionStatus } from "~~/composables/self_paced";

const STATUS_POLL_MS = 3000;

const props = defineProps({
  // Initial GET /active_quizzes/:id/status payload; refreshed by polling here.
  status: {
    type: Object,
    required: true,
  },
  terminating: {
    default: false,
    type: Boolean,
    required: false,
  },
});
const emits = defineEmits(["endQuiz", "ended"]);

const { apiUrl, baseUrl } = useRuntimeConfig().public;

const current = ref(props.status);
const qrOverlay = ref(null);
let pollTimer = null;

const code = computed(() =>
  current.value?.code != null ? String(current.value.code) : ""
);
const joinUrl = computed(() => `${baseUrl}/join`);
const joinUrlDisplay = computed(() =>
  joinUrl.value.replace(/^https?:\/\//, "")
);
const joinUrlWithCode = computed(() => `${joinUrl.value}?code=${code.value}`);

const quizTitle = computed(() => current.value?.title || "Kvíz");
const totalQuestions = computed(
  () => Number(current.value?.total_questions) || 0
);
const players = computed(() => current.value?.players || []);
const completedCount = computed(
  () => players.value.filter((player) => player.completed).length
);

const playerName = (player) => player.first_name || player.username || "Hráč";

const stopPolling = () => {
  clearInterval(pollTimer);
  pollTimer = null;
};

// Ended from elsewhere (another tab, the reports page): leave for results.
const checkEnded = () => {
  if (current.value?.is_active !== false) return false;
  stopPolling();
  emits("ended");
  return true;
};

const refresh = async () => {
  try {
    const status = await fetchSessionStatus(apiUrl, current.value.session_id);
    if (!status) return;
    current.value = status;
    checkEnded();
  } catch (error) {
    console.error("failed to refresh self-paced session status", error);
  }
};

const openQrFullscreen = () => qrOverlay.value?.open();

onMounted(() => {
  if (checkEnded()) return;
  pollTimer = setInterval(refresh, STATUS_POLL_MS);
});

onUnmounted(stopPolling);
</script>

<template>
  <main
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
            class="inline-flex items-center gap-2 font-body text-[11px] font-black uppercase tracking-[0.18em] text-jv-muted sm:justify-self-center sm:text-[13px]"
          >
            <Footprints class="size-4" :stroke-width="2.4" />
            Vlastním tempem
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
          Kvíz běží. Hráči začnou hrát hned po připojení, každý svým tempem.
        </p>
      </header>

      <div
        class="grid flex-1 gap-6 md:gap-7 xl:grid-cols-[minmax(340px,0.8fr)_minmax(0,1.2fr)] xl:items-start"
      >
        <!-- Join details -->
        <section
          class="relative flex flex-col overflow-hidden bg-jv-white shadow-soft-sm jv-border-rough sm:shadow-soft-lg"
        >
          <span
            class="absolute left-1/2 top-[-2px] z-10 h-4 w-12 -translate-x-1/2 bg-jv-coral"
            aria-hidden="true"
          ></span>

          <div class="bg-jv-yellow px-4 pb-6 pt-6 text-center sm:px-8">
            <p
              class="inline-flex items-center gap-2 font-body text-[12px] font-black uppercase tracking-[0.12em] text-jv-muted"
            >
              <Keyboard class="size-4" :stroke-width="2.4" />
              Kód kvízu
            </p>
            <p
              class="mt-2 break-all font-feature text-[52px] font-black leading-none tracking-[0.12em] text-jv-ink sm:text-[68px]"
            >
              {{ code }}
            </p>
            <p
              class="mt-3 font-body text-[14px] font-bold text-jv-muted sm:text-[16px]"
            >
              Připoj se na
              <strong class="text-jv-ink">{{ joinUrlDisplay }}</strong>
            </p>
          </div>

          <div class="flex flex-col items-center gap-4 px-4 py-6 sm:px-8">
            <div class="flex w-full flex-col gap-2 sm:flex-row">
              <button
                type="button"
                class="inline-flex h-11 flex-1 items-center justify-center gap-2 rounded-full border border-jv-line bg-jv-coral px-5 font-body text-[15px] font-bold text-white shadow-soft-sm transition-transform"
                @click="usecopyToClipboard(code)"
              >
                <Copy class="size-4" :stroke-width="2.5" />
                Kopírovat kód
              </button>
              <button
                type="button"
                class="inline-flex h-11 flex-1 items-center justify-center gap-2 rounded-full border border-jv-line bg-jv-white px-5 font-body text-[15px] font-bold text-jv-ink shadow-soft-sm transition-transform"
                @click="usecopyToClipboard(joinUrlWithCode)"
              >
                <Link2 class="size-4" :stroke-width="2.5" />
                Kopírovat odkaz
              </button>
            </div>

            <div
              class="grid size-[176px] place-items-center bg-jv-white p-3 shadow-soft-sm jv-border-rough min-[420px]:size-[196px] sm:size-[220px] sm:p-4"
            >
              <QrCode :scan-u-r-l="joinUrl" :quiz-code="code" :size="190" />
            </div>
            <button
              type="button"
              class="inline-flex h-11 items-center justify-center gap-2 rounded-full border border-jv-line bg-jv-mint px-5 font-body text-[14px] font-black text-jv-ink shadow-soft-sm transition-transform"
              @click="openQrFullscreen"
            >
              <Maximize2 class="size-4" :stroke-width="2.4" />
              Zvětšit QR kód
            </button>
          </div>
        </section>

        <!-- Player progress -->
        <aside
          class="flex flex-col bg-jv-white px-4 py-6 shadow-soft-sm jv-border-rough sm:px-8 sm:py-8 sm:shadow-soft-lg"
        >
          <div class="flex items-start justify-between gap-3 sm:items-center">
            <div class="flex min-w-0 items-center gap-3 sm:gap-4">
              <span
                class="grid size-10 shrink-0 place-items-center jv-border-rough bg-jv-mint sm:size-12"
              >
                <Users class="size-5 sm:size-6" :stroke-width="2.4" />
              </span>
              <div class="min-w-0">
                <h2
                  class="font-headings text-[28px] leading-tight text-jv-ink sm:text-[38px]"
                >
                  Hráči
                </h2>
                <p
                  v-if="players.length"
                  class="font-body text-[13px] font-bold text-jv-muted sm:text-[14px]"
                >
                  Hotovo {{ completedCount }} z
                  {{ czCount(players.length, "hráče", "hráčů", "hráčů") }}
                </p>
              </div>
            </div>
            <span
              class="grid size-11 shrink-0 place-items-center border border-jv-line bg-jv-yellow font-feature text-[24px] font-black shadow-soft-sm sm:size-14 sm:text-[30px]"
            >
              {{ players.length }}
            </span>
          </div>

          <div class="my-5 border-t border-dashed border-jv-line"></div>

          <div v-if="players.length" class="-mx-1 overflow-x-auto px-1">
            <table
              class="w-full min-w-[440px] border-separate border-spacing-y-2"
            >
              <thead>
                <tr
                  class="text-left font-body text-[11px] font-black uppercase tracking-[0.12em] text-jv-muted"
                >
                  <th class="px-3 font-black">Jméno</th>
                  <th class="px-3 text-center font-black">Odpovězeno</th>
                  <th class="px-3 text-right font-black">Skóre</th>
                  <th class="px-3 text-center font-black">Hotovo</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="player in players"
                  :key="player.user_played_quiz_id"
                  class="bg-jv-white shadow-soft-sm"
                >
                  <td
                    class="rounded-l-2xl border-y border-l border-jv-line px-3 py-2"
                  >
                    <div class="flex min-w-0 items-center gap-3">
                      <span
                        class="grid size-9 shrink-0 place-items-center overflow-hidden rounded-full border border-jv-line bg-jv-mint"
                      >
                        <NuxtImg
                          v-if="player.img_key"
                          :src="getAvatarUrlByName(player.img_key)"
                          :alt="`Avatar hráče ${playerName(player)}`"
                          width="36"
                          height="36"
                          class="size-full object-cover"
                        />
                        <UserRound v-else class="size-5" :stroke-width="2.2" />
                      </span>
                      <span
                        class="truncate font-body text-[15px] font-bold text-jv-ink sm:text-[17px]"
                      >
                        {{ playerName(player) }}
                      </span>
                    </div>
                  </td>
                  <td
                    class="border-y border-jv-line px-3 py-2 text-center font-feature text-[15px] font-black"
                  >
                    {{ player.answered }}/{{ totalQuestions }}
                  </td>
                  <td
                    class="border-y border-jv-line px-3 py-2 text-right font-feature text-[16px] font-black"
                  >
                    {{ player.score }}
                  </td>
                  <td
                    class="rounded-r-2xl border-y border-r border-jv-line px-3 py-2 text-center"
                  >
                    <span
                      v-if="player.completed"
                      class="inline-grid size-8 place-items-center rounded-full bg-jv-mint"
                      aria-label="Hotovo"
                    >
                      <Check class="size-4" :stroke-width="3" />
                    </span>
                    <span
                      v-else
                      class="font-body text-[13px] font-bold text-jv-muted"
                    >
                      Hraje
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <div
            v-else
            class="grid min-h-[220px] place-items-center border border-dashed border-jv-line bg-jv-canvas/50 p-6 text-center"
          >
            <div>
              <div
                class="mx-auto grid size-14 place-items-center rounded-full border border-jv-line bg-jv-mint shadow-soft-sm"
              >
                <Users class="size-7" :stroke-width="2.4" />
              </div>
              <p
                class="mt-5 font-body text-[18px] font-bold text-jv-muted sm:text-[20px]"
              >
                Čekáme na první hráče…
              </p>
            </div>
          </div>

          <div class="mt-6 border-t border-dashed border-jv-line pt-5">
            <button
              type="button"
              :disabled="terminating"
              class="mx-auto flex h-12 w-full max-w-[320px] items-center justify-center gap-2 rounded-full border border-jv-line bg-jv-coral px-5 font-body text-[16px] font-black text-white shadow-soft-sm transition-transform disabled:cursor-not-allowed disabled:opacity-60"
              @click="emits('endQuiz')"
            >
              <Ban class="size-4" :stroke-width="2.4" />
              <span>{{ terminating ? "Ukončuji…" : "Ukončit kvíz" }}</span>
            </button>
          </div>
        </aside>
      </div>
    </section>
    <QrFullscreenOverlay ref="qrOverlay" :join-url="joinUrl" :code="code" />
  </main>
</template>
