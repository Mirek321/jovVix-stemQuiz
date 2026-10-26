<template>
  <article
    class="relative min-h-[64px] sm:min-h-[72px] bg-jv-white px-3 sm:px-4 md:px-5 xl:px-7 py-3 text-center shadow-soft-sm jv-border-rough"
  >
    <h3
      class="font-body text-[15px] sm:text-[17px] md:text-[19px] xl:text-[20px] font-black leading-tight text-jv-ink"
    >
      {{ title }}
    </h3>
    <p
      class="mt-1 text-[11px] sm:text-[14px] leading-tight text-jv-muted break-words"
    >
      <a
        v-if="href"
        :href="href"
        class="font-semibold text-jv-ink underline decoration-jv-orange underline-offset-2"
        >{{ description }}</a
      >
      <template v-else>{{ description }}</template>
    </p>
  </article>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true,
  },
  description: {
    type: String,
    required: true,
  },
  href: {
    type: String,
    default: "",
  },
});
</script>
