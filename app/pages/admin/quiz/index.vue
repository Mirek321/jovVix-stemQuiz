<template>
  <div class="min-h-screen bg-jv-canvas px-4 py-8 sm:px-6 sm:py-10">
    <div class="mx-auto max-w-6xl">
      <header class="mb-8">
        <h1
          class="font-headings text-[28px] leading-tight text-jv-ink sm:text-[36px]"
        >
          Kvízy
        </h1>
        <p class="mt-2 font-body text-sm text-jv-muted sm:text-base">
          Co chceš dnes dělat?
        </p>
      </header>

      <div
        class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 lg:gap-8"
      >
        <ActionCard
          action-title="Vytvořit kvíz"
          @card-click="navigateTo('/admin/quiz/list-quiz?create=1')"
        >
          <template #image>
            <NuxtImg
              src="/images/add-file.webp"
              alt=""
              width="256"
              height="256"
              loading="lazy"
              class="max-h-full max-w-full object-contain"
            />
          </template>
        </ActionCard>
        <ActionCard
          action-title="Moje kvízy"
          @card-click="navigateTo('/admin/quiz/list-quiz')"
        >
          <template #image>
            <NuxtImg
              src="/images/list.webp"
              alt=""
              width="256"
              height="256"
              loading="lazy"
              class="max-h-full max-w-full object-contain"
            />
          </template>
        </ActionCard>
        <ActionCard
          action-title="Sdílené kvízy"
          @card-click="navigateTo('/admin/quiz/shared-quiz')"
        >
          <template #image>
            <NuxtImg
              src="/images/shared-file.webp"
              alt=""
              width="256"
              height="256"
              loading="lazy"
              class="max-h-full max-w-full object-contain"
            />
          </template>
        </ActionCard>
        <ActionCard
          action-title="Sdíleno se mnou"
          @card-click="navigateTo('/admin/quiz/share-with-me')"
        >
          <template #image>
            <NuxtImg
              src="/images/share-with-me.webp"
              alt=""
              width="256"
              height="256"
              loading="lazy"
              class="max-h-full max-w-full object-contain"
            />
          </template>
        </ActionCard>
      </div>
    </div>
  </div>
</template>

<script setup>
import ActionCard from "~~/components/Quiz/ActionCard.vue";

definePageMeta({
  layout: "empty",
});

useSeoMeta({
  title: "Moje kvízy – STEM_HUB",
  description:
    "Vytvářej, upravuj a spravuj své interaktivní kvízy v systému STEM_HUB.",
  robots: "noindex, nofollow",
});
</script>
