<template>
  <form
    class="jv-border-rough bg-jv-white p-4 shadow-soft-sm sm:p-5 lg:p-6"
    @submit.prevent="submitForm"
  >
    <div
      class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
    >
      <div>
        <p
          class="text-[12px] font-bold uppercase tracking-[0.14em] text-jv-coral"
        >
          {{ eyebrow }}
        </p>
        <h2
          class="mt-1 font-body text-[22px] font-bold leading-tight text-jv-ink sm:text-[26px]"
        >
          {{ title }}
        </h2>
      </div>

      <select
        v-model.number="form.type"
        class="h-10 rounded-md border border-jv-line bg-jv-white px-3 text-[14px] font-semibold text-jv-ink outline-none focus:border-jv-line"
      >
        <option :value="1">Jedna správná odpověď</option>
        <option :value="2">Anketní otázka</option>
      </select>
    </div>

    <div class="mt-6 space-y-2">
      <label
        class="text-[13px] font-semibold uppercase tracking-wide text-jv-muted"
      >
        Otázka
      </label>
      <input
        v-model.trim="form.question"
        type="text"
        required
        placeholder="Zadej otázku…"
        class="h-11 w-full rounded-md border border-jv-line bg-jv-white px-3 text-[15px] font-medium text-jv-ink outline-none placeholder:text-jv-ink/40 focus:border-jv-line focus:shadow-soft-sm"
      />

      <div class="flex flex-wrap items-center gap-1.5 pt-1">
        <NavigationLink
          v-for="choice in mediaChoices"
          :key="`question-${choice.value}`"
          variant="toggle"
          type="button"
          :class="
            form.question_media === choice.value
              ? 'border-jv-line bg-jv-ink text-jv-white'
              : 'border-jv-line bg-jv-white text-jv-muted hover:border-jv-line hover:text-jv-ink'
          "
          @click="form.question_media = choice.value"
        >
          <component :is="choice.icon" class="size-3.5" :stroke-width="2.3" />
          <span>{{ choice.label }}</span>
        </NavigationLink>
      </div>

      <textarea
        v-if="form.question_media === 'code'"
        v-model="form.resource"
        rows="5"
        placeholder="Zadej kód…"
        class="mt-1 w-full resize-none rounded-md border border-jv-line bg-jv-canvas px-3 py-2 font-mono text-[14px] text-jv-ink outline-none focus:border-jv-line focus:shadow-soft-sm"
      ></textarea>

      <label
        v-else-if="form.question_media === 'image'"
        class="mt-1 flex min-h-12 cursor-pointer items-center justify-center gap-2 rounded-md border border-dashed border-jv-line bg-jv-canvas px-3 py-3 text-[14px] font-medium text-jv-muted transition-colors hover:bg-jv-yellow/20"
      >
        <ImageIcon class="size-4" :stroke-width="2.3" />
        <span>{{
          processingImage
            ? "Zpracovávám obrázek…"
            : questionFileName || "Nahrát obrázek k otázce"
        }}</span>
        <input
          type="file"
          class="hidden"
          accept="image/*"
          :disabled="processingImage"
          @change="handleQuestionImage"
        />
      </label>
    </div>

    <div class="mt-6 space-y-3">
      <div class="flex flex-wrap items-center justify-between gap-2">
        <label
          class="text-[13px] font-semibold uppercase tracking-wide text-jv-muted"
        >
          Možnosti
        </label>
        <div class="flex flex-wrap items-center gap-1.5">
          <NavigationLink
            v-for="choice in mediaChoices"
            :key="`options-${choice.value}`"
            variant="toggle"
            type="button"
            :class="
              form.options_media === choice.value
                ? 'border-jv-line bg-jv-ink text-jv-white'
                : 'border-jv-line bg-jv-white text-jv-muted hover:border-jv-line hover:text-jv-ink'
            "
            @click="form.options_media = choice.value"
          >
            <component :is="choice.icon" class="size-3.5" :stroke-width="2.3" />
            <span>{{ choice.label }}</span>
          </NavigationLink>
        </div>
      </div>

      <div class="flex flex-col">
        <div
          v-for="key in optionKeys"
          :key="key"
          class="flex min-w-0 items-center gap-3 border-b border-jv-line py-2.5 pl-2 pr-1 last:border-b-0"
          :class="
            form.type === 1 && Number(key) === Number(correctAnswer)
              ? 'border-l-4 border-l-jv-accent-green bg-jv-accent-green/25 pl-1'
              : 'border-l-4 border-l-transparent'
          "
        >
          <input
            v-if="form.type === 1"
            v-model.number="correctAnswer"
            type="radio"
            :value="Number(key)"
            class="size-5 shrink-0 accent-jv-accent-green"
            :aria-label="`Správná odpověď: možnost ${optionLetter(key)}`"
          />

          <span class="w-6 shrink-0 text-[14px] font-bold text-jv-coral">
            {{ optionLetter(key) }}.
          </span>

          <input
            v-if="form.options_media === 'text'"
            v-model.trim="form.options[key]"
            type="text"
            :required="Number(key) <= 2"
            :placeholder="`Možnost ${optionLetter(key)}`"
            class="h-10 w-full min-w-0 flex-1 rounded-md border border-jv-line bg-jv-white px-3 text-[15px] font-medium text-jv-ink outline-none placeholder:text-jv-ink/40 focus:border-jv-line focus:shadow-soft-sm"
          />

          <textarea
            v-else-if="form.options_media === 'code'"
            v-model="form.options[key]"
            rows="2"
            :placeholder="`Kód pro možnost ${optionLetter(key)}`"
            class="min-w-0 flex-1 resize-none rounded-md border border-jv-line bg-jv-white px-3 py-2 font-mono text-[13px] text-jv-ink outline-none focus:border-jv-line focus:shadow-soft-sm"
          ></textarea>

          <label
            v-else
            class="flex h-10 min-w-0 flex-1 cursor-pointer items-center gap-2 rounded-md border border-dashed border-jv-line bg-jv-white px-3 text-[13px] font-medium text-jv-muted focus-within:shadow-soft-sm hover:bg-jv-yellow/10"
          >
            <ImageIcon class="size-4 shrink-0" :stroke-width="2.3" />
            <span class="truncate">{{
              optionFileNames[key] || existingImageLabel(key)
            }}</span>
            <input
              type="file"
              class="hidden"
              accept="image/*"
              :disabled="processingImage"
              @change="handleOptionImage($event, key)"
            />
          </label>
        </div>
      </div>
    </div>

    <div class="mt-6 flex flex-col-reverse gap-3 sm:flex-row sm:justify-end">
      <NavigationLink
        v-if="showCancel"
        class="h-11 bg-jv-white text-jv-ink"
        @click="$emit('cancel')"
      >
        Zrušit
      </NavigationLink>
      <NavigationLink class="h-11 bg-jv-accent-green text-white">
        {{ saving ? "Ukládání…" : submitLabel }}
      </NavigationLink>
    </div>
  </form>
