<script setup>
import { Crown, Medal, Award } from "lucide-vue-next";
import { getAvatarUrlByName } from "~~/composables/avatar";

const props = defineProps({
  winner: {
    type: Object,
    required: true,
    default: () => ({}),
  },
});

const avatar = computed(() => getAvatarUrlByName(props.winner?.img_key));

const rankConfig = computed(() => {
  const rank = Number(props.winner?.rank);
  if (rank === 1) {
    return {
      bg: "bg-jv-yellow",
      icon: Crown,
      label: "1.",
      rotate: "",
      scale: "sm:scale-110",
    };
  }
  if (rank === 2) {
    return {
      bg: "bg-jv-salmon",
      icon: Medal,
      label: "2.",
      rotate: "",
      scale: "",
    };
  }
  return {
    bg: "bg-jv-mint",
    icon: Award,
    label: "3.",
    rotate: "",
    scale: "",
  };
});

const getOrdinal = (rank) => `${rank}.`;
</script>

<template>
  <article
    :class="[
      'relative flex w-full max-w-[280px] flex-col items-center gap-3 jv-border-rough p-5 shadow-soft-lg sm:gap-4 sm:p-6',
      rankConfig.bg,
      rankConfig.rotate,
      rankConfig.scale,
    ]"
    :aria-label="`${getOrdinal(props.winner.rank)} místo: ${
      props.winner.firstname
    }, skóre ${props.winner.score}`"
    role="article"
  >
    <div
      class="inline-flex items-center gap-2 rounded-full border border-jv-line bg-jv-white px-3 py-1 font-feature text-[14px] font-black text-jv-ink shadow-soft-sm sm:text-[15px]"
    >
      <component
        :is="rankConfig.icon"
        class="size-4 text-jv-ink"
        :stroke-width="2.4"
      />
      {{ rankConfig.label }} místo
    </div>

    <div
      class="grid size-24 place-items-center overflow-hidden rounded-full border border-jv-line bg-jv-white shadow-soft-sm sm:size-28"
    >
      <NuxtImg
        :src="avatar"
        :alt="`Avatar hráče ${props.winner.firstname}`"
        width="96"
        height="96"
        class="size-full object-cover"
      />
    </div>

    <div class="flex w-full flex-col items-center gap-1 text-center">
      <h2
        :id="`winner-name-${props.winner.rank}`"
        class="break-words font-headings text-[22px] uppercase leading-tight text-jv-ink sm:text-[26px]"
      >
        {{ props.winner.firstname }}
      </h2>
      <p
        v-if="props.winner.username"
        class="font-body text-[12px] font-bold text-jv-muted sm:text-[13px]"
      >
        @{{ props.winner.username }}
      </p>
    </div>

    <div
      class="mt-1 inline-flex items-baseline gap-2 rounded-full border border-jv-line bg-jv-white px-4 py-1.5 shadow-soft-sm"
    >
      <span class="font-body text-[11px] font-bold text-jv-muted">SKÓRE</span>
      <span
        class="font-feature text-[22px] font-black tabular-nums text-jv-ink sm:text-[26px]"
      >
        {{ props.winner.score }}
      </span>
    </div>
  </article>
</template>
