<template>
  <main class="min-h-screen bg-jv-canvas px-4 py-5 sm:px-6 md:px-8 md:py-6">
    <div
      v-if="quizError?.data?.code == 401"
      class="jv-border-rough bg-jv-white p-5 text-[18px] font-semibold text-jv-coral shadow-soft-sm"
    >
      {{ navigateTo("/account/login") }}
    </div>

    <div
      v-else-if="quizError"
      class="jv-border-rough bg-jv-white p-5 text-[18px] font-semibold text-jv-coral shadow-soft-sm"
    >
      {{ quizError.message || quizError }}
    </div>

    <section v-else class="mx-auto flex flex-col gap-6">
      <header
        class="flex flex-col gap-4 md:flex-row md:items-start md:justify-between"
      >
        <div class="min-w-0">
          <div class="flex flex-wrap items-center gap-3">
            <h1
              class="break-words font-headings text-[38px] leading-none text-jv-ink min-[420px]:text-[44px] sm:text-[52px] md:text-[56px]"
            >
              {{ quizTitle }}
            </h1>
            <span
              v-if="isPublicQuiz"
              class="inline-flex items-center gap-1.5 border border-jv-line bg-jv-mint px-3 py-1.5 text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink shadow-soft-sm sm:text-[14px]"
              aria-label="Veřejný kvíz"
            >
              <Globe class="size-4" :stroke-width="2.6" />
              Veřejný
            </span>
          </div>
          <div
            class="ml-1 mt-1 h-3 w-24 rounded-full border-b border-jv-coral sm:ml-2 sm:w-32"
          ></div>
          <p
            class="mt-3 max-w-2xl text-[16px] leading-[1.5] text-jv-muted sm:text-[18px]"
          >
            {{ quizDescription || "Zatím bez popisu." }}
          </p>
        </div>
        <NavigationLink
          url-name="Spustit kvíz"
          class="rounded-[999px] bg-jv-coral text-white"
          @click="askStartMode"
        >
          <Play class="size-4 fill-current" :stroke-width="2.4" />
        </NavigationLink>
      </header>

      <div
        class="flex flex-col gap-4 jv-border-rough bg-jv-white p-3 shadow-soft-sm sm:p-4 lg:flex-row lg:items-center lg:justify-between"
      >
        <div class="flex flex-wrap gap-2">
          <span
            class="inline-flex items-center gap-1.5 rounded-[6px] border border-dashed border-jv-line bg-jv-white px-3 py-2 text-[14px] font-semibold text-jv-muted"
          >
            Počet otázek:
            <strong class="text-jv-ink">{{ questions.length }}</strong>
          </span>
          <span
            class="inline-flex items-center gap-1.5 rounded-[6px] border border-dashed border-jv-line bg-jv-white px-3 py-2 text-[14px] font-semibold text-jv-muted"
          >
            Anketní otázky:
            <strong class="text-jv-ink">{{ totalSurveyQuestion }}</strong>
          </span>
          <span
            class="inline-flex items-center gap-1.5 rounded-[6px] border border-dashed border-jv-line bg-jv-white px-3 py-2 text-[14px] font-semibold text-jv-muted"
          >
            Odehráno:
            <strong class="text-jv-ink">{{
              quizData?.data?.quiz_played_count || 0
            }}</strong>
          </span>
        </div>

        <div class="flex flex-wrap gap-2 sm:justify-end">
          <NavigationLink
            url-name="Smazat kvíz"
            class="bg-jv-white text-jv-coral border-none shadow-none"
            @click="deleteModalOpen = true"
          >
            <Trash2 class="size-4" :stroke-width="2.4" />
          </NavigationLink>
          <NavigationLink
            url-name="Sdílet kvíz"
            class="bg-jv-white"
            @click="handleShareQuiz"
          >
            <Share2 class="size-4" :stroke-width="2.4" />
          </NavigationLink>
        </div>
      </div>

      <div
        v-if="canEditPublicMeta && canEditQuiz"
        class="grid gap-5 sm:grid-cols-2 sm:gap-x-8"
      >
        <label class="grid content-start gap-2">
          <span
            class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink"
          >
            Kategorie
          </span>
          <select
            v-model="settings.category_id"
            :disabled="settingsPending"
            class="h-14 border border-jv-line bg-jv-canvas px-4 text-[16px] font-semibold text-jv-ink outline-none transition-shadow focus:shadow-soft-sm disabled:opacity-60"
          >
            <option value="">Bez kategorie (zobrazí se v „Ostatní“)</option>
            <option
              v-for="category in categories"
              :key="category.id"
              :value="category.id"
            >
              {{ category.name }}
            </option>
          </select>
        </label>

        <div class="grid content-start gap-2">
          <span
            class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink"
          >
            Titulní obrázek
          </span>
          <label
            v-if="!settings.cover_image"
            class="flex h-14 cursor-pointer items-center justify-center gap-2 border border-dashed border-jv-line bg-jv-canvas px-4 text-[15px] font-semibold text-jv-muted transition-colors hover:bg-jv-yellow/20"
          >
            <ImageIcon class="size-4 shrink-0" :stroke-width="2.3" />
            <span class="truncate">{{
              processingCover
                ? "Zpracovávám obrázek…"
                : "Nahrát titulní obrázek"
            }}</span>
            <input
              type="file"
              class="hidden"
              accept="image/*"
              :disabled="settingsPending || processingCover"
              @change="handleCoverImage"
            />
          </label>
          <div
            v-else
            class="flex h-14 items-center gap-3 border border-jv-line bg-jv-canvas px-3"
          >
            <NuxtImg
              :src="settings.cover_image"
              alt="Náhled titulního obrázku"
              class="h-10 w-14 shrink-0 border border-jv-line object-cover"
            />
            <span
              class="min-w-0 flex-1 truncate text-[14px] font-semibold text-jv-ink"
            >
              {{ coverImageName || "Aktuální titulní obrázek" }}
            </span>
            <button
              type="button"
              class="grid size-8 shrink-0 place-items-center text-jv-ink transition-colors hover:text-jv-coral"
              aria-label="Odebrat titulní obrázek"
              :disabled="settingsPending"
              @click="removeCoverImage"
            >
              <X class="size-4" :stroke-width="2.4" />
            </button>
          </div>
        </div>
      </div>

      <div class="grid gap-4 lg:grid-cols-[minmax(0,1fr)_auto] lg:items-center">
        <div
          v-if="questions.length > 0"
          class="flex flex-wrap items-center gap-4"
        >
          <TooltipProvider>
            <label
              class="inline-flex items-center gap-2 text-[15px] font-semibold text-jv-muted"
            >
              <span class="inline-flex items-center gap-1.5">
                Body
                <Tooltip>
                  <TooltipTrigger as-child>
                    <button
                      type="button"
                      class="inline-flex size-5 items-center justify-center rounded-full border border-jv-line bg-jv-white text-jv-ink shadow-soft-sm transition-transform"
                      aria-label="Co znamenají body?"
                    >
                      <Info class="size-3" :stroke-width="2.6" />
                    </button>
                  </TooltipTrigger>
                  <TooltipContent side="bottom" align="start">
                    Body se udělují celé za každou správnou odpověď a ke skóre
                    hráče přidají body &times; 100. Navíc každá správná odpověď
                    přináší 500 základních bodů plus bonus za rychlost až 400 a
                    další body přidávají série správných odpovědí.
                  </TooltipContent>
                </Tooltip>
              </span>
              <input
                v-model.number="settings.points"
                type="number"
                :min="constants.MinimumPoints"
                :max="constants.MaximumPoints"
                :disabled="!canEditQuiz || settingsPending"
                :aria-invalid="!!settingsErrors.points"
                :aria-describedby="
                  settingsErrors.points ? 'quiz-points-error' : undefined
                "
                class="h-10 w-14 border border-jv-line bg-jv-white px-2 text-center text-[16px] font-black text-jv-ink shadow-soft-sm outline-none disabled:opacity-60"
              />
              <span
                v-if="settingsErrors.points"
                id="quiz-points-error"
                :title="settingsErrors.points"
                class="text-[20px] font-black leading-none text-jv-coral"
              >
                <span aria-hidden="true">*</span>
                <span class="sr-only">{{ settingsErrors.points }}</span>
              </span>
            </label>
            <label
              class="inline-flex items-center gap-2 text-[15px] font-semibold text-jv-muted"
            >
              <span class="inline-flex items-center gap-1.5">
                Časový limit (s)
                <Tooltip>
                  <TooltipTrigger as-child>
                    <button
                      type="button"
                      class="inline-flex size-5 items-center justify-center rounded-full border border-jv-line bg-jv-white text-jv-ink shadow-soft-sm transition-transform"
                      aria-label="Co znamená časový limit?"
                    >
                      <Info class="size-3" :stroke-width="2.6" />
                    </button>
                  </TooltipTrigger>
                  <TooltipContent side="bottom" align="start">
                    Kolik sekund má každý hráč na odpověď, než se otázka
                    automaticky přeskočí. Kratší limit kvíz zrychlí, ale zároveň
                    ztíží.
                  </TooltipContent>
                </Tooltip>
              </span>
              <input
                v-model.number="settings.duration_in_seconds"
                type="number"
                :min="constants.MinimumDurationInSeconds"
                :max="constants.MaximumDurationInSeconds"
                :disabled="!canEditQuiz || settingsPending"
                :aria-invalid="!!settingsErrors.duration"
                :aria-describedby="
                  settingsErrors.duration ? 'quiz-duration-error' : undefined
                "
                class="h-10 w-16 border border-jv-line bg-jv-white px-2 text-center text-[16px] font-black text-jv-ink shadow-soft-sm outline-none disabled:opacity-60"
              />
              <span
                v-if="settingsErrors.duration"
                id="quiz-duration-error"
                :title="settingsErrors.duration"
                class="text-[20px] font-black leading-none text-jv-coral"
              >
                <span aria-hidden="true">*</span>
                <span class="sr-only">{{ settingsErrors.duration }}</span>
              </span>
            </label>
          </TooltipProvider>
        </div>

        <div
          v-if="canEditQuiz"
          class="flex flex-col gap-3 sm:flex-row sm:justify-end"
        >
          <NavigationLink
            v-if="canEditQuiz && hasUnsavedSettings"
            url-name="Uložit změny"
            class="rounded-[999px] bg-jv-white"
            :disabled="settingsPending"
            @click="saveSettings"
          />
          <NavigationLink
            url-name="Importovat otázky z CSV"
            class="rounded-[999px] bg-jv-white"
            @click="importModalOpen = true"
          >
            <Upload class="size-4" :stroke-width="2.4" />
          </NavigationLink>
          <NavigationLink
            url-name="Přidat otázku"
            class="rounded-[999px] bg-jv-accent-green text-white"
            @click="openNewQuestionForm"
          >
            <Plus class="size-4" :stroke-width="2.4" />
          </NavigationLink>
        </div>
      </div>

      <QuestionFormCard
        v-if="showNewQuestionForm && canEditQuiz"
        mode="create"
        :saving="savingNewQuestion"
        @save="saveNewQuestion"
        @cancel="showNewQuestionForm = false"
      />

      <section
        v-if="questions.length === 0 && !showNewQuestionForm"
        class="grid min-h-[260px] place-items-center text-center"
      >
        <div class="max-w-md py-8">
          <h2
            class="font-headings text-[30px] leading-tight text-jv-ink sm:text-[36px]"
          >
            Zatím žádné otázky
          </h2>
          <p class="mt-2 text-[17px] text-jv-muted">
            Vytvoř otázku ručně nebo naimportuj CSV a začni tvořit kvíz.
          </p>
        </div>
      </section>

      <draggable
        v-else
        v-model="orderedQuestions"
        :component-data="{
          tag: 'section',
          name: 'flip-list',
          type: 'transition-group',
        }"
        class="grid gap-5"
        item-key="question_id"
        handle=".drag-handle"
        ghost-class="question-ghost"
        chosen-class="question-chosen"
        drag-class="question-drag"
        :disabled="!canEditQuiz"
        :animation="180"
        :scroll="scrollContainer"
        :bubble-scroll="true"
        :force-autoscroll-fallback="true"
        :scroll-sensitivity="80"
        :scroll-speed="14"
      >
        <template #item="{ element: question, index }">
          <article
            :class="
              editingQuestionId === question.question_id
                ? ''
                : 'jv-border-rough bg-jv-white p-4 shadow-soft-sm transition-transform sm:p-5'
            "
          >
            <QuestionFormCard
              v-if="editingQuestionId === question.question_id"
              mode="edit"
              :question="question"
              :saving="savingQuestionId === question.question_id"
              @save="(data) => saveExistingQuestion(question, data)"
              @cancel="editingQuestionId = ''"
              @click.stop
            />

            <template v-else>
              <div class="flex items-start justify-between gap-4">
                <div class="flex min-w-0 items-start gap-2">
                  <button
                    v-if="canEditQuiz"
                    type="button"
                    class="drag-handle mt-1 grid size-7 shrink-0 cursor-grab place-items-center rounded text-jv-muted hover:bg-jv-canvas hover:text-jv-ink active:cursor-grabbing"
                    aria-label="Přetažením změníš pořadí"
                    @click.stop
                  >
                    <GripVertical class="size-4" :stroke-width="2.4" />
                  </button>
                  <div class="min-w-0">
                    <p
                      class="text-[12px] font-bold uppercase tracking-[0.14em] text-jv-coral"
                    >
                      Otázka {{ index + 1 }}
                    </p>
                    <h2
                      class="mt-1 break-words text-[20px] font-bold leading-snug text-jv-ink sm:text-[22px]"
                    >
                      {{ question.question }}
                    </h2>
                  </div>
                </div>
                <div v-if="canEditQuiz" class="flex shrink-0 gap-2">
                  <NavigationLink
                    class="rounded-full size-10 px-2 py-1 sm:px-2 sm:py-1 md:px-2 md:py-1"
                    aria-label="Upravit otázku"
                    @click.stop="startEditQuestion(question)"
                  >
                    <Pencil class="size-4" :stroke-width="2.4" />
                  </NavigationLink>
                  <NavigationLink
                    class="rounded-full bg-jv-coral text-white size-10 px-2 py-1 sm:px-2 sm:py-1 md:px-2 md:py-1"
                    aria-label="Smazat otázku"
                    @click.stop="deleteQuestion(question.question_id)"
                  >
                    <Trash2 class="size-4" :stroke-width="2.4" />
                  </NavigationLink>
                </div>
              </div>

              <div
                v-if="question.question_media === 'image' && question.resource"
                class="mt-3 flex justify-center rounded-md bg-jv-canvas p-2"
                @click.stop
              >
                <NuxtImg
                  :src="question.resource"
                  :alt="question.question"
                  class="max-h-64 w-auto max-w-full object-contain"
                />
              </div>

              <div
                v-else-if="
                  question.question_media === 'code' && question.resource
                "
                class="mt-3 min-w-0 overflow-x-auto"
                @click.stop
              >
                <CodeBlockComponent :code="question.resource" />
              </div>

              <ul class="mt-4 flex flex-col">
                <li
                  v-for="(option, key) in question.options"
                  :key="key"
                  class="flex min-w-0 items-center gap-3 border-b border-jv-line py-3 pl-3 pr-2 text-[15px] font-medium text-jv-ink last:border-b-0"
                  :class="
                    isCorrectAnswer(question, key)
                      ? 'border-l-4 border-l-jv-accent-green bg-jv-accent-green/25 pl-2'
                      : 'border-l-4 border-l-transparent'
                  "
                >
                  <span
                    class="w-5 shrink-0 text-[14px] font-bold text-jv-coral"
                  >
                    {{ optionLabel(key) }}.
                  </span>

                  <div
                    v-if="question.options_media === 'image' && option"
                    class="flex min-w-0 flex-1 justify-start"
                    @click.stop
                  >
                    <NuxtImg
                      :src="option"
                      :alt="`Možnost ${optionLabel(key)}`"
                      class="max-h-32 w-auto max-w-full object-contain"
                    />
                  </div>

                  <div
                    v-else-if="question.options_media === 'code' && option"
                    class="min-w-0 flex-1 overflow-x-auto"
                    @click.stop
                  >
                    <CodeBlockComponent :code="option" />
                  </div>

                  <span v-else class="min-w-0 flex-1 break-words">
                    {{ option }}
                  </span>

                  <Check
                    v-if="isCorrectAnswer(question, key)"
                    class="size-5 shrink-0 text-jv-accent-green"
                    :stroke-width="3"
                  />
                </li>
              </ul>
            </template>
          </article>
        </template>
      </draggable>
    </section>

    <Teleport to="body">
      <div
        v-if="importModalOpen"
        class="fixed inset-0 z-50 grid place-items-center bg-jv-ink/35 px-4 py-6 backdrop-blur-[2px]"
        @click.self="closeImportModal"
      >
        <form
          class="w-full max-w-[680px] border-[4px] border-jv-line bg-jv-white shadow-soft-lg"
          @submit.prevent="importCsv"
        >
          <div
            class="flex items-center justify-between gap-4 border-b border-jv-line bg-jv-ink px-5 py-4 text-jv-white sm:px-6"
          >
            <h2
              class="font-body text-[24px] font-black leading-none text-jv-white sm:text-[28px]"
            >
              Přidat otázky do kvízu
            </h2>
            <button
              type="button"
              class="grid size-9 place-items-center text-jv-white transition-transform"
              aria-label="Zavřít okno importu CSV"
              @click="closeImportModal"
            >
              <X class="size-6" :stroke-width="2.4" />
            </button>
          </div>

          <div class="grid gap-4 px-5 py-6 sm:px-8">
            <label class="grid gap-2">
              <span
                class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink"
              >
                Vyber soubor CSV <span class="text-jv-coral">*</span>
              </span>
              <span
                class="flex h-14 cursor-pointer border border-jv-line bg-jv-canvas"
              >
                <span
                  class="inline-flex h-full items-center gap-2 bg-jv-ink px-4 text-[16px] font-black text-jv-white"
                >
                  <Upload class="size-4" :stroke-width="2.4" />
                  Vybrat soubor
                </span>
                <span
                  class="flex min-w-0 flex-1 items-center px-4 text-[15px] font-semibold text-jv-muted"
                >
                  {{ csvFileName || "Není vybrán žádný soubor" }}
                </span>
                <input
                  type="file"
                  class="hidden"
                  accept=".csv,text/csv"
                  required
                  @change="handleCsvFile"
                />
              </span>
            </label>
            <p class="text-[15px] leading-[1.6] text-jv-muted">
              Nahraj soubor CSV s dalšími otázkami. Nové řádky se přidají na
              konec kvízu.
            </p>
            <NavigationLink
              url="/files/demo.csv"
              url-name="Stáhnout ukázku"
              external
              download="demo.csv"
              class="w-fit rounded-[999px] bg-jv-white"
            >
              <Download class="size-4" :stroke-width="2.4" />
            </NavigationLink>
          </div>

          <div
            class="flex flex-col-reverse gap-3 border-t border-jv-line bg-jv-canvas px-5 py-4 sm:flex-row sm:justify-end sm:px-8"
          >
            <button
              type="button"
              class="inline-flex h-12 items-center justify-center border border-jv-line bg-jv-white px-6 text-[17px] font-black text-jv-ink shadow-soft-sm transition-transform"
              @click="closeImportModal"
            >
              Zrušit
            </button>
            <button
              type="submit"
              class="inline-flex h-12 items-center justify-center border border-jv-line bg-jv-accent-green px-6 text-[17px] font-black text-jv-ink shadow-soft-sm transition-transform disabled:cursor-not-allowed disabled:opacity-60"
              :disabled="importPending"
            >
              {{ importPending ? "Přidávám…" : "Přidat otázky" }}
            </button>
          </div>
        </form>
      </div>
    </Teleport>

    <ShareQuizModal v-model="shareModalOpen" :quiz-id="quizId" />
    <StartModeModal v-model="startModeOpen" @select="handleStartQuiz" />
    <DeleteDialog
      v-model="deleteModalOpen"
      title="Smazat kvíz"
      :message="deleteQuizMessage"
      confirm-label="Smazat kvíz"
      @confirm-delete="deleteQuiz"
    />
  </main>
