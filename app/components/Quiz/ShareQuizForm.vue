<script setup>
import { ref, watch } from "vue";
import { Send, UserCog } from "lucide-vue-next";

const props = defineProps({
  formTitle: {
    type: String,
    required: true,
    default: "",
  },
  id: {
    type: String,
    required: false,
    default: "",
  },
  email: {
    type: String,
    required: false,
    default: "",
  },
  permission: {
    type: String,
    required: false,
    default: "",
  },
});

const emits = defineEmits(["shareQuiz", "updateUserPermission"]);

const email = ref(props.email);
const permission = ref(props.permission);

// Keep local state in sync when editing an existing permission
watch(
  () => props.email,
  (newEmail) => {
    email.value = newEmail;
  }
);

watch(
  () => props.permission,
  (newPermission) => {
    permission.value = newPermission;
  }
);

const handleSubmit = () => {
  if (props.id) {
    emits("updateUserPermission", props.id, email.value, permission.value);
  } else {
    emits("shareQuiz", email.value, permission.value);
  }

  email.value = "";
  permission.value = "";
};
</script>

<template>
  <form
    class="mt-2 jv-border-rough border border-jv-line bg-jv-canvas p-4 shadow-soft-sm sm:p-5"
    @submit.prevent="handleSubmit"
  >
    <h3
      class="mb-4 flex items-center gap-2 font-body text-[15px] font-black uppercase tracking-[0.14em] text-jv-ink"
    >
      <UserCog class="size-4" :stroke-width="2.6" />
      {{ props.formTitle }}
    </h3>

    <!-- Email -->
    <label class="mb-3 grid gap-2">
      <span
        class="text-[12px] font-black uppercase tracking-[0.16em] text-jv-ink"
      >
        E-mail <span class="text-jv-coral">*</span>
      </span>
      <input
        v-model="email"
        type="email"
        name="identifier"
        required
        :disabled="props.id !== ''"
        class="h-12 border border-jv-line bg-jv-white px-3 text-[15px] font-semibold text-jv-ink outline-none transition-shadow focus:shadow-soft-sm disabled:cursor-not-allowed disabled:opacity-60"
      />
    </label>

    <!-- Permission -->
    <label class="mb-4 grid gap-2">
      <span
        class="text-[12px] font-black uppercase tracking-[0.16em] text-jv-ink"
      >
        Oprávnění <span class="text-jv-coral">*</span>
      </span>
      <select
        v-model="permission"
        required
        class="h-12 border border-jv-line bg-jv-white px-3 text-[15px] font-semibold text-jv-ink outline-none transition-shadow focus:shadow-soft-sm"
      >
        <option value="" disabled>Vyber úroveň oprávnění</option>
        <option value="read">Jen čtení</option>
        <option value="write">Může upravovat</option>
        <option value="share">Může spravovat</option>
      </select>
    </label>

    <!-- Submit -->
    <button
      type="submit"
      class="inline-flex h-11 w-full items-center justify-center gap-2 rounded-full border border-jv-line bg-jv-coral px-5 font-body text-[15px] font-black text-white shadow-soft-sm transition-transform"
    >
      <Send class="size-4" :stroke-width="2.4" />
      {{ props.id ? "Upravit přístup" : "Sdílet kvíz" }}
    </button>
  </form>
</template>
