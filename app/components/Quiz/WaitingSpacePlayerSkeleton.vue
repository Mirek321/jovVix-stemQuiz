<script setup>
import { useRouter } from "nuxt/app";
import { Hourglass, LogOut, Volume2, VolumeX } from "lucide-vue-next";
import { useMusicStore } from "~~/store/music";
import { Skeleton } from "@/components/ui/skeleton";

const router = useRouter();
const musicStore = useMusicStore();
const { getMusic, setMusic } = musicStore;

const music = computed(() => getMusic());
const toggleMusic = () => setMusic(!music.value);
const leaveLobby = () => router.push("/join");
</script>

<template>
  <main
    class="flex min-h-screen flex-col bg-jv-canvas text-jv-ink"
    aria-busy="true"
    aria-live="polite"
  >
    <header
      class="flex flex-wrap items-center justify-between gap-3 px-4 py-4 sm:gap-4 sm:px-6 sm:py-5 md:px-10"
    >
      <div class="flex min-w-0 items-center gap-3 sm:gap-4">
        <div class="flex min-w-0 flex-col gap-1.5">
          <p
            class="font-body text-[10px] font-bold uppercase tracking-[0.12em] text-jv-muted sm:text-[12px]"
          >
            Live Quiz Session
          </p>
          <Skeleton
            class="h-[24px] w-[160px] rounded-[6px] bg-jv-ink/10 sm:h-[30px] sm:w-[220px]"
          />
        </div>
      </div>

      <div class="flex shrink-0 items-center gap-2 sm:gap-3">
        <button
          type="button"
          class="grid size-11 place-items-center rounded-[8px] border border-jv-line bg-jv-white text-jv-ink shadow-soft-sm transition-transform sm:size-12"
          :aria-label="music ? 'Mute music' : 'Unmute music'"
          @click="toggleMusic"
        >
          <Volume2 v-if="music" class="size-5" :stroke-width="2.4" />
          <VolumeX v-else class="size-5" :stroke-width="2.4" />
        </button>
        <button
          type="button"
          class="inline-flex h-11 items-center justify-center gap-2 rounded-[8px] border border-jv-line bg-jv-white px-4 font-body text-[15px] font-bold text-jv-ink shadow-soft-sm transition-transform sm:h-12 sm:px-5 sm:text-[18px]"
          @click="leaveLobby"
        >
          <LogOut class="size-4 sm:size-5" :stroke-width="2.4" />
          <span>Leave Lobby</span>
        </button>
      </div>
    </header>

    <div class="flex justify-center px-4 sm:px-6 md:px-10">
      <div
        class="inline-flex max-w-full items-center gap-3 jv-border-rough bg-jv-white px-3 py-2 shadow-soft-sm sm:gap-4 sm:px-5 sm:py-2.5"
      >
        <Skeleton
          class="size-10 shrink-0 rounded-full border border-jv-line bg-jv-mint/50 sm:size-12"
        />
        <div class="flex min-w-0 flex-col gap-1.5">
          <p
            class="font-body text-[10px] font-bold uppercase tracking-[0.12em] text-jv-muted sm:text-[12px]"
          >
            Playing as
          </p>
          <Skeleton
            class="h-[18px] w-[120px] rounded-[6px] bg-jv-ink/10 sm:h-[22px] sm:w-[160px]"
          />
        </div>
      </div>
    </div>

    <section
      class="flex flex-1 items-center justify-center px-4 py-6 sm:px-6 sm:py-10 md:px-10 md:py-12"
    >
      <div class="relative w-full max-w-[980px]">
        <div
          class="jv-border-rough bg-jv-white p-6 shadow-soft-lg sm:p-10 md:p-12"
        >
          <div class="border-b border-dashed border-jv-line pb-6 sm:pb-7">
            <h2
              class="font-headings text-[28px] leading-tight sm:text-[36px] md:text-[48px]"
            >
              Ready, Steady, Go!
            </h2>
            <p
              class="mt-1 font-body text-[14px] font-bold text-jv-muted sm:text-[18px]"
            >
              Connecting to lobby...
            </p>
          </div>

          <div
            class="flex flex-wrap items-center justify-center gap-4 px-2 py-12 text-center sm:gap-5 sm:py-20 md:py-24"
          >
            <Skeleton
              class="h-[36px] w-[260px] max-w-full rounded-[8px] bg-jv-ink/10 sm:h-[52px] sm:w-[420px] md:h-[68px] md:w-[520px]"
            />
            <span
              class="relative grid size-14 shrink-0 place-items-center rounded-full border border-jv-line bg-jv-mint shadow-soft sm:size-16"
            >
              <Hourglass
                class="size-6 animate-hourglass-flip sm:size-7"
                :stroke-width="2.4"
              />
              <span
                class="absolute inset-1 rounded-full border border-dashed border-jv-line"
                aria-hidden="true"
              ></span>
            </span>
          </div>
        </div>
      </div>
    </section>
    <span class="sr-only">Connecting to quiz lobby…</span>
  </main>
</template>