</template>

<script setup>
import { computed, onMounted, ref, watch } from "vue";
import {
  Check,
  Download,
  Globe,
  GripVertical,
  Image as ImageIcon,
  Info,
  Pencil,
  Play,
  Plus,
  Share2,
  Trash2,
  Upload,
  X,
} from "lucide-vue-next";
import { usePush } from "notivue";
import draggable from "vuedraggable";
import QuestionFormCard from "@/components/quiz-manage/QuestionFormCard.vue";
import CodeBlockComponent from "@/components/CodeBlockComponent.vue";
import ShareQuizModal from "@/components/Quiz/ShareQuizModal.vue";
import StartModeModal from "@/components/Quiz/StartModeModal.vue";
import constants from "~~/config/constants";
import { readApiError } from "@/composables/apiError";
import { useListUserstore } from "~/store/userlist";
import { useSessionStore } from "~~/store/session";
import NavigationLink from "~/components/common/NavigationLink.vue";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";

definePageMeta({
  layout: "empty",
});

useSeoMeta({
  title: "Upravit kvíz – STEM_HUB",
  description: "Uprav detaily kvízu a spravuj jeho otázky v systému STEM_HUB.",
  robots: "noindex, nofollow",
});

const toast = usePush();
const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const route = useRoute();
const router = useRouter();
const sessionStore = useSessionStore();
const listUserStore = useListUserstore();

