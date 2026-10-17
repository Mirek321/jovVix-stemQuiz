<template>
  <div class="relative z-10 mx-auto w-full max-w-[460px]">
    <div class="relative">
      <span
        class="jv-card absolute -top-[8px] left-1/2 z-20 h-3 w-12 -translate-x-1/2 border border-jv-line bg-jv-slate shadow-soft-sm"
        aria-hidden="true"
      ></span>

      <div
        class="jv-card border border-jv-line bg-jv-white px-6 py-7 shadow-soft-lg sm:px-8 sm:py-9"
      >
        <header class="mb-6 flex flex-col items-center gap-1.5">
          <div class="relative inline-block">
            <h1
              class="m-0 font-headings text-[26px] leading-none text-jv-ink sm:text-[32px]"
            >
              Změnit heslo
            </h1>
            <svg
              class="absolute -bottom-2 left-1/2 -translate-x-1/2 text-jv-mint"
              width="160"
              height="14"
              viewBox="0 0 160 14"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
              aria-hidden="true"
            >
              <path
                d="M3 9 Q 30 1, 60 7 T 110 6 T 157 4"
                stroke="currentColor"
                stroke-width="2.5"
                stroke-linecap="round"
                fill="none"
              />
            </svg>
          </div>
          <p
            class="mt-2 text-center font-body text-sm text-jv-ink/70 sm:text-[15px]"
          >
            Zadej své nové heslo.
          </p>
        </header>

        <form
          class="flex flex-col gap-4"
          @submit.prevent="handleChangePassword"
        >
          <div class="flex flex-col gap-1.5">
            <label
              for="newPassword"
              class="px-0.5 font-body text-xs font-bold uppercase tracking-wide text-jv-ink sm:text-[13px]"
            >
              Nové heslo
            </label>
            <div
              class="jv-card flex items-center gap-2.5 border border-jv-line bg-jv-white px-3 py-2.5 shadow-soft-sm transition-all focus-within:translate-x-[1px] focus-within:translate-y-[1px] focus-within:shadow-none"
            >
              <Lock
                class="size-[18px] shrink-0 text-jv-ink/70"
                :stroke-width="2.2"
              />
              <input
                id="newPassword"
                v-model="password.new"
                type="password"
                placeholder="Zadej nové heslo"
                class="min-w-0 flex-1 border-0 bg-transparent font-body text-sm text-jv-ink outline-none placeholder:text-jv-ink/40 sm:text-base"
                required
              />
            </div>
          </div>

          <div class="flex flex-col gap-1.5">
            <label
              for="confirmPassword"
              class="px-0.5 font-body text-xs font-bold uppercase tracking-wide text-jv-ink sm:text-[13px]"
            >
              Potvrzení hesla
            </label>
            <div
              class="jv-card flex items-center gap-2.5 border border-jv-line bg-jv-white px-3 py-2.5 shadow-soft-sm transition-all focus-within:translate-x-[1px] focus-within:translate-y-[1px] focus-within:shadow-none"
            >
              <Lock
                class="size-[18px] shrink-0 text-jv-ink/70"
                :stroke-width="2.2"
              />
              <input
                id="confirmPassword"
                v-model="password.confirm"
                type="password"
                placeholder="Zopakuj nové heslo"
                class="min-w-0 flex-1 border-0 bg-transparent font-body text-sm text-jv-ink outline-none placeholder:text-jv-ink/40 sm:text-base"
                required
              />
            </div>
          </div>

          <div
            v-if="
              passwordRequestError ||
              (passwordSubmitted && passwordErrors.length)
            "
            class="flex flex-col gap-1 rounded-[8px] border border-jv-coral bg-jv-coral/10 px-3 py-2"
          >
            <p
              v-if="passwordRequestError"
              class="flex items-start gap-1.5 font-body text-xs font-bold text-jv-coral"
            >
              <AlertCircle
                class="mt-0.5 size-3.5 shrink-0"
                :stroke-width="2.4"
              />
              {{ passwordRequestError }}
            </p>
            <p
              v-for="(err, i) in passwordErrors"
              v-show="passwordSubmitted"
              :key="i"
              class="flex items-start gap-1.5 font-body text-xs font-bold text-jv-coral"
            >
              <AlertCircle
                class="mt-0.5 size-3.5 shrink-0"
                :stroke-width="2.4"
              />
              {{ err }}
            </p>
          </div>

          <div class="mt-2 flex flex-col gap-3">
            <button
              type="submit"
              class="jv-card inline-flex h-12 items-center justify-center gap-2 border border-jv-line bg-jv-coral font-headings text-base text-white shadow-soft-sm transition-transform sm:text-lg"
            >
              Změnit heslo
            </button>
            <button
              type="button"
              class="jv-card inline-flex h-11 items-center justify-center border border-jv-line bg-jv-white font-headings text-base text-jv-ink shadow-soft-sm transition-transform"
              @click="handleCancel"
            >
              Zrušit
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { usePush } from "notivue";
import { toRef } from "vue";
import { Lock, AlertCircle } from "lucide-vue-next";
import { useUserPasswordRules } from "@/composables/user_password_rules";
import { kratosMessageText } from "~/composables/kratosMessages";

