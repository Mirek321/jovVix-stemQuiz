<script setup>
import { Trash2 } from "lucide-vue-next";
import { Modal } from "@/components/ui/modal";

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  title: {
    type: String,
    default: "Smazat položku",
  },
  message: {
    type: String,
    default: "Smazanou položku už nepůjde obnovit. Chceš pokračovat?",
  },
  confirmLabel: {
    type: String,
    default: "Smazat",
  },
  cancelLabel: {
    type: String,
    default: "Zrušit",
  },
});

const emit = defineEmits(["update:modelValue", "confirmDelete"]);

function close() {
  emit("update:modelValue", false);
}

function confirm() {
  emit("confirmDelete");
  close();
}
</script>

<template>
  <Modal
    :model-value="props.modelValue"
    :title="props.title"
    :description="props.message"
    size="sm"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <template #footer>
      <button
        type="button"
        class="inline-flex h-10 items-center justify-center rounded-full border border-jv-line bg-jv-white px-5 font-body text-[14px] font-black text-jv-ink shadow-soft-sm transition-transform sm:h-11 sm:px-6 sm:text-[15px]"
        @click="close"
      >
        {{ props.cancelLabel }}
      </button>
      <button
        type="button"
        class="inline-flex h-10 items-center justify-center gap-2 rounded-full border border-jv-line bg-jv-coral px-5 font-body text-[14px] font-black text-white shadow-soft-sm transition-transform sm:h-11 sm:px-6 sm:text-[15px]"
        @click="confirm"
      >
        <Trash2 class="size-4" :stroke-width="2.4" />
        {{ props.confirmLabel }}
      </button>
    </template>
  </Modal>
</template>