</template>

<script setup>
import { computed, reactive, ref, watch } from "vue";
import { Code2, ImageIcon, Type } from "lucide-vue-next";
import { usePush } from "notivue";
import NavigationLink from "../common/NavigationLink.vue";

const { pickImageFile } = useImagePicker();

const props = defineProps({
  question: {
    type: Object,
    default: null,
  },
  mode: {
    type: String,
    default: "create",
  },
  saving: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["save", "cancel"]);

const mediaChoices = [
  { label: "Text", value: "text", icon: Type },
  { label: "Obrázek", value: "image", icon: ImageIcon },
  { label: "Kód", value: "code", icon: Code2 },
];

const form = reactive({
  question: "",
  type: 1,
  question_media: "text",
  options_media: "text",
  resource: "",
  options: {
    1: "",
    2: "",
    3: "",
    4: "",
  },
});

const correctAnswer = ref(1);
const questionFileName = ref("");
const optionFileNames = ref({});
const processingImage = ref(false);

const parseAnswers = (value) => {
  if (Array.isArray(value)) return value;
  if (!value) return [1];
  try {
    const parsed = JSON.parse(value);
    return Array.isArray(parsed) ? parsed : [1];
  } catch {
    return [1];
  }
};

const resetForm = () => {
  const question = props.question;
  form.question = question?.question || "";
  form.type = Number(question?.question_type_id || 1);
  form.question_media = question?.question_media || "text";
  form.options_media = question?.options_media || "text";
  form.resource = question?.resource || "";
  form.options = {
    1: question?.options?.["1"] || "",
    2: question?.options?.["2"] || "",
    3: question?.options?.["3"] || "",
    4: question?.options?.["4"] || "",
  };

  if (question?.options?.["5"]) {
    form.options[5] = question.options["5"];
  }

  const answers = parseAnswers(question?.correct_answer);
  correctAnswer.value = Number(answers[0] || 1);
  questionFileName.value = "";
  optionFileNames.value = {};
};

watch(
  () => props.question,
  () => resetForm(),
  { immediate: true, deep: true }
);

const optionKeys = computed(() =>
  Object.keys(form.options)
    .sort((a, b) => Number(a) - Number(b))
    .filter((key) => Number(key) <= 5)
);

const optionLetter = (key) => String.fromCharCode(64 + Number(key));

const title = computed(() =>
  props.mode === "edit" ? "Upravit otázku" : "Nová otázka"
);
const eyebrow = computed(() => (props.mode === "edit" ? "Otázka" : "Vytvořit"));
const submitLabel = computed(() =>
  props.mode === "edit" ? "Uložit změny" : "Přidat otázku"
);
const showCancel = computed(
  () => props.mode === "edit" || props.mode === "create"
);

const pickImage = async (event) => {
  const file = event.target.files?.[0];
  if (!file) return null;

  processingImage.value = true;
  const picked = await pickImageFile(file).finally(() => {
    processingImage.value = false;
  });

  if (!picked) event.target.value = "";
  return picked;
};

const handleQuestionImage = async (event) => {
  const picked = await pickImage(event);
  if (!picked) return;

  form.resource = picked.dataUrl;
  questionFileName.value = picked.name;
};

const handleOptionImage = async (event, key) => {
  const picked = await pickImage(event);
  if (!picked) return;

  form.options[key] = picked.dataUrl;
  optionFileNames.value = {
    ...optionFileNames.value,
    [key]: picked.name,
  };
};

const existingImageLabel = (key) => {
  if (form.options[key]) return `Obrázek možnosti ${optionLetter(key)}`;
  return `Nahrát obrázek pro možnost ${optionLetter(key)}`;
};

const toast = usePush();

const submitForm = () => {
  const filledKeys = optionKeys.value.filter((key) =>
    String(form.options[key] ?? "").trim()
  );

  if (!form.question.trim() && !form.resource) {
    toast.error("Zadej text otázky.");
    return;
  }
  if (filledKeys.length < 2) {
    toast.error("Vyplň aspoň dvě možnosti.");
    return;
  }
  if (form.type === 1 && !filledKeys.includes(String(correctAnswer.value))) {
    toast.error("Správná odpověď musí být jedna z vyplněných možností.");
    return;
  }

  // Blank options are dropped and the rest renumbered 1..n, so players never
  // see empty choices and the correct answer still points at the same text.
  const options = {};
  const renumbered = {};
  filledKeys.forEach((key, index) => {
    options[index + 1] = form.options[key];
    renumbered[key] = index + 1;
  });
  const answers =
    form.type === 1
      ? [renumbered[String(correctAnswer.value)]]
      : filledKeys.map((key) => renumbered[key]);

  emit("save", {
    payload: {
      question: form.question,
      type: Number(form.type),
      options,
      answers,
      question_media: form.question_media,
      options_media: form.options_media,
      resource: form.resource,
    },
  });
};
</script>