const quizId = computed(() => route.params.quiz_id || "");
const importModalOpen = ref(false);
const shareModalOpen = ref(false);
const deleteModalOpen = ref(false);
const importPending = ref(false);
const csvFile = ref(null);
const csvFileName = ref("");
const showNewQuestionForm = ref(false);
const showAddAnother = ref(false);
const savingNewQuestion = ref(false);
const editingQuestionId = ref("");
const savingQuestionId = ref("");
const settingsPending = ref(false);
const startingQuiz = ref(false);
const startModeOpen = ref(false);

const {
  refresh,
  data: quizData,
  error: quizError,
} = await useFetch(`${url.apiUrl}/quizzes/${quizId.value}/questions`, {
  method: "GET",
  headers: headers,
  mode: "cors",
  credentials: "include",
});

const settings = ref({
  points: Number(quizData.value?.data?.points ?? constants.DefaultPoints),
  duration_in_seconds: Number(
    quizData.value?.data?.duration_in_seconds ??
      constants.DefaultDurationInSeconds
  ),
  category_id: nullableString(quizData.value?.data?.category_id),
  cover_image: nullableString(quizData.value?.data?.cover_image),
});
const coverImageName = ref("");
const processingCover = ref(false);
const { pickCoverImage } = useCoverImage();

