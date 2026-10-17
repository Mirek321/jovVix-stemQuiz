<template>
  <NuxtLink
    :to="props.navigateTo"
    class="jv-card inline-flex h-11 items-center justify-center gap-2 border border-jv-line bg-jv-coral px-6 font-headings text-base text-white shadow-soft-sm transition-transform"
  >
    {{ props.title }}
  </NuxtLink>
</template>

<script setup>
const props = defineProps({
  title: {
    type: String,
    required: true,
    default: "",
  },
  navigateTo: {
    type: String,
    required: true,
    default: "",
  },
});
</script>
