<script setup>
import { ChevronRight, Footprints, Radio } from "lucide-vue-next";
import { Modal } from "@/components/ui/modal";
import { QUIZ_MODES } from "~~/composables/self_paced";

defineProps({
  modelValue: {
    type: Boolean,
    required: true,
  },
});
const emits = defineEmits(["update:modelValue", "select"]);

const MODES = [
  {
    value: QUIZ_MODES.live,
    title: "Živě s moderátorem",
    hint: "Otázky posouváš ty a všichni odpovídají najednou.",
    icon: Radio,
    badge: "bg-jv-yellow",
  },
  {
    value: QUIZ_MODES.selfPaced,
    title: "Vlastním tempem",
    hint: "Hráči hrají hned po připojení, bez čekání na moderátora.",
    icon: Footprints,
    badge: "bg-jv-mint",
  },
];

const choose = (mode) => {
  emits("update:modelValue", false);
  emits("select", mode);
};
</script>

<template>
  <Modal
    :model-value="modelValue"
    title="Jak chceš kvíz spustit?"
    size="sm"
    @update:model-value="(open) => emits('update:modelValue', open)"
  >
    <div class="flex flex-col gap-3">
      <button
        v-for="mode in MODES"
        :key="mode.value"
        type="button"
        class="group flex w-full items-center gap-3 rounded-2xl border-2 border-jv-line bg-jv-white px-3 py-3 text-left shadow-soft-sm transition-all duration-150 hover:-translate-y-0.5 hover:border-jv-primary hover:shadow-soft sm:gap-4 sm:px-4 sm:py-4"
        @click="choose(mode.value)"
      >
        <span
          :class="[
            'grid size-11 shrink-0 place-items-center rounded-xl text-jv-ink sm:size-12',
            mode.badge,
          ]"
        >
          <component :is="mode.icon" class="size-5" :stroke-width="2.4" />
        </span>
        <span class="min-w-0 flex-1">
          <span
            class="block font-headings text-[17px] leading-tight text-jv-ink sm:text-[19px]"
          >
            {{ mode.title }}
          </span>
          <span
            class="mt-1 block font-body text-[13px] font-bold leading-snug text-jv-muted sm:text-[14px]"
          >
            {{ mode.hint }}
          </span>
        </span>
        <ChevronRight
          class="size-5 shrink-0 text-jv-muted transition-transform group-hover:translate-x-0.5"
          :stroke-width="2.4"
        />
      </button>
    </div>
  </Modal>
</template>