const { data: categoriesData } = await useFetch(`${url.apiUrl}/categories`, {
  method: "GET",
  headers: headers,
  credentials: "include",
});
const categories = computed(() => categoriesData.value?.data || []);

const orderedQuestions = ref([...(quizData.value?.data?.data || [])]);
const originalOrderIds = ref(
  (quizData.value?.data?.data || []).map((q) => q.question_id)
);
const scrollContainer = ref(null);

onMounted(() => {
  scrollContainer.value =
    document.querySelector(".lg\\:overflow-y-auto") ||
    document.scrollingElement;
});

const questions = computed(() => quizData.value?.data?.data || []);
const quizTitle = computed(() => quizData.value?.data?.quiz_title || "Kvíz");

const deleteQuizMessage = computed(
  () =>
    `Smazáním kvízu „${quizTitle.value}“ odstraníš i všechny jeho otázky, všechny odehrané hry a skóre i odpovědi hráčů. Tuto akci nejde vrátit.`
);
const quizDescription = computed(
  () => quizData.value?.data?.quiz_description?.String || ""
);
const isPublicQuiz = computed(() => !!quizData.value?.data?.is_public);
// The server decides this: public quiz + the configured admin allowlist.
const canEditPublicMeta = computed(
  () => !!quizData.value?.data?.can_edit_public_meta
);
const totalSurveyQuestion = computed(() =>
  questions.value.reduce(
    (count, item) => (item.question_type_id === 2 ? count + 1 : count),
    0
  )
);
const canEditQuiz = computed(() => {
  const permission = quizData.value?.data?.permission;
  const isEditable = quizData.value?.data?.is_quiz_editable;
  const hasPermission = permission === "write" || permission === "share";
  return hasPermission && (isPublicQuiz.value || isEditable);
});

