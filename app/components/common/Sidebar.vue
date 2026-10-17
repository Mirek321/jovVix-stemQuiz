<template>
  <!-- Desktop sidebar (lg and up) -->
  <aside
    class="hidden lg:flex bg-jv-ivory min-w-56 h-screen flex-col justify-between border-r-[4px] border-jv-line px-8 py-5 overflow-y-hidden"
  >
    <div class="flex flex-col gap-4">
      <NuxtLink
        to="/"
        class="flex items-center gap-2 text-[20px] sm:text-[22px] md:text-[24px] font-black tracking-[-0.5px] text-jv-ink no-underline h-auto"
      >
        <NuxtImg
          src="/images/stem-logo.png"
          alt=""
          class="w-11 shrink-0 lg:w-12"
          width="80"
          height="80"
          sizes="56px lg:80px"
          preload
        />
        <span
          class="font-headings text-[18px] font-extrabold leading-none tracking-[0.02em] text-jv-ink whitespace-nowrap lg:text-[20px]"
          >STEM_HUB</span
        >
      </NuxtLink>

      <template v-if="mounted">
        <NavigationLink
          url="/admin/quiz/list-quiz?create=1"
          url-name="Vytvořit kvíz"
        >
          <Plus class="size-[18px]" />
        </NavigationLink>
      </template>
      <template v-else>
        <Skeleton class="h-12 w-full rounded-[12px] bg-jv-ink/10" />
      </template>

      <nav class="mt-5 flex flex-col gap-2">
        <template v-if="mounted">
          <NavigationLink
            v-for="item in navItems"
            :key="item.url"
            :url="item.url"
            :url-name="item.label"
            :class="navItemClass(item)"
            class="relative flex gap-3 justify-start px-4 py-2.5 text-jv-ink/65 no-underline transition-transform group hover:text-jv-ink"
          >
            <component
              :is="item.icon"
              class="size-5 shrink-0 group-hover:text-jv-ink"
              :class="item.active ? 'text-jv-coral' : 'text-jv-ink/45'"
              :stroke-width="2.3"
            />
            <span
              v-if="item.needsAttention"
              class="absolute right-3 top-1/2 size-2 -translate-y-1/2 rounded-full bg-jv-coral"
              aria-label="E-mail není ověřený"
            ></span>
          </NavigationLink>
        </template>
        <template v-else>
          <div
            v-for="i in 4"
            :key="i"
            class="flex items-center gap-3 px-4 py-2.5"
          >
            <Skeleton class="size-5 rounded-full bg-jv-ink/10" />
            <Skeleton class="h-4 w-24 bg-jv-ink/10" />
          </div>
        </template>
      </nav>
    </div>

    <div class="flex flex-col gap-4">
      <div v-if="!mounted" class="flex items-center gap-3 px-2 py-1">
        <Skeleton class="size-10 rounded-full bg-jv-ink/10" />
        <Skeleton class="h-4 flex-1 bg-jv-ink/10" />
        <Skeleton class="size-5 rounded-full bg-jv-ink/10" />
      </div>
      <div
        v-else-if="showAdminNav"
        class="flex items-center gap-3 jv-border-uneven bg-jv-white px-3 py-2 shadow-soft-sm"
      >
        <NuxtImg
          :src="userAvatar"
          :alt="userName"
          width="40"
          height="40"
          loading="lazy"
          class="size-10 shrink-0 rounded-full border border-jv-line bg-jv-slate object-cover"
        />
        <span
          class="min-w-0 flex-1 truncate font-headings text-[15px] text-jv-ink"
        >
          {{ userName }}
        </span>
        <Popover v-model:open="desktopMenuOpen">
          <PopoverTrigger as-child>
            <NavigationLink
              type="button"
              aria-label="Nabídka profilu"
              class="rounded-full bg-white !p-2 border-0 shadow-none hover:bg-jv-ink/10 text-sm sm:text-base"
            >
              <MoreVertical class="size-5" :stroke-width="2.4" />
            </NavigationLink>
          </PopoverTrigger>
          <PopoverContent
            align="end"
            side="top"
            :side-offset="8"
            class="w-44 bg-jv-white p-1.5 shadow-soft"
          >
            <NavigationLink
              type="button"
              aria-label="Odhlásit se"
              class="px-3 py-2 justify-start w-full bg-white !p-2 border-0 shadow-none text-sm hover:bg-jv-coral hover:border hover:jv-border-even hover:rounded-md"
              @click="handleDesktopLogout"
            >
              <LogOut class="size-4" :stroke-width="2.4" />
              <span class="text-sm">Odhlásit se</span>
            </NavigationLink>
          </PopoverContent>
        </Popover>
      </div>
      <template v-else>
        <NavigationLink url="/account/login" url-name="Přihlásit se" />
        <NavigationLink url="/account/register" url-name="Zaregistrovat se" />
      </template>
    </div>
  </aside>

  <!-- Mobile/Tablet header (md and below) -->
  <header
    class="lg:hidden sticky top-0 z-40 border-b border-jv-line bg-jv-ivory/95 px-4 sm:px-6 py-3 backdrop-blur"
  >
    <div class="flex items-center justify-between gap-3">
      <NuxtLink
        to="/"
        class="flex items-center gap-2 text-[20px] sm:text-[22px] md:text-[24px] font-black tracking-[-0.5px] text-jv-ink no-underline h-auto"
      >
        <NuxtImg
          src="/images/stem-logo.png"
          alt=""
          class="w-11 shrink-0 lg:w-12"
          width="56"
          height="56"
          sizes="48px lg:56px"
        />
        <span
          class="font-headings text-[18px] font-extrabold leading-none tracking-[0.02em] text-jv-ink whitespace-nowrap lg:text-[20px]"
          >STEM_HUB</span
        >
      </NuxtLink>

      <button
        type="button"
        class="grid size-10 place-items-center rounded-[10px] border border-jv-line bg-jv-yellow text-jv-ink shadow-soft-sm transition-transform"
        :aria-expanded="open"
        :aria-label="open ? 'Zavřít nabídku' : 'Otevřít nabídku'"
        aria-controls="mobile-nav"
        @click="open = !open"
      >
        <X v-if="open" class="size-5" :stroke-width="2.4" />
        <Menu v-else class="size-5" :stroke-width="2.4" />
      </button>
    </div>

    <nav
      v-show="open"
      id="mobile-nav"
      class="mt-4 grid grid-cols-1 min-[420px]:grid-cols-2 sm:grid-cols-3 gap-3"
    >
      <template v-if="mounted">
        <NuxtLink
          v-for="item in mobileNavItems"
          :key="item.url"
          :to="item.url"
          :class="[
            'inline-flex items-center justify-center gap-1.5 jv-border-uneven px-4 py-2 text-sm font-headings text-jv-ink no-underline shadow-soft-sm transition-transform ',
            item.active ? 'bg-jv-yellow/70' : 'bg-jv-white',
            item.highlight ? 'bg-jv-coral !text-white' : '',
          ]"
          @click="open = false"
        >
          <component
            :is="item.icon"
            v-if="item.icon"
            class="size-4"
            :stroke-width="2.4"
          />
          <span>{{ item.label }}</span>
          <span
            v-if="item.needsAttention"
            class="size-2 shrink-0 rounded-full bg-jv-coral"
            aria-label="E-mail není ověřený"
          ></span>
        </NuxtLink>
        <div
          v-if="showAdminNav"
          class="col-span-full mt-1 flex items-center gap-3 jv-border-uneven bg-jv-white px-3 py-2 shadow-soft-sm"
        >
          <NuxtImg
            :src="userAvatar"
            :alt="userName"
            width="36"
            height="36"
            loading="lazy"
            class="size-9 shrink-0 rounded-full border border-jv-line bg-jv-slate object-cover"
          />
          <span
            class="min-w-0 flex-1 truncate font-headings text-[15px] text-jv-ink"
          >
            {{ userName }}
          </span>
          <Popover v-model:open="mobileMenuOpen">
            <PopoverTrigger as-child>
              <NavigationLink
                type="button"
                aria-label="Nabídka profilu"
                class="rounded-full bg-white !p-2 border-0 shadow-none hover:bg-jv-ink/10 text-sm sm:text-base"
              >
                <MoreVertical class="size-5" :stroke-width="2.4" />
              </NavigationLink>
            </PopoverTrigger>
            <PopoverContent
              align="end"
              :side-offset="8"
              class="w-44 bg-jv-white p-1.5 shadow-soft"
            >
              <NavigationLink
                type="button"
                aria-label="Odhlásit se"
                class="px-3 py-2 justify-start w-full bg-white !p-2 border-0 shadow-none text-sm hover:bg-jv-coral hover:border hover:jv-border-even hover:rounded-md"
                @click="handleMobileLogout"
              >
                <LogOut class="size-4" :stroke-width="2.4" />
                <span class="text-sm">Odhlásit se</span>
              </NavigationLink>
            </PopoverContent>
          </Popover>
        </div>
      </template>
      <template v-else>
        <div
          v-for="i in 4"
          :key="i"
          class="inline-flex items-center justify-center gap-1.5 jv-border-uneven bg-jv-white px-4 py-2 shadow-soft-sm"
        >
          <Skeleton class="size-4 rounded-full bg-jv-ink/10" />
          <Skeleton class="h-3 w-16 bg-jv-ink/10" />
        </div>
        <div
          class="col-span-full mt-1 flex items-center gap-3 jv-border-uneven bg-jv-white px-3 py-2 shadow-soft-sm"
        >
          <Skeleton class="size-9 rounded-full bg-jv-ink/10" />
          <Skeleton class="h-4 flex-1 bg-jv-ink/10" />
          <Skeleton class="size-5 rounded-full bg-jv-ink/10" />
        </div>
      </template>
    </nav>
  </header>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import {
  BarChart3,
  HelpCircle,
  Home,
  LogOut,
  Menu,
  MoreVertical,
  Plus,
  Tag,
  UserRound,
  X,
} from "lucide-vue-next";
import NavigationLink from "@/components/common/NavigationLink.vue";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Skeleton } from "@/components/ui/skeleton";
import { handleLogout, setUserDataStore } from "@/composables/auth";
import {
  refreshEmailVerified,
  useEmailVerified,
} from "@/composables/email_verification";
import { getAvatarUrlByName } from "@/composables/avatar";
import { useUsersStore } from "~~/store/users";

