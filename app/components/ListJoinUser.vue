<script setup>
import { Users } from "lucide-vue-next";
import { storeToRefs } from "pinia";
import { useListUserstore } from "~/store/userlist";
import { getAvatarUrlByName } from "~~/composables/avatar";

const listUserStore = useListUserstore();
const { listUsers } = storeToRefs(listUserStore);
</script>

<template>
  <div class="mx-auto w-full max-w-[800px]">
    <div class="mb-4 flex justify-center">
      <div
        class="jv-border-rough inline-flex items-center gap-3 border border-jv-line bg-jv-yellow px-5 py-2.5 shadow-soft-sm"
      >
        <Users class="size-5 text-jv-ink" :stroke-width="2.4" />
        <h5
          v-if="listUsers.length === 0"
          class="m-0 font-headings text-base text-jv-ink sm:text-lg"
        >
          Waiting for Participants...
        </h5>
        <h5 v-else class="m-0 font-headings text-base text-jv-ink sm:text-lg">
          {{ listUsers.length }} Participants
        </h5>
      </div>
    </div>

    <div v-if="listUsers.length" class="flex flex-wrap justify-center gap-3">
      <div
        v-for="user in listUsers"
        :key="user.UserId"
        class="jv-card flex items-center gap-3 border border-jv-line bg-jv-white pr-4 shadow-soft-sm"
      >
        <NuxtImg
          :src="getAvatarUrlByName(user?.Avatar)"
          alt=""
          width="48"
          height="48"
          loading="lazy"
          class="size-12 rounded-full border border-jv-line object-cover"
        />
        <span class="pr-2 font-body text-sm font-bold text-jv-ink sm:text-base">
          {{ user.UserName }}
        </span>
      </div>
    </div>
  </div>
</template>