// Mirrors the server rules on ReqUpdateQuizSettings so an out-of-range value is
// caught here instead of coming back as a bare 400.
const settingsErrors = computed(() => {
  const errors = {};
  const points = settings.value.points;
  const duration = settings.value.duration_in_seconds;

  if (
    !Number.isInteger(points) ||
    points < constants.MinimumPoints ||
    points > constants.MaximumPoints
  ) {
    errors.points = `Body musí být v rozmezí ${constants.MinimumPoints} až ${constants.MaximumPoints}.`;
  }
  if (
    !Number.isInteger(duration) ||
    duration < constants.MinimumDurationInSeconds ||
    duration > constants.MaximumDurationInSeconds
  ) {
    errors.duration = `Časový limit musí být v rozmezí ${constants.MinimumDurationInSeconds} až ${constants.MaximumDurationInSeconds} sekund.`;
  }

  return errors;
});

const settingsError = computed(
  () => settingsErrors.value.duration || settingsErrors.value.points || ""
);

const hasUnsavedSettings = computed(() => {
  const data = quizData.value?.data;
  if (!data) return false;
  const currentIds = orderedQuestions.value.map((q) => q.question_id);
  const orderChanged =
    currentIds.length !== originalOrderIds.value.length ||
    currentIds.some((id, i) => id !== originalOrderIds.value[i]);
  return (
    Number(settings.value.points) !==
      Number(data.points ?? constants.DefaultPoints) ||
    Number(settings.value.duration_in_seconds) !==
      Number(data.duration_in_seconds ?? constants.DefaultDurationInSeconds) ||
    settings.value.category_id !== nullableString(data.category_id) ||
    settings.value.cover_image !== nullableString(data.cover_image) ||
    orderChanged
  );
});

