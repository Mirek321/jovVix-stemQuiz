<script setup>
import { computed, ref, watch, watchEffect } from "vue";
import { onClickOutside } from "@vueuse/core";
import {
  ChevronDown,
  Filter,
  Image as ImageIcon,
  Search,
  X,
} from "lucide-vue-next";
import { usePush } from "notivue";
import AdminQuizListCard from "@/components/quiz-list/AdminQuizListCard.vue";
import ShareQuizModal from "@/components/Quiz/ShareQuizModal.vue";
import StartModeModal from "@/components/Quiz/StartModeModal.vue";
import { useListUserstore } from "~/store/userlist";
import { useSessionStore } from "~~/store/session";
import { useUsersStore } from "~~/store/users";
import NavigationLink from "@/components/common/NavigationLink.vue";

definePageMeta({
  layout: "empty",
  middleware: async (to) => {
    const { apiUrl } = useRuntimeConfig().public;
    const requestHeaders = useRequestHeaders(["cookie"]);

    try {
      const who = await $fetch(`${apiUrl}/user/who`, {
        method: "GET",
        credentials: "include",
        headers: requestHeaders,
      });

      const role = who?.data?.role;
      if (role && role !== "guest-user") {
        return;
      }
    } catch (error) {
      console.log("unauthenticated visitor on the quiz list", error.message);
    }

    return navigateTo(
      `/account/login?returnTo=${encodeURIComponent(to.fullPath)}`
    );
  },
});

useSeoMeta({
  title: "Knihovna kvízů – STEM_HUB",
  description: "Procházej a spravuj všechny své kvízy v systému STEM_HUB.",
  robots: "noindex, nofollow",
});

const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const toast = usePush();
const router = useRouter();
const route = useRoute();
const { pickCoverImage } = useCoverImage();
const sessionStore = useSessionStore();
const listUserStore = useListUserstore();
const usersStore = useUsersStore();

const canCreatePublicQuiz = computed(
  () => !!usersStore.userData?.canCreatePublicQuiz
);

const searchQuery = ref("");
const startingQuizId = ref("");
const startModeQuizId = ref("");
const startModeOpen = ref(false);
const shareModalOpen = ref(false);
const shareQuizId = ref("");
const deleteModalOpen = ref(false);
const quizPendingDelete = ref(null);
const createQuizOpen = ref(false);
const createQuizPending = ref(false);
const coverImageName = ref("");
const processingCover = ref(false);
const createQuizForm = ref({
  title: "",
  description: "",
  is_public: false,
  category_id: "",
  cover_image: "",
});
const selectedFilter = ref("Všechny kvízy");
const filterOpen = ref(false);
const filterMenuRef = ref(null);
const filterOptions = [
  "Všechny kvízy",
  "Moje kvízy",
  "Sdílím já",
  "Sdíleno se mnou",
];
const visibilityFilter = ref("Soukromé");
const visibilityOptions = ["Soukromé", "Veřejné"];

const { coverUrl, fallbackUrl } = useQuizCover();

const emptyState = computed(() => {
  if (selectedFilter.value === "Moje kvízy") {
    return {
      title: "Zatím nemáš žádný kvíz",
      message: "Vytvoř svůj první kvíz.",
      showCreateAction: true,
    };
  }

  if (selectedFilter.value === "Sdílím já") {
    return {
      title: "Zatím nesdílíš žádný kvíz",
      message: "Nasdílej svůj první kvíz.",
      showCreateAction: false,
    };
  }

  if (selectedFilter.value === "Sdíleno se mnou") {
    return {
      title: "Zatím s tebou nikdo nesdílel žádný kvíz",
      message: "Požádej někoho, ať s tebou nasdílí kvíz.",
      showCreateAction: false,
    };
  }

  return {
    title: "Zatím žádné kvízy",
    message:
      "Vytvoř svůj první kvíz, nebo požádej někoho, ať s tebou nějaký nasdílí.",
    showCreateAction: true,
  };
});

const quizListOptions = {
  method: "GET",
  headers: headers,
  mode: "cors",
  credentials: "include",
};

const {
  data: ownQuizList,
  pending: ownPending,
  error: ownError,
  refresh: refreshOwn,
} = useFetch(`${url.apiUrl}/quizzes`, quizListOptions);

