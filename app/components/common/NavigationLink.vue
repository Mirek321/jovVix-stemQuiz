<template>
  <NuxtLink
    v-if="url"
    :to="url"
    v-bind="$attrs"
    :class="cn(sharedClasses, variantClasses, props.class)"
  >
    <slot />
    <span v-if="urlName">{{ urlName }}</span>
  </NuxtLink>
  <button
    v-else
    v-bind="$attrs"
    :class="cn(sharedClasses, variantClasses, props.class)"
  >
    <slot />
    <span v-if="urlName">{{ urlName }}</span>
  </button>
</template>

<script setup>
import { computed } from "vue";
import { cn } from "@/lib/utils";

defineOptions({ inheritAttrs: false });

const props = defineProps({
  url: {
    type: [String, Object],
    default: "",
  },
  urlName: {
    type: String,
    default: "",
  },
  variant: {
    type: String,
    default: "pill",
  },
  class: {
    type: [String, Array, Object],
    default: "",
  },
});

const sharedClasses =
  "inline-flex items-center justify-center whitespace-nowrap disabled:opacity-60 disabled:cursor-not-allowed aria-disabled:opacity-60 aria-disabled:cursor-not-allowed aria-disabled:pointer-events-none";

const pillClasses =
  "font-headings font-semibold text-center jv-border-uneven shadow-soft-sm transition-all duration-150 hover:-translate-y-px hover:shadow-soft active:translate-y-0 active:scale-[0.98] disabled:hover:translate-y-0 bg-jv-yellow text-jv-ink gap-2 px-5 sm:px-6 py-2.5 text-sm md:text-[15px]";

const toggleClasses =
  "gap-1.5 rounded-md border px-2.5 h-8 text-[13px] font-semibold transition-colors";

const variantClasses = computed(() =>
  props.variant === "toggle" ? toggleClasses : pillClasses
);
</script>