watch(
  () => quizData.value?.data,
  (data) => {
    if (!data) return;
    settings.value = {
      points: Number(data.points ?? constants.DefaultPoints),
      duration_in_seconds: Number(
        data.duration_in_seconds ?? constants.DefaultDurationInSeconds
      ),
      category_id: nullableString(data.category_id),
      cover_image: nullableString(data.cover_image),
    };
    // The server value is authoritative after a refresh, and we no longer
    // have the filename that produced it.
    coverImageName.value = "";
    const serverQuestions = data.data || [];
    orderedQuestions.value = [...serverQuestions];
    originalOrderIds.value = serverQuestions.map((q) => q.question_id);
  }
);

const parseAnswers = (value) => {
  if (Array.isArray(value)) return value;
  if (!value) return [];
  try {
    const parsed = JSON.parse(value);
    return Array.isArray(parsed) ? parsed : [];
  } catch {
    return [];
  }
};

const optionLabel = (key) => String.fromCharCode(64 + Number(key));

const isCorrectAnswer = (question, key) =>
  parseAnswers(question.correct_answer).includes(Number(key));

const openNewQuestionForm = () => {
  showAddAnother.value = false;
  showNewQuestionForm.value = true;
  editingQuestionId.value = "";
};

const startEditQuestion = (question) => {
  showNewQuestionForm.value = false;
  showAddAnother.value = false;
  editingQuestionId.value = question.question_id;
};