const route = useRoute();
const open = ref(false);
const desktopMenuOpen = ref(false);
const mobileMenuOpen = ref(false);
const mounted = ref(false);
const userDataStore = useUsersStore();
const emailVerified = useEmailVerified();

const isQuizListPage = computed(() =>
  route.path.startsWith("/admin/quiz/list-quiz")
);

const isAdminPage = computed(() => route.path.startsWith("/admin"));

const isLoggedInAdmin = computed(
  () => userDataStore.getUserData()?.role === "admin-user"
);

const showAdminNav = computed(
  () => isLoggedInAdmin.value || isQuizListPage.value || isAdminPage.value
);

const currentUser = computed(() => userDataStore.getUserData());

// Categories manage public-catalog grouping, so the nav item only shows for
// the configured public-quiz admins (same gate as the API enforces).
const canManageCategories = computed(
  () => !!currentUser.value?.canCreatePublicQuiz
);

const userName = computed(
  () => currentUser.value?.firstname || currentUser.value?.username || "Profil"
);

const userAvatar = computed(() =>
  getAvatarUrlByName(currentUser.value?.avatar)
);

const isActiveRoute = (url) => {
  if (url === "/") return route.path === "/";
  if (url === "/admin/quiz/list-quiz")
    return route.path.startsWith("/admin/quiz");
  if (url === "/admin/reports") return route.path.startsWith("/admin/reports");
  if (url === "/admin/categories")
    return route.path.startsWith("/admin/categories");
  if (url === "/admin") return route.path === "/admin";

  return route.path === url;
};

