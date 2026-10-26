<template>
  <div class="lg:flex lg:h-screen lg:overflow-hidden">
    <Sidebar v-if="!hideSidebar" />
    <div
      class="flex min-h-screen flex-col lg:min-h-0 lg:flex-1 lg:overflow-y-auto"
    >
      <div class="flex-1">
        <slot></slot>
      </div>
      <SiteFooter v-if="!hideFooter" />
    </div>
  </div>
</template>

<script setup>
import Sidebar from "@/components/common/Sidebar.vue";
import SiteFooter from "@/components/common/SiteFooter.vue";

const route = useRoute();
const hideSidebar = computed(() => route.meta.hideSidebar === true);
// Game screens (question, timer) stay full-height without the footer.
const hideFooter = computed(() => route.meta.hideFooter === true);
</script>
