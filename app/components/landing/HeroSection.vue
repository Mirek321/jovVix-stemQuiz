<template>
  <section class="grid items-start">
    <div class="flex flex-col gap-3 items-start w-full min-w-0">
      <div
        class="rounded-[8px] border border-jv-line bg-jv-yellow/30 px-2.5 py-1 sm:px-3 sm:py-1.5 text-[12px] sm:text-[14px] md:small-text font-semibold uppercase shadow-soft-sm"
      >
        OBJEVUJ, UČ SE, HRAJ!
      </div>

      <h1 class="main-heading text-jv-ink">
        <span class="block">Open-source kvízy</span>
        <span class="block"
          >stvořené pro
          <span
            class="inline-block jv-highlight mt-2 sm:mt-3 px-3 sm:px-5 md:px-6 py-1 sm:py-2 rounded-[4px]"
            >živé kvízové zážitky</span
          >
        </span>
      </h1>

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
    title: "Demo instance",
    description: "Až 100 hráčů najednou",
  },
  {
    title: "Vlastní nasazení",
    description: "Neomezený počet uživatelů a plná kontrola.",
  },
  {
    title: "Napiš nám",
    description: "contact@jovvix.com",
  },
];
</script>