definePageMeta({
  layout: "auth",
});

useSeoMeta({
  title: "Změnit heslo – STEM_HUB",
  description:
    "Změň si heslo k účtu v systému STEM_HUB, ať jsou tvoje kvízy a data v bezpečí.",
  robots: "noindex, nofollow",
});

const toast = usePush();
const url = useRuntimeConfig().public;

const password = reactive({
  new: "",
  confirm: "",
});

const newPasswordRef = toRef(password, "new");

const { passwordErrors } = useUserPasswordRules(newPasswordRef);

const passwordRequestError = ref("");
const passwordSubmitted = ref(false);
const flow = ref("");
const csrfToken = ref("");

const fetchFlowIdAndCsrfToken = async () => {
  try {
    const response = await fetch(
      `${url.kratosUrl}/self-service/settings/browser?aal=&refresh=&return_to=`,
      {
        method: "GET",
        headers: {
          Accept: "application/json",
        },
        credentials: "include",
      }
    );

    if (!response.ok) {
      throw new Error(`Error: ${response.statusText}`);
    }

    const data = await response.json();
    flow.value = data.id;
    csrfToken.value = data.ui.nodes.find(
      (node) => node.attributes.name === "csrf_token"
    ).attributes.value;
  } catch (error) {
    console.error("Error fetching flow ID and CSRF token:", error);
  }
};

const handleChangePassword = async () => {
  passwordSubmitted.value = true;

  try {
    if (passwordErrors.value.length > 0) {
      return;
    }

    if (password.new !== password.confirm) {
      passwordRequestError.value = "Hesla se neshodují.";
      return;
    }

    await fetchFlowIdAndCsrfToken();

    const response = await fetch(
      `${url.kratosUrl}/self-service/settings?flow=${flow.value}`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        credentials: "include",
        body: JSON.stringify({
          password: password.new,
          csrf_token: csrfToken.value,
          method: "password",
        }),
      }
    );

    if (!response.ok) {
      const errorData = await response.json();
      if (errorData?.error?.id === "session_refresh_required") {
        window.location.href = errorData.redirect_browser_to;
        return;
      }

      const nodeError = kratosMessageText(
        errorData?.ui?.nodes
          ?.flatMap((n) => n.messages || [])
          ?.find((m) => m.type === "error")
      );
      const uiError = kratosMessageText(
        errorData?.ui?.messages?.find((m) => m.type === "error")
      );

      throw new Error(
        nodeError ||
          uiError ||
          errorData?.error?.message ||
          errorData?.message ||
          "Heslo se nepodařilo změnit."
      );
    }

    passwordRequestError.value = "";

    password.new = "";
    password.confirm = "";

    toast.success("Heslo bylo úspěšně změněno.");

    navigateTo("/admin");
  } catch (error) {
    passwordRequestError.value = error.message;
  }
};

const handleCancel = () => {
  navigateTo("/admin");
};
</script>
