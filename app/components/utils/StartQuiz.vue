<script setup>
import { usePush } from "notivue";
import { useListUserstore } from "~/store/userlist";
import { useSessionStore } from "~~/store/session";
import NavigationLink from "@/components/common/NavigationLink.vue";
import StartModeModal from "@/components/Quiz/StartModeModal.vue";
const sessionStore = useSessionStore();
const { setSession } = sessionStore;
const listUserStore = useListUserstore();
const { removeAllUsers } = listUserStore;
let urls = useRuntimeConfig().public;
const router = useRouter();
const toast = usePush();
const requestPending = ref(false);
const activeQuizId = ref(false);
const startModeOpen = ref(false);
const props = defineProps({
  quizId: {
    default: () => {
      return "";
    },
    type: String,
    required: true,
  },
});

const handleStartDemo = async (mode) => {
  try {
    requestPending.value = true;
    await $fetch(
      `${urls.apiUrl}/quizzes/${props.quizId}/demo_session?mode=${mode}`,
      {
        method: "POST",
        credentials: "include",
        headers: {
          Accept: "application/json",
        },
        onResponse({ response }) {
          if (response.status != 202) {
            requestPending.value = false;
            toast.error("Kvíz se nepodařilo spustit");
            return;
          }
          if (response.status == 202) {
            activeQuizId.value = response._data.data;
            requestPending.value = false;
          }
        },
      }
    );
  } catch (error) {
    toast.error(error.message);
    requestPending.value = false;
    return;
  }

  removeAllUsers();
  setSocketObject(null);
  setSession(activeQuizId.value);
  router.push(`/admin/arrange/${activeQuizId.value}`);
};
</script>
<template>
  <NavigationLink
    :url-name="requestPending ? 'Spouštím…' : 'Spustit kvíz'"
    class="h-8 rounded-full shadow-none"
    :disabled="requestPending"
    @click="startModeOpen = true"
  />
  <StartModeModal v-model="startModeOpen" @select="handleStartDemo" />
</template>