const {
  data: sharedByMeList,
  pending: sharedByMePending,
  error: sharedByMeError,
  refresh: refreshSharedByMe,
} = useFetch(`${url.apiUrl}/shared_quizzes?type=shared_by_me`, quizListOptions);

const {
  data: sharedWithMeList,
  pending: sharedWithMePending,
  error: sharedWithMeError,
  refresh: refreshSharedWithMe,
} = useFetch(
  `${url.apiUrl}/shared_quizzes?type=shared_with_me`,
  quizListOptions
);

const quizPending = computed(
  () => ownPending.value || sharedByMePending.value || sharedWithMePending.value
);

const quizError = computed(
  () => ownError.value || sharedByMeError.value || sharedWithMeError.value
);

const refresh = () =>
  Promise.all([refreshOwn(), refreshSharedByMe(), refreshSharedWithMe()]);

const { data: categoriesData } = useFetch(`${url.apiUrl}/categories`, {
  method: "GET",
  headers: headers,
  credentials: "include",
});
const categories = computed(() => categoriesData.value?.data || []);

// Category and cover image only apply to public quizzes; drop them if the
// admin unticks the checkbox after filling them in.
watch(
  () => createQuizForm.value.is_public,
  (isPublic) => {
    if (!isPublic) {
      createQuizForm.value.category_id = "";
      createQuizForm.value.cover_image = "";
      coverImageName.value = "";
    }
  }
);

watchEffect(() => {
  if (quizError.value?.data?.code === 401) {
    navigateTo(
      `/account/login?returnTo=${encodeURIComponent(route.fullPath)}`,
      { replace: true }
    );
  }
});

watchEffect(() => {
  if (route.query.create === "1") {
    createQuizOpen.value = true;
  }
});

const formatCreatedAt = (value) => {
  if (!value) return "Vytvořeno nedávno";

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "Vytvořeno nedávno";

  return `Vytvořeno ${new Intl.DateTimeFormat("cs", {
    month: "short",
    day: "2-digit",
    year: "numeric",
  }).format(date)}`;
};

const rawQuizzes = computed(() => {
  const own = ownQuizList.value?.data || [];
  const byMe = sharedByMeList.value?.data || [];
  const withMe = sharedWithMeList.value?.data || [];

  if (selectedFilter.value === "Moje kvízy") return own;
  if (selectedFilter.value === "Sdílím já") return byMe;
  if (selectedFilter.value === "Sdíleno se mnou") return withMe;

  const seen = new Set();

  return [...own, ...byMe, ...withMe]
    .filter((quiz) => !seen.has(quiz.id) && seen.add(quiz.id))
    .sort(
      (a, b) =>
        new Date(b.created_at) - new Date(a.created_at) ||
        a.id.localeCompare(b.id)
    );
});

const quizzes = computed(() =>
  rawQuizzes.value.map((quiz) => ({
    id: quiz.id,
    title: quiz.title,
    description: quiz.description?.String || "",
    createdAt: formatCreatedAt(quiz.created_at),
    questionCount: quiz.total_questions || 0,
    image: coverUrl(quiz),
    fallback: fallbackUrl(quiz),
    viewUrl: `/admin/quiz/list-quiz/${quiz.id}`,
    isPublic: !!quiz.is_public,
  }))
);

const visibilityCounts = computed(() => {
  const list = quizzes.value;
  const publicCount = list.filter((quiz) => quiz.isPublic).length;

  return {
    Soukromé: list.length - publicCount,
    Veřejné: publicCount,
  };
});

const filteredQuizzes = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();
  const isPublicWanted = visibilityFilter.value === "Veřejné";

  return quizzes.value.filter((quiz) => {
    if (quiz.isPublic !== isPublicWanted) return false;
    if (!query) return true;

    return (
      quiz.title.toLowerCase().includes(query) ||
      quiz.description.toLowerCase().includes(query)
    );
  });
});

const noResultsMessage = computed(() => {
  if (searchQuery.value.trim()) return "Hledání neodpovídá žádný kvíz.";

  return visibilityFilter.value === "Veřejné"
    ? "Zatím tu nejsou žádné veřejné kvízy."
    : "Zatím tu nejsou žádné soukromé kvízy.";
});

