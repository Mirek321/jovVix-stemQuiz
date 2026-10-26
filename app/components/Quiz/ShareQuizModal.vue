<script setup>
import { computed, ref, watch } from "vue";
import { usePush } from "notivue";
import { Lock, UserPlus, Users } from "lucide-vue-next";
import { Modal } from "@/components/ui/modal";
import { useUsersStore } from "@/store/users";
import ShareQuizAuthorizeUser from "./ShareQuizAuthorizeUser.vue";
import ShareQuizForm from "./ShareQuizForm.vue";

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  quizId: {
    type: String,
    required: true,
    default: "",
  },
});

const emits = defineEmits(["update:modelValue"]);

const toast = usePush();
const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const usersStore = useUsersStore();

// edit-form state
const id = ref("");
const email = ref("");
const permission = ref("");
const component = ref("");

// Authorized users for this quiz — reactive to quizId so the modal works
// from both the quiz-detail page and the per-quiz list cards.
const {
  refresh: quizAuthorizedUsersDataRefresh,
  data: quizAuthorizedUsersData,
  pending: quizAuthorizedUsersPending,
  error: quizAuthorizedUsersError,
} = useFetch(() => `${url.apiUrl}/shared_quizzes/${props.quizId}`, {
  method: "GET",
  headers,
  mode: "cors",
  credentials: "include",
  immediate: false,
  watch: [() => props.quizId],
});

const authorizedUsers = computed(
  () => quizAuthorizedUsersData.value?.data || []
);

const isSelf = (user) =>
  !!usersStore.userData?.email &&
  user.shared_to?.toLowerCase() === usersStore.userData.email.toLowerCase();

const isPermissionDenied = computed(() => {
  const code =
    quizAuthorizedUsersError.value?.data?.code ??
    quizAuthorizedUsersError.value?.statusCode;
  return code === 401 || code === 403;
});

const accessErrorMessage = computed(() => {
  if (!quizAuthorizedUsersError.value) return "";
  if (isPermissionDenied.value) {
    return "Nemáš oprávnění sdílet tento kvíz. Pokud chceš zvát další lidi, požádej autora kvízu o oprávnění „Může spravovat“.";
  }
  return "Nepodařilo se načíst, kdo má ke kvízu přístup. Zkus to znovu.";
});

function close() {
  emits("update:modelValue", false);
}

// Reset the form whenever the modal opens, and (re)load the user list.
watch(
  () => props.modelValue,
  (open) => {
    if (open) {
      component.value = "";
      id.value = "";
      email.value = "";
      permission.value = "";
      if (props.quizId) quizAuthorizedUsersDataRefresh();
    }
  }
);

const showEditForm = (idVal, emailVal, permissionVal) => {
  id.value = idVal;
  email.value = emailVal;
  permission.value = permissionVal;
  component.value = "edit-permission";
};

const shareQuiz = async (emailVal, permissionVal) => {
  try {
    await $fetch(`${url.apiUrl}/shared_quizzes/${props.quizId}`, {
      method: "POST",
      headers,
      body: { email: emailVal, permission: permissionVal },
      credentials: "include",
    });
    toast.success("Kvíz byl nasdílen");
    quizAuthorizedUsersDataRefresh();
    component.value = "";
  } catch (error) {
    console.error("Failed to share the quiz.", error);
    toast.error(error.data?.message || "Kvíz se nepodařilo nasdílet.");
  }
};

const updateUserPermission = async (idVal, emailVal, permissionVal) => {
  try {
    await $fetch(
      `${url.apiUrl}/shared_quizzes/${props.quizId}?shared_quiz_id=${idVal}`,
      {
        method: "PUT",
        headers,
        body: { email: emailVal, permission: permissionVal },
        credentials: "include",
      }
    );
    toast.success("Oprávnění bylo upraveno");
    quizAuthorizedUsersDataRefresh();
    component.value = "";
  } catch (error) {
    console.error("Failed to update user permission.", error);
    toast.error(error.data?.message || "Oprávnění se nepodařilo upravit.");
  }
};