const handleCoverImage = async (event) => {
  processingCover.value = true;
  const picked = await pickCoverImage(event).finally(() => {
    processingCover.value = false;
  });
  if (!picked) return;
  settings.value.cover_image = picked.dataUrl;
  coverImageName.value = picked.name;
};

const removeCoverImage = () => {
  settings.value.cover_image = "";
  coverImageName.value = "";
};

const saveSettings = async () => {
  if (!canEditQuiz.value) return;

  if (settingsError.value) {
    toast.error(settingsError.value);
    return;
  }

  const body = {
    points: Number(settings.value.points),
    duration_in_seconds: Number(settings.value.duration_in_seconds),
    question_ids: orderedQuestions.value.map((q) => q.question_id),
  };

  // Only send category/cover when they actually changed — omitting them leaves
  // the columns untouched, which keeps a points tweak from re-uploading the
  // whole base64 cover image (and from tripping the admin-only check).
  if (canEditPublicMeta.value) {
    const server = quizData.value?.data || {};
    if (settings.value.category_id !== nullableString(server.category_id)) {
      body.category_id = settings.value.category_id;
    }
    if (settings.value.cover_image !== nullableString(server.cover_image)) {
      body.cover_image = settings.value.cover_image;
    }
  }

  try {
    settingsPending.value = true;
    await $fetch(`${url.apiUrl}/quizzes/${quizId.value}/settings`, {
      method: "PUT",
      headers,
      credentials: "include",
      body,
    });
    toast.success("Nastavení kvízu bylo uloženo.");
    refresh();
  } catch (error) {
    toast.error(readApiError(error, "Nastavení kvízu se nepodařilo uložit."));
  } finally {
    settingsPending.value = false;
  }
};

const saveNewQuestion = async ({ payload }) => {
  try {
    savingNewQuestion.value = true;
    const payloadWithDefaults = {
      ...payload,
      points: Number(settings.value.points),
      duration_in_seconds: Number(settings.value.duration_in_seconds),
    };
    await $fetch(`${url.apiUrl}/quizzes/${quizId.value}/questions`, {
      method: "POST",
      headers,
      credentials: "include",
      body: payloadWithDefaults,
    });

    toast.success("Otázka byla přidána.");
    showNewQuestionForm.value = false;
    showAddAnother.value = true;
    refresh();
  } catch (error) {
    toast.error(readApiError(error, "Otázku se nepodařilo přidat."));
  } finally {
    savingNewQuestion.value = false;
  }
};

