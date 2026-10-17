<script setup>
import { useRouter } from "nuxt/app";
import { usePush } from "notivue";
import { Volume2, VolumeX } from "lucide-vue-next";
import { useMusicStore } from "~~/store/music";

const musicStore = useMusicStore();
const { getMusic, setMusic } = musicStore;

const props = defineProps({
  pageTitle: {
    type: String,
    required: true,
    default: "title",
  },
  pageMessage: {
    type: String,
    required: false,
    default: null,
  },
  musicComponent: {
    type: Boolean,
    required: false,
    default: false,
  },
});

const router = useRouter();
const toast = usePush();

const toastError = () => {
  const errorQueryParam = router.currentRoute.value.query?.error;
  if (errorQueryParam) {
    toast.error(errorQueryParam);
  }
};

onMounted(() => {
  if (process.client) {
    toastError();
  }
});

watch(
  () => router.currentRoute.value.query,
  () => toastError()
);

const music = computed(() => getMusic());
</script>

<template>
  <div class="flex min-h-[60vh] justify-center px-4 py-6 sm:px-6 sm:py-10">
    <div
      class="jv-border-rough relative w-full max-w-[800px] bg-jv-white p-5 shadow-soft-lg sm:p-8"
    >
      <div class="flex flex-row items-start gap-4">
        <div class="flex-1">
          <h1
            class="font-headings text-[26px] leading-tight text-jv-ink sm:text-[34px]"
          >
            {{ props.pageTitle }}
          </h1>
          <p
            v-if="props.pageMessage"
            class="mt-2 font-body text-[14px] font-bold text-jv-muted sm:text-[15px]"
          >
            {{ props.pageMessage }}
          </p>
        </div>
        <div>
          <slot name="sub-title"></slot>
        </div>
      </div>

      <hr class="my-4 border-t border-dashed border-jv-line" />

      <slot></slot>

      <div v-if="props.musicComponent" class="mt-4 flex justify-end">
        <button
          type="button"
          class="grid size-10 place-items-center rounded-[8px] border border-jv-line bg-jv-yellow text-jv-ink shadow-soft-sm transition-transform"
          :aria-label="music ? 'Vypnout hudbu' : 'Zapnout hudbu'"
          @click="setMusic(!music)"
        >
          <Volume2 v-if="music" class="size-5" :stroke-width="2.4" />
          <VolumeX v-else class="size-5" :stroke-width="2.4" />
        </button>
      </div>
    </div>
  </div>
</template>
