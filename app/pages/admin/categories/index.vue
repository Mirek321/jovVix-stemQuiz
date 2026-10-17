<script setup>
import { computed, onMounted, ref } from "vue";
import { Check, Pencil, Search, Tag, Trash2, X } from "lucide-vue-next";
import { usePush } from "notivue";
import NavigationLink from "@/components/common/NavigationLink.vue";
import { useUsersStore } from "~~/store/users";

definePageMeta({
  layout: "empty",
});

useSeoMeta({
  title: "Kategorie kvízů – STEM_HUB",
  description:
    "Spravuj kategorie, do kterých se v systému STEM_HUB řadí veřejné kvízy.",
  robots: "noindex, nofollow",
});

const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const toast = usePush();
const usersStore = useUsersStore();

const canCreatePublicQuiz = computed(
  () => !!usersStore.userData?.canCreatePublicQuiz
);

// The API enforces the public-quiz admin allowlist on every write; this guard
// just redirects non-admins away instead of showing them a broken page.
onMounted(async () => {
  if (!usersStore.userData) {
    await setUserDataStore();
  }
  if (!canCreatePublicQuiz.value) {
    navigateTo("/admin/quiz/list-quiz");
  }
});

const {
  data: categoriesData,
  pending: categoriesPending,
  error: categoriesError,
  refresh,
} = useFetch(`${url.apiUrl}/categories`, {
  method: "GET",
  headers: headers,
  credentials: "include",
});

const categories = computed(() => categoriesData.value?.data || []);

const searchQuery = ref("");
const creating = ref(false);
const editingId = ref("");
const editingName = ref("");
const savingEdit = ref(false);
const deletingId = ref("");

const formatDate = (iso) => {
  if (!iso) return "";
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) return "";
  return date.toLocaleDateString("cs-CZ", {
    day: "numeric",
    month: "numeric",
    year: "numeric",
  });
};

const filteredCategories = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();
  if (!query) return categories.value;

  return categories.value.filter((category) =>
    category.name.toLowerCase().includes(query)
  );
});

// The create option only shows when the searched name doesn't exist yet.
const hasExactMatch = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();
  if (!query) return true;
  return categories.value.some(
    (category) => category.name.toLowerCase() === query
  );
});

const handleCreateFromSearch = async () => {
  const name = searchQuery.value.trim();
  if (!name || creating.value) return;

  try {
    creating.value = true;
    await $fetch(`${url.apiUrl}/categories`, {
      method: "POST",
      body: { name },
      credentials: "include",
      headers: {
        Accept: "application/json",
      },
    });
    toast.success("Kategorie byla vytvořena.");
    searchQuery.value = "";
    refresh();
  } catch (error) {
    toast.error(
      error?.data?.message ||
        error?.message ||
        "Kategorii se nepodařilo vytvořit."
    );
  } finally {
    creating.value = false;
  }
};

const startEdit = (category) => {
  editingId.value = category.id;
  editingName.value = category.name;
};

const cancelEdit = () => {
  editingId.value = "";
  editingName.value = "";
};

const handleUpdateCategory = async (category) => {
  const name = editingName.value.trim();
  if (!name || savingEdit.value) return;
  if (name === category.name) {
    cancelEdit();
    return;
  }

  try {
    savingEdit.value = true;
    await $fetch(`${url.apiUrl}/categories/${category.id}`, {
      method: "PUT",
      body: { name },
      credentials: "include",
      headers: {
        Accept: "application/json",
      },
    });
    toast.success("Kategorie byla upravena.");
    cancelEdit();
    refresh();
  } catch (error) {
    toast.error(
      error?.data?.message ||
        error?.message ||
        "Kategorii se nepodařilo upravit."
    );
  } finally {
    savingEdit.value = false;
  }
};