const saveExistingQuestion = async (question, { payload }) => {
  const questionId = question.question_id;

  try {
    savingQuestionId.value = questionId;
    await $fetch(
      `${url.apiUrl}/quizzes/${quizId.value}/questions/${questionId}`,
      {
        method: "PUT",
        headers,
        credentials: "include",
        body: {
          ...payload,
          points: Number(settings.value.points),
          duration_in_seconds: Number(settings.value.duration_in_seconds),
        },
      }
    );

    toast.success("Otázka byla uložena.");
    editingQuestionId.value = "";
    refresh();
  } catch (error) {
    toast.error(readApiError(error, "Otázku se nepodařilo uložit."));
  } finally {
    savingQuestionId.value = "";
  }
};

const deleteQuestion = async (questionId) => {
  if (!window.confirm("Opravdu chceš smazat tuto otázku?")) return;

  try {
    await $fetch(
      `${url.apiUrl}/quizzes/${quizId.value}/questions/${questionId}`,
      {
        method: "DELETE",
        headers,
        credentials: "include",
      }
    );
    toast.success("Otázka byla smazána.");
    refresh();
  } catch (error) {
    toast.error(readApiError(error, "Otázku se nepodařilo smazat."));
  }
};

const deleteQuiz = async () => {
  try {
    await $fetch(`${url.apiUrl}/quizzes/${quizId.value}`, {
      method: "DELETE",
      headers,
      credentials: "include",
    });
    toast.success("Kvíz byl smazán.");
    router.push("/admin/quiz/list-quiz");
  } catch (error) {
    toast.error(readApiError(error, "Kvíz se nepodařilo smazat."));
  }
};

const handleCsvFile = (event) => {
  const file = event.target.files?.[0];
  csvFile.value = file || null;
  csvFileName.value = file?.name || "";
};

const closeImportModal = () => {
  importModalOpen.value = false;
  csvFile.value = null;
  csvFileName.value = "";
};

const importCsv = async () => {
  if (!csvFile.value) {
    toast.error("Vyber soubor CSV.");
    return;
  }

  try {
    importPending.value = true;
    const formData = new FormData();
    formData.append("attachment", csvFile.value);

    await $fetch(`${url.apiUrl}/quizzes/${quizId.value}/questions/upload`, {
      method: "POST",
      headers: {
        Accept: "application/json",
      },
      body: formData,
      credentials: "include",
    });

    toast.success("Otázky byly naimportovány.");
    closeImportModal();
    refresh();
  } catch (error) {
    toast.error(readApiError(error, "CSV se nepodařilo naimportovat."));
  } finally {
    importPending.value = false;
  }
};

const handleShareQuiz = () => {
  shareModalOpen.value = true;
};

const askStartMode = () => {
  if (questions.value.length === 0) {
    toast.error("Před spuštěním kvízu přidej aspoň jednu otázku.");
    return;
  }
  startModeOpen.value = true;
};

const handleStartQuiz = async (mode) => {
  try {
    startingQuiz.value = true;
    const response = await $fetch(
      `${url.apiUrl}/quizzes/${quizId.value}/demo_session?mode=${mode}`,
      {
        method: "POST",
        credentials: "include",
        headers: {
          Accept: "application/json",
        },
      }
    );

    const activeQuizId = response?.data;
    if (!activeQuizId) {
      toast.error("Kvíz se nepodařilo spustit.");
      return;
    }

    if (quizTitle.value) {
      sessionStore.setActiveQuizTitle(quizTitle.value);
    }

    listUserStore.removeAllUsers();
    setSocketObject(null);
    sessionStore.setSession(activeQuizId);
    router.push(`/admin/arrange/${activeQuizId}`);
  } catch (error) {
    toast.error(readApiError(error, "Kvíz se nepodařilo spustit."));
  } finally {
    startingQuiz.value = false;
  }
};
</script>

<style scoped>
.flip-list-move {
  transition: transform 0.4s ease;
}

/* Drop indicator: faded preview of the card in its drop position. */
.question-ghost {
  position: relative;
  opacity: 0.4 !important;
  background-color: var(--jv-white) !important;
  border: 2px dashed var(--jv-violet) !important;
  box-shadow: none !important;
  transform: none !important;
  transition: opacity 0.15s ease, border-color 0.15s ease;
}

/* Disable the hover lift on the ghost while it's acting as drop placeholder. */
.question-ghost:hover {
  transform: none !important;
}

/* Floating clone: the faded preview that follows the cursor. */
.question-drag {
  opacity: 0.92 !important;
  transform: scale(1.02) !important;
  box-shadow: var(--shadow-soft-lg) !important;
  cursor: grabbing !important;
}

/* Origin item: subtly dimmed so users see where it came from. */
.question-chosen {
  cursor: grabbing;
}
</style>