const navItems = computed(() => {
  if (showAdminNav.value) {
    return [
      { label: "Domů", url: "/", icon: Home, active: isActiveRoute("/") },
      {
        label: "Kvízy",
        url: "/admin/quiz/list-quiz",
        icon: HelpCircle,
        active: isActiveRoute("/admin/quiz/list-quiz"),
      },
      {
        label: "Reporty",
        url: "/admin/reports",
        icon: BarChart3,
        active: isActiveRoute("/admin/reports"),
      },
      ...(canManageCategories.value
        ? [
            {
              label: "Kategorie",
              url: "/admin/categories",
              icon: Tag,
              active: isActiveRoute("/admin/categories"),
            },
          ]
        : []),
      {
        label: "Profil",
        url: "/admin",
        icon: UserRound,
        active: isActiveRoute("/admin"),
        needsAttention: emailVerified.value === false,
      },
    ];
  }

  return [{ label: "Domů", url: "/", icon: Home, active: true }];
});

const mobileNavItems = computed(() => {
  if (showAdminNav.value) {
    return [
      { label: "Domů", url: "/", icon: Home, active: isActiveRoute("/") },
      {
        label: "Kvízy",
        url: "/admin/quiz/list-quiz",
        icon: HelpCircle,
        active: isActiveRoute("/admin/quiz/list-quiz"),
      },
      {
        label: "Reporty",
        url: "/admin/reports",
        icon: BarChart3,
        active: isActiveRoute("/admin/reports"),
      },
      ...(canManageCategories.value
        ? [
            {
              label: "Kategorie",
              url: "/admin/categories",
              icon: Tag,
              active: isActiveRoute("/admin/categories"),
            },
          ]
        : []),
      {
        label: "Profil",
        url: "/admin",
        icon: UserRound,
        active: isActiveRoute("/admin"),
        needsAttention: emailVerified.value === false,
      },
      {
        label: "Vytvořit kvíz",
        url: "/admin/quiz/list-quiz?create=1",
        icon: Plus,
      },
    ];
  }

  return [
    { label: "Domů", url: "/", icon: Home, active: true },
    { label: "Zadat kód", url: "/join" },
    {
      label: "Vytvořit kvíz",
      url: "/admin/quiz/list-quiz?create=1",
      icon: Plus,
    },
    { label: "Přihlásit se", url: "/account/login" },
    { label: "Zaregistrovat se", url: "/account/register" },
  ];
});

const navItemClass = (item) =>
  item.active
    ? "rounded-full border border-dashed border-jv-line bg-jv-yellow/20 text-jv-ink shadow-none "
    : "border-0 bg-transparent shadow-none ";

const handleDesktopLogout = async () => {
  desktopMenuOpen.value = false;
  await handleLogout();
};

const handleMobileLogout = async () => {
  mobileMenuOpen.value = false;
  open.value = false;
  await handleLogout();
};

let breakpointMql = null;
const closeAllMenus = () => {
  desktopMenuOpen.value = false;
  mobileMenuOpen.value = false;
  open.value = false;
};

onMounted(async () => {
  await setUserDataStore();
  mounted.value = true;
  refreshEmailVerified();
  if (typeof window !== "undefined") {
    breakpointMql = window.matchMedia("(min-width: 1024px)");
    breakpointMql.addEventListener("change", closeAllMenus);
  }
});

onBeforeUnmount(() => {
  if (breakpointMql) {
    breakpointMql.removeEventListener("change", closeAllMenus);
  }
});
</script>
