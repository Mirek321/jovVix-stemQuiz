<template>
  <section class="grid items-start">
    <div class="flex flex-col gap-3 items-start w-full min-w-0">
      <div
        class="rounded-[8px] border border-jv-line bg-jv-yellow/30 px-2.5 py-1 sm:px-3 sm:py-1.5 text-[12px] sm:text-[14px] md:small-text font-semibold uppercase shadow-soft-sm"
      >
        OBJEVUJ, UČ SE, HRAJ!
      </div>

      <h1 class="main-heading text-jv-ink">
        <span class="block">Kvízy STEM_HUB</span>
        <span class="block"
          >pro naše
          <span
            class="inline-block jv-highlight mt-2 sm:mt-3 px-3 sm:px-5 md:px-6 py-1 sm:py-2 rounded-[4px]"
            >živé akce</span
          >
        </span>
      </h1>

      <p
        class="max-w-[640px] text-[16px] leading-relaxed text-jv-muted sm:text-[18px]"
      >
        Náš systém na kvízy na akcích STEM Hubu. Víc o nás najdeš na
        <a
          href="https://stemhub.cz"
          target="_blank"
          rel="noopener noreferrer"
          class="font-semibold text-jv-ink underline decoration-jv-orange decoration-2 underline-offset-4 hover:decoration-jv-ink"
          >stemhub.cz</a
        >.
      </p>

      <!-- Mobile/sm carousel (under md) -->
      <Carousel
        :opts="{ align: 'start', loop: false }"
        class="group mt-4 sm:mt-6 w-full md:hidden"
      >
        <CarouselContent class="-ml-4 py-2">
          <CarouselItem
            v-for="card in actionCards"
            :key="card.title"
            class="pl-4 basis-[80%] min-[420px]:basis-1/2 sm:basis-1/3"
          >
            <HeroActionCard
              :title="card.title"
              :description="card.description"
              :href="card.href"
            />
          </CarouselItem>
        </CarouselContent>
        <CarouselPrevious class="jv-carousel-nav !left-2" />
        <CarouselNext class="jv-carousel-nav !right-2" />
      </Carousel>

      <!-- md+ static grid -->
      <div
        class="hidden md:grid mt-4 sm:mt-6 gap-5 sm:gap-6 md:gap-8 lg:gap-10 grid-cols-3 w-full"
      >
        <HeroActionCard
          v-for="card in actionCards"
          :key="card.title"
          :title="card.title"
          :description="card.description"
          :href="card.href"
        />
      </div>
    </div>
  </section>
</template>

<script setup>
import HeroActionCard from "@/components/landing/HeroActionCard.vue";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel";

const actionCards = [
  {
    title: "Připoj se kódem",
    description: "Kód kvízu dostaneš přímo na akci.",
  },
  {
    title: "Soutěž naživo",
    description: "Body za správné odpovědi i rychlost.",
  },
  {
    title: "Napiš nám",
    description: "info@stemhub.cz",
    href: "mailto:info@stemhub.cz",
  },
];
</script>