const deleteUserPermission = async (idVal) => {
  try {
    await $fetch(
      `${url.apiUrl}/shared_quizzes/${props.quizId}?shared_quiz_id=${idVal}`,
      {
        method: "DELETE",
        headers,
        credentials: "include",
      }
    );
    toast.success("Oprávnění bylo odebráno");
    quizAuthorizedUsersDataRefresh();
  } catch (error) {
    console.error("Failed to delete user permission.", error);
    toast.error(error.data?.message || "Oprávnění se nepodařilo odebrat.");
  }
};
</script>

<template>
  <Modal
    :model-value="props.modelValue"
    title="Sdílet kvíz"
    description="Dej ostatním přístup ke kvízu a spravuj jejich oprávnění."
    size="md"
    @update:model-value="emits('update:modelValue', $event)"
  >
    <div
      class="mb-4 h-[3px] w-full rounded-full border-b border-dashed border-jv-line"
      aria-hidden="true"
    ></div>

    <div v-if="quizAuthorizedUsersPending" class="py-6 text-center">
      <p class="text-[15px] font-semibold text-jv-muted">Načítání…</p>
    </div>

    <div
      v-else-if="quizAuthorizedUsersError"
      class="jv-border-rough flex items-start gap-3 border border-jv-line p-4 text-[15px] font-semibold"
      :class="
        isPermissionDenied
          ? 'bg-jv-canvas text-jv-ink'
          : 'bg-jv-white text-jv-coral'
      "
    >
      <Lock class="mt-[2px] size-5 shrink-0" :stroke-width="2.6" />
      <p>{{ accessErrorMessage }}</p>
    </div>

    <template v-else>
      <!-- People with access -->
      <div class="flex items-center justify-between gap-3">
        <h3
          class="flex items-center gap-2 font-body text-[15px] font-black uppercase tracking-[0.14em] text-jv-ink"
        >
          <Users class="size-4" :stroke-width="2.6" />
          Lidé s přístupem
        </h3>
        <button
          type="button"
          title="Přidat lidi"
          aria-label="Přidat lidi"
          class="inline-flex items-center gap-2 rounded-full border border-jv-line bg-jv-white px-3 py-1.5 text-[13px] font-black text-jv-ink shadow-soft-sm transition-transform"
          @click="component = 'give-permission'"
        >
          <UserPlus class="size-4" :stroke-width="2.6" />
          Přidat
        </button>
      </div>

      <ul class="mt-4 flex max-h-[320px] flex-col gap-3 overflow-y-auto pr-1">
        <li v-for="(user, i) in authorizedUsers" :key="i">
          <ShareQuizAuthorizeUser
            :user="user"
            :is-self="isSelf(user)"
            @show-edit-form="showEditForm"
            @delete-user-permission="deleteUserPermission"
          />
        </li>
        <li
          v-if="authorizedUsers.length === 0"
          class="py-4 text-center text-[14px] font-semibold text-jv-muted"
        >
          Zatím nikdo nemá přístup. Začni tím, že někoho přidáš.
        </li>
      </ul>

      <!-- Add-permission form -->
      <ShareQuizForm
        v-if="component === 'give-permission'"
        form-title="Přidat lidi"
        @share-quiz="shareQuiz"
      />

      <!-- Edit-permission form -->
      <ShareQuizForm
        v-if="component === 'edit-permission'"
        :id="id"
        form-title="Upravit oprávnění"
        :email="email"
        :permission="permission"
        @update-user-permission="updateUserPermission"
      />
    </template>

    <template #footer>
      <button
        type="button"
        class="inline-flex h-10 items-center justify-center rounded-full border border-jv-line bg-jv-white px-5 font-body text-[14px] font-black text-jv-ink shadow-soft-sm transition-transform sm:h-11 sm:px-6 sm:text-[15px]"
        @click="close"
      >
        Hotovo
      </button>
    </template>
  </Modal>
</template>