const handleDeleteCategory = async (category) => {
  if (deletingId.value) return;

  try {
    deletingId.value = category.id;
    await $fetch(`${url.apiUrl}/categories/${category.id}`, {
      method: "DELETE",
      credentials: "include",
      headers: {
        Accept: "application/json",
      },
    });
    toast.success("Kategorie byla smazána.");
    refresh();
  } catch (error) {
    toast.error(
      error?.data?.message ||
        error?.message ||
        "Kategorii se nepodařilo smazat."
    );
  } finally {
    deletingId.value = "";
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
          Kategorie kvízů
        </h1>
        <div
          class="ml-1 mt-1 h-3 w-24 rounded-full border-b border-jv-yellow sm:ml-2 sm:w-28"
          aria-hidden="true"
        ></div>
      </div>
    </div>

    <div
      class="flex flex-col gap-3 sm:flex-row sm:items-center sm:gap-4 justify-between"
    >
      <label
        class="flex h-11 w-full min-w-0 items-center gap-2.5 jv-border-rough bg-jv-white px-3 shadow-soft-sm sm:h-12 sm:max-w-[448px] sm:gap-3 sm:px-4"
      >
        <Search class="size-5 shrink-0 text-jv-ink/40" :stroke-width="2.4" />
        <input
          v-model="searchQuery"
          type="search"
          maxlength="50"
          class="h-full min-w-0 flex-1 bg-transparent text-[15px] text-jv-ink outline-none placeholder:text-jv-ink/45 sm:text-[17px]"
          placeholder="Hledej nebo přidej kategorii…"
          @keyup.enter="!hasExactMatch && handleCreateFromSearch()"
        />
      </label>
      <NavigationLink
        :url-name="creating ? 'Vytvářím…' : 'Vytvořit kategorii'"
        class="w-full bg-jv-mint py-2 font-[500] sm:w-fit"
        :disabled="creating || hasExactMatch"
        @click="handleCreateFromSearch"
      />
    </div>

    <div
      v-if="categoriesPending"
      class="jv-border-rough bg-jv-white p-5 text-center text-[18px] font-semibold text-jv-muted shadow-soft-sm"
    >
      Načítání…
    </div>

    <div
      v-else-if="categoriesError"
      class="jv-border-rough bg-jv-white p-5 text-[18px] font-semibold text-jv-coral shadow-soft-sm"
    >
      {{ categoriesError.message }}
    </div>

    <section
      v-else-if="categories.length < 1 && !searchQuery.trim()"
      class="grid min-h-[320px] place-items-center px-2 text-center sm:min-h-[420px]"
    >
      <div class="max-w-md py-8">
        <h2
          class="font-headings text-[28px] leading-tight text-jv-ink sm:text-[34px]"
        >
          Zatím žádné kategorie
        </h2>
        <p class="mt-2 text-[17px] text-jv-muted">
          Napiš název do vyhledávání nahoře a vytvoř svou první kategorii.
        </p>
      </div>
    </section>

    <section
      v-else-if="filteredCategories.length < 1"
      class="grid min-h-[320px] place-items-center px-2 text-center sm:min-h-[420px]"
    >
      <div class="max-w-md py-8">
        <h2
          class="font-headings text-[28px] leading-tight text-jv-ink sm:text-[34px]"
        >
          Kategorie nenalezena
        </h2>
        <p class="mt-2 text-[17px] text-jv-muted">
          Kategorie „{{ searchQuery.trim() }}“ zatím neexistuje. Přidáš ji
          tlačítkem Vytvořit kategorii.
        </p>
      </div>
    </section>

    <div v-else class="flex flex-col gap-3">
      <div
        class="hidden lg:grid grid-cols-[2fr_1fr_180px] gap-4 px-5 py-2 text-[12px] font-black uppercase tracking-[0.1em] text-jv-muted"
      >
        <div>Název kategorie</div>
        <div class="text-center">Vytvořeno</div>
        <div class="text-right">Akce</div>
      </div>

      <div
        v-for="category in filteredCategories"
        :key="category.id"
        class="grid grid-cols-1 lg:grid-cols-[2fr_1fr_180px] gap-4 items-center jv-border-rough bg-jv-white px-5 py-4 shadow-soft-sm hover:translate-y-[-2px] hover:shadow-soft-md transition-all"
      >
        <template v-if="editingId === category.id">
          <div class="flex min-w-0 items-center gap-3">
            <span
              class="grid size-10 shrink-0 place-items-center border border-jv-line bg-jv-yellow/60 text-jv-ink"
            >
              <Tag class="size-5" :stroke-width="2.4" />
            </span>
            <input
              v-model.trim="editingName"
              type="text"
              required
              maxlength="50"
              class="h-10 min-w-0 flex-1 border border-jv-line bg-jv-canvas px-3 text-[16px] font-semibold text-jv-ink outline-none"
              @keyup.enter="handleUpdateCategory(category)"
              @keyup.esc="cancelEdit"
            />
          </div>
          <div class="hidden lg:block"></div>
          <div class="flex items-center gap-2 lg:justify-end">
            <button
              type="button"
              :disabled="savingEdit || !editingName"
              class="grid size-9 shrink-0 place-items-center border border-jv-line bg-jv-mint text-jv-ink shadow-soft-sm transition-transform disabled:cursor-not-allowed disabled:opacity-60"
              aria-label="Uložit název kategorie"
              @click="handleUpdateCategory(category)"
            >
              <Check class="size-4" :stroke-width="2.5" />
            </button>
            <button
              type="button"
              class="grid size-9 shrink-0 place-items-center border border-jv-line bg-jv-white text-jv-ink shadow-soft-sm transition-transform"
              aria-label="Zrušit přejmenování"
              @click="cancelEdit"
            >
              <X class="size-4" :stroke-width="2.5" />
            </button>
          </div>
        </template>

        <template v-else>
          <div class="flex min-w-0 items-center gap-3">
            <span
              class="grid size-10 shrink-0 place-items-center border border-jv-line bg-jv-yellow/60 text-jv-ink"
            >
              <Tag class="size-5" :stroke-width="2.4" />
            </span>
            <h3
              class="min-w-0 truncate font-headings text-[22px] font-bold leading-tight text-jv-ink sm:text-[24px]"
            >
              {{ category.name }}
            </h3>
          </div>

          <div class="lg:text-center">
            <span
              class="lg:hidden mr-2 text-[14px] font-bold uppercase tracking-wide text-jv-muted"
            >
              Vytvořeno:
            </span>
            <span class="text-[15px] font-bold text-jv-ink">
              {{ formatDate(category.created_at) }}
            </span>
          </div>

          <div class="flex items-center gap-2 lg:justify-end">
            <button
              type="button"
              class="grid size-9 shrink-0 place-items-center border border-jv-line bg-jv-white text-jv-ink shadow-soft-sm transition-transform"
              aria-label="Přejmenovat kategorii"
              @click="startEdit(category)"
            >
              <Pencil class="size-4" :stroke-width="2.4" />
            </button>
            <button
              type="button"
              :disabled="deletingId === category.id"
              class="grid size-9 shrink-0 place-items-center border border-jv-line bg-jv-white text-jv-ink shadow-soft-sm transition-all hover:text-jv-coral disabled:cursor-wait disabled:opacity-60"
              aria-label="Smazat kategorii"
              @click="handleDeleteCategory(category)"
            >
              <Trash2 class="size-4" :stroke-width="2.4" />
            </button>
          </div>
        </template>
      </div>
    </div>
  </main>
</template>
