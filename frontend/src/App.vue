<script setup lang="ts">
import { computed } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useCardStore } from "@/stores/cardStore";
import { useAuthStore } from "@/stores/authStore";

const router = useRouter();
const route  = useRoute();
const card   = useCardStore();
const auth   = useAuthStore();

const isHome    = computed(() => route.path === "/");
const isLogin   = computed(() => route.path === "/login");
const isChat    = computed(() => route.path.startsWith("/chat/"));
const isProfile = computed(() => route.path === "/profile");
</script>

<template>
  <div
    class="min-h-screen w-full relative overflow-hidden flex flex-col items-center"
    style="background: #0d0d0d;"
  >
    <!-- Ambient blobs -->
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <div
        class="absolute -top-40 -left-40 w-[700px] h-[700px] rounded-full opacity-25 blur-3xl transition-all duration-700"
        :style="`background: radial-gradient(circle, ${card.blobColor1}, transparent)`"
      />
      <div
        class="absolute -bottom-40 -right-40 w-[600px] h-[600px] rounded-full opacity-20 blur-3xl transition-all duration-700"
        :style="`background: radial-gradient(circle, ${card.blobColor2}, transparent)`"
      />
      <div
        class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[500px] h-[500px] rounded-full opacity-15 blur-3xl transition-all duration-700"
        :style="`background: radial-gradient(circle, ${card.blobColor3}, transparent)`"
      />
    </div>

    <!-- Header — скрываем на странице логина и чата -->
    <header
      v-if="!isLogin && !isChat"
      class="relative z-10 w-full max-w-2xl px-4 pt-8 pb-4 flex items-center justify-between"
    >
      <button @click="router.push('/')" class="flex flex-col leading-none">
        <span class="text-white font-bold text-2xl tracking-tight">
          neko<span style="color: #ec4899">swipe</span>
        </span>
        <span class="text-white/25 text-[10px] tracking-[0.3em] uppercase font-medium">find your catgirl</span>
      </button>

      <div class="flex items-center gap-2">
        <!-- username -->
        <span v-if="auth.user" class="text-white/30 text-xs hidden sm:block">
          {{ auth.user.username }}
        </span>

        <!-- Liked -->
        <button
          class="flex items-center gap-1.5 px-4 py-2 rounded-2xl border transition-all duration-200 text-sm font-semibold"
          :class="isHome
            ? 'bg-white/5 border-white/10 text-white/60 hover:bg-white/10 hover:text-white'
            : 'bg-pink-500/20 border-pink-400/30 text-pink-300'"
          @click="router.push(isHome ? '/liked' : '/')"
        >
          <span v-if="isHome" class="flex items-center gap-1.5">💚 Лайки</span>
          <span v-else class="flex items-center gap-1.5">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
            </svg>
            Свайпать
          </span>
        </button>

        <!-- Profile -->
        <button
          class="w-9 h-9 rounded-2xl border transition-all duration-200 flex items-center justify-center"
          :class="isProfile
            ? 'bg-violet-500/20 border-violet-400/30 text-violet-300'
            : 'bg-white/5 border-white/10 text-white/40 hover:text-white/80 hover:bg-white/10'"
          title="Профиль"
          @click="router.push('/profile')"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
        </button>

        <!-- Logout -->
        <button
          class="w-9 h-9 rounded-2xl bg-white/5 border border-white/10 text-white/40 hover:text-white/80 hover:bg-white/10 transition flex items-center justify-center"
          title="Выйти"
          @click="auth.logout().then(() => router.push('/login'))"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a2 2 0 01-2 2H5a2 2 0 01-2-2V7a2 2 0 012-2h6a2 2 0 012 2v1" />
          </svg>
        </button>
      </div>
    </header>

    <!-- Page content -->
    <main
      class="relative z-10 w-full flex justify-center flex-1 pb-8"
      :class="{ 'items-center': isLogin }"
    >
      <RouterView />
    </main>
  </div>
</template>