const askStartMode = (quizId) => {
  startModeQuizId.value = quizId;
  startModeOpen.value = true;
};

const handleStartQuiz = async (quizId, mode) => {
  try {
    startingQuizId.value = quizId;
    const response = await $fetch(
      `${url.apiUrl}/quizzes/${quizId}/demo_session?mode=${mode}`,
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

    const startedQuiz = quizzes.value.find((q) => q.id === quizId);
    if (startedQuiz?.title) {
      sessionStore.setActiveQuizTitle(startedQuiz.title);
    }

    listUserStore.removeAllUsers();
    setSocketObject(null);
    router.push(`/admin/arrange/${activeQuizId}`);

    setTimeout(() => {
      sessionStore.setSession(activeQuizId);
    }, 1000);
  } catch (error) {
    toast.error(error?.message || "Kvíz se nepodařilo spustit.");
  } finally {
    startingQuizId.value = "";
  }
};

const handleShareQuiz = (quiz) => {
  shareQuizId.value = quiz.id;
  shareModalOpen.value = true;
};

const deleteQuizMessage = computed(() => {
  const title = quizPendingDelete.value?.title;
  return `Smazáním ${
    title ? `kvízu „${title}“` : "tohoto kvízu"
  } odstraníš i všechny jeho otázky, všechny odehrané hry a skóre i odpovědi hráčů. Tuto akci nejde vrátit.`;
});

const handleDeleteQuiz = (quiz) => {
  quizPendingDelete.value = quiz;
  deleteModalOpen.value = true;
};

const confirmDeleteQuiz = async () => {
  const quizId = quizPendingDelete.value?.id;
  if (!quizId) return;

  try {
    await $fetch(`${url.apiUrl}/quizzes/${quizId}`, {
      method: "DELETE",
      headers: headers,
      credentials: "include",
    });
    toast.success("Kvíz byl smazán.");
    refresh();
  } catch (error) {
    console.error("Failed to delete quiz", error);
    toast.error(
      error?.data?.message || error?.message || "Kvíz se nepodařilo smazat."
    );
  } finally {
    quizPendingDelete.value = null;
  }
};

const toggleQuizVisibility = async (quiz) => {
  const isPublic = !quiz.isPublic;
  try {
    await $fetch(`${url.apiUrl}/quizzes/${quiz.id}/visibility`, {
      method: "PUT",
      headers: headers,
      credentials: "include",
      body: { is_public: isPublic },
    });
    toast.success(
      isPublic
        ? `Kvíz „${quiz.title}“ je teď veřejný.`
        : `Kvíz „${quiz.title}“ je teď soukromý.`
    );
    refresh();
  } catch (error) {
    console.error("Failed to change quiz visibility", error);
    toast.error(
      error?.data?.data ||
        error?.data?.message ||
        "Viditelnost kvízu se nepodařilo změnit."
    );
  }
};

const handleSelectFilter = (option) => {
  selectedFilter.value = option;
  filterOpen.value = false;
};

const closeFilter = () => {
  filterOpen.value = false;
};

onClickOutside(filterMenuRef, closeFilter);

const resetCreateQuizForm = () => {
  createQuizOpen.value = false;
  createQuizForm.value = {
    title: "",
    description: "",
    is_public: false,
    category_id: "",
    cover_image: "",
  };
  coverImageName.value = "";
};

const closeCreateQuizModal = () => {
  resetCreateQuizForm();
  if (route.query.create === "1") {
    router.replace({ path: route.path, query: {} });
  }
};

const handleCoverImage = async (event) => {
  processingCover.value = true;
  const picked = await pickCoverImage(event).finally(() => {
    processingCover.value = false;
  });
  if (!picked) return;
  createQuizForm.value.cover_image = picked.dataUrl;
  coverImageName.value = picked.name;
};

const removeCoverImage = () => {
  createQuizForm.value.cover_image = "";
  coverImageName.value = "";
};

const handleCreateQuiz = async () => {
  try {
    createQuizPending.value = true;
    const response = await $fetch(`${url.apiUrl}/quizzes`, {
      method: "POST",
      headers: {
        Accept: "application/json",
      },
      body: createQuizForm.value,
      credentials: "include",
    });

    const quizId = response?.data;
    if (!quizId) {
      toast.error("Kvíz se nepodařilo vytvořit.");
      return;
    }

    toast.success("Kvíz byl vytvořen.");
    // One navigation only: a separate replace() to strip ?create=1 would
    // cancel this push and leave the form open. Replacing the ?create=1
    // entry also keeps Back from reopening the form.
    resetCreateQuizForm();
    const target = `/admin/quiz/list-quiz/${quizId}`;
    if (route.query.create === "1") {
      router.replace(target);
    } else {
      router.push(target);
    }
  } catch (error) {
    toast.error(
      error?.data?.data ||
        error?.data?.message ||
        error?.message ||
        "Kvíz se nepodařilo vytvořit."
    );
  } finally {
    createQuizPending.value = false;
  }
};
</script>

<template>
  <main
    class="flex min-h-screen flex-col gap-8 bg-jv-canvas px-4 py-5 sm:gap-10 sm:px-6 md:px-8 md:py-6"
  >
    <div
      class="flex flex-col gap-5 md:flex-row md:items-start md:justify-between"
    >
      <div class="min-w-0">
        <h1
          class="font-headings text-[38px] leading-none text-jv-ink min-[420px]:text-[44px] sm:text-[52px] md:text-[56px]"
        >
          Seznam kvízů
        </h1>
        <div
          class="ml-1 mt-1 h-3 w-24 rounded-full border-b border-jv-yellow sm:ml-2 sm:w-28"
          aria-hidden="true"
        ></div>
      </div>
      <div class="flex w-full flex-col gap-3 sm:w-fit sm:flex-row md:mt-1">
        <NavigationLink
          url-name="Vytvořit kvíz"
          class="w-full bg-jv-coral py-2 font-[500] text-white sm:w-fit"
          @click="createQuizOpen = true"
        />
      </div>
    </div>

    <div
      class="grid gap-3 sm:gap-4 md:grid-cols-[minmax(0,1fr)_minmax(280px,448px)] md:items-center md:justify-between"
    >
      <div class="flex flex-wrap items-center gap-3">
        <div ref="filterMenuRef" class="relative w-full sm:w-fit">
          <button
            type="button"
            class="inline-flex h-11 w-full items-center justify-between gap-2 jv-border-rough bg-jv-white px-3 text-[16px] font-semibold text-jv-ink shadow-soft-sm sm:h-12 sm:w-fit sm:justify-start sm:px-4 sm:text-[18px]"
            :aria-expanded="filterOpen"
            aria-haspopup="listbox"
            @click="filterOpen = !filterOpen"
          >
            <span class="flex min-w-0 items-center gap-2">
              <Filter
                class="size-5 shrink-0 text-jv-coral"
                :stroke-width="2.2"
              />
              <span class="truncate">{{ selectedFilter }}</span>
            </span>
            <ChevronDown
              class="size-4 shrink-0 text-jv-ink/60 transition-transform"
              :class="filterOpen ? 'rotate-180' : ''"
              :stroke-width="2.4"
            />
          </button>

          <div
            v-if="filterOpen"
            class="absolute left-0 top-[52px] z-30 w-full min-w-[190px] jv-border-rough bg-jv-white p-2 shadow-soft-sm sm:top-14 sm:w-52"
            role="listbox"
          >
            <button
              v-for="option in filterOptions"
              :key="option"
              type="button"
              class="block w-full rounded-[6px] px-3 py-2 text-left text-[15px] font-semibold text-jv-ink transition-colors hover:bg-jv-yellow/40"
              :class="option === selectedFilter ? 'bg-jv-yellow/60' : ''"
              role="option"
              :aria-selected="option === selectedFilter"
              @click="handleSelectFilter(option)"
            >
              {{ option }}
            </button>
          </div>
        </div>

        <div
          class="flex h-11 w-full items-center gap-1 jv-border-rough bg-jv-white p-1 shadow-soft-sm sm:h-12 sm:w-fit"
          role="group"
          aria-label="Filtrovat podle viditelnosti"
        >
          <button
            v-for="option in visibilityOptions"
            :key="option"
            type="button"
            class="flex-1 rounded-[6px] px-3 py-1.5 text-[14px] font-semibold text-jv-ink transition-colors hover:bg-jv-yellow/40 sm:flex-none sm:text-[15px]"
            :class="option === visibilityFilter ? 'bg-jv-yellow/60' : ''"
            :aria-pressed="option === visibilityFilter"
            @click="visibilityFilter = option"
          >
            {{ option }}
            <span class="text-jv-ink/50">{{ visibilityCounts[option] }}</span>
          </button>
        </div>
      </div>

      <label
        class="flex h-11 min-w-0 items-center gap-2.5 jv-border-rough bg-jv-white px-3 shadow-soft-sm sm:h-12 sm:gap-3 sm:px-4"
      >
        <Search class="size-5 shrink-0 text-jv-ink/40" :stroke-width="2.4" />
        <input
          v-model="searchQuery"
          type="search"
          class="h-full min-w-0 flex-1 bg-transparent text-[15px] text-jv-ink outline-none placeholder:text-jv-ink/45 sm:text-[17px]"
          placeholder="Hledej podle názvu nebo klíčového slova…"
        />
      </label>
    </div>

    <UtilsQuizListWaiting v-if="quizPending" />

    <div
      v-else-if="quizError"
      class="jv-border-rough bg-jv-white p-5 text-[18px] font-semibold text-jv-coral shadow-soft-sm"
    >
      {{ quizError.message }}
    </div>

    <section
      v-else-if="quizzes.length < 1"
      class="grid min-h-[320px] place-items-center px-2 text-center sm:min-h-[420px]"
    >
      <div class="max-w-md py-8">
        <h2
          class="font-headings text-[28px] leading-tight text-jv-ink sm:text-[34px]"
        >
          {{ emptyState.title }}
        </h2>
        <p class="mt-2 text-[17px] text-jv-muted">
          {{ emptyState.message }}
        </p>
        <div
          v-if="emptyState.showCreateAction"
          class="mt-6 flex flex-wrap justify-center gap-3"
        >
          <NavigationLink
            url-name="Vytvořit kvíz"
            class="bg-jv-coral text-white font-[500] py-2"
            @click="createQuizOpen = true"
          />
        </div>
      </div>
    </section>

    <section
      v-else-if="filteredQuizzes.length < 1"
      class="jv-border-rough bg-jv-white p-5 text-[18px] font-semibold text-jv-muted shadow-soft-sm"
    >
      {{ noResultsMessage }}
    </section>

    <section
      v-else
      class="grid items-start justify-items-center gap-7 sm:grid-cols-2 sm:gap-8 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 xl:gap-7"
    >
      <AdminQuizListCard
        v-for="quiz in filteredQuizzes"
        :key="quiz.id"
        :title="quiz.title"
        :created-at="quiz.createdAt"
        :question-count="quiz.questionCount"
        :image="quiz.image"
        :fallback="quiz.fallback"
        :view-url="quiz.viewUrl"
        :starting="startingQuizId === quiz.id"
        :is-public="quiz.isPublic"
        :can-change-visibility="canCreatePublicQuiz"
        @start-quiz="askStartMode(quiz.id)"
        @share="handleShareQuiz(quiz)"
        @delete="handleDeleteQuiz(quiz)"
        @toggle-visibility="toggleQuizVisibility(quiz)"
      />
    </section>

    <Teleport to="body">
      <div
        v-if="createQuizOpen"
        class="fixed inset-0 z-50 grid place-items-center bg-jv-ink/35 px-4 py-6 backdrop-blur-[2px]"
        @click.self="closeCreateQuizModal"
      >
        <form
          class="w-full max-w-[680px] border-[4px] border-jv-line bg-jv-white shadow-soft-lg"
          @submit.prevent="handleCreateQuiz"
        >
          <div
            class="flex items-center justify-between gap-4 border-b border-jv-line bg-jv-ink px-5 py-4 text-jv-white sm:px-6"
          >
            <h2
              class="font-body text-[24px] font-black leading-none text-jv-white sm:text-[28px]"
            >
              Vytvořit kvíz
            </h2>
            <button
              type="button"
              class="grid size-9 place-items-center text-jv-white transition-transform"
              aria-label="Zavřít okno pro vytvoření kvízu"
              @click="closeCreateQuizModal"
            >
              <X class="size-6" :stroke-width="2.4" />
            </button>
          </div>

          <div class="grid gap-5 px-5 py-6 sm:px-8">
            <label class="grid gap-2">
              <span
                class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink"
              >
                Název kvízu <span class="text-jv-coral">*</span>
              </span>
              <input
                v-model.trim="createQuizForm.title"
                type="text"
                class="h-14 border border-jv-line bg-jv-canvas px-4 text-[17px] font-semibold text-jv-ink outline-none transition-shadow focus:shadow-soft-sm"
              />
            </label>

            <label class="grid gap-2">
              <span
                class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink"
              >
                Popis kvízu
              </span>
              <textarea
                v-model.trim="createQuizForm.description"
                maxlength="150"
                rows="5"
                class="resize-none border border-jv-line bg-jv-canvas px-4 py-3 text-[17px] font-semibold text-jv-ink outline-none transition-shadow focus:shadow-soft-sm"
              ></textarea>
            </label>

            <label
              v-if="canCreatePublicQuiz"
              class="flex cursor-pointer items-start gap-3 jv-border-rough bg-jv-canvas p-4"
            >
              <input
                v-model="createQuizForm.is_public"
                type="checkbox"
                class="mt-1 size-5 shrink-0 cursor-pointer accent-jv-coral"
              />
              <span class="flex flex-col gap-1">
                <span
                  class="text-[14px] font-black uppercase tracking-[0.14em] text-jv-ink"
                >
                  Zveřejnit tento kvíz
                </span>
                <span class="text-[14px] font-medium text-jv-muted">
                  Veřejné kvízy se zobrazují na domovské stránce mezi veřejnými
                  kvízy. Najít je může kdokoli.
                </span>
              </span>
            </label>

            <div
              v-if="canCreatePublicQuiz && createQuizForm.is_public"
              class="grid gap-5 sm:grid-cols-2 sm:gap-x-8"
            >
              <label class="grid content-start gap-2">
                <span
                  class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink"
                >
                  Kategorie
                </span>
                <select
                  v-model="createQuizForm.category_id"
                  class="h-14 border border-jv-line bg-jv-canvas px-4 text-[16px] font-semibold text-jv-ink outline-none transition-shadow focus:shadow-soft-sm"
                >
                  <option value="">
                    Bez kategorie (zobrazí se v „Ostatní“)
                  </option>
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
                  v-if="!createQuizForm.cover_image"
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
                    :disabled="processingCover"
                    @change="handleCoverImage"
                  />
                </label>
                <div
                  v-else
                  class="flex h-14 items-center gap-3 border border-jv-line bg-jv-canvas px-3"
                >
                  <NuxtImg
                    :src="createQuizForm.cover_image"
                    alt="Náhled titulního obrázku"
                    class="h-10 w-14 shrink-0 border border-jv-line object-cover"
                  />
                  <span
                    class="min-w-0 flex-1 truncate text-[14px] font-semibold text-jv-ink"
                  >
                    {{ coverImageName }}
                  </span>
                  <button
                    type="button"
                    class="grid size-8 shrink-0 place-items-center text-jv-ink transition-colors hover:text-jv-coral"
                    aria-label="Odebrat titulní obrázek"
                    @click="removeCoverImage"
                  >
                    <X class="size-4" :stroke-width="2.4" />
                  </button>
                </div>
              </div>
            </div>
          </div>

          <div
            class="flex flex-col-reverse gap-3 border-t border-jv-line bg-jv-canvas px-5 py-4 sm:flex-row sm:justify-end sm:px-8"
          >
            <NavigationLink
              url-name="Zrušit"
              class="bg-jv-white font-[500]"
              @click="closeCreateQuizModal"
            />
            <NavigationLink
              url-name="Vytvořit kvíz"
              class="bg-jv-mint font-[500]"
            />
          </div>
        </form>
      </div>
    </Teleport>

    <ShareQuizModal v-model="shareModalOpen" :quiz-id="shareQuizId" />
    <StartModeModal
      v-model="startModeOpen"
      @select="(mode) => handleStartQuiz(startModeQuizId, mode)"
    />
    <DeleteDialog
      v-model="deleteModalOpen"
      title="Smazat kvíz"
      :message="deleteQuizMessage"
      confirm-label="Smazat kvíz"
      @confirm-delete="confirmDeleteQuiz"
    />
  </main>
</template>
