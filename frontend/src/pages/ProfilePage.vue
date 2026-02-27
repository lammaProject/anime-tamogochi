<script setup lang="ts">
import {computed, onMounted, ref} from "vue";
import {useRouter} from "vue-router";
import type {ChatEntry} from "@/api/backend";
import {getChats} from "@/api/backend";
import {useAuthStore} from "@/stores/authStore";

const router = useRouter();
const auth = useAuthStore();
const chats = ref<ChatEntry[]>([]);
const loading = ref(true);

onMounted(async () => {
  try {
    const { data } = await getChats();
    chats.value = data;
  } catch {
    chats.value = [];
  } finally {
    loading.value = false;
  }
});

function goToChat(entry: ChatEntry) {
  const img =
    entry.data?.image.compressed?.url || entry.data?.image.original?.url || "";
  const name = entry.data?.anime.character || "";
  router.push({ path: `/chat/${entry.catgirl_id}`, query: { name, img, routePath: '/profile'  } });
}

function formatTime(iso: string) {
  const d = new Date(iso);
  const now = new Date();
  const diffMs = now.getTime() - d.getTime();
  const diffMin = Math.floor(diffMs / 60000);
  const diffH = Math.floor(diffMin / 60);
  const diffD = Math.floor(diffH / 24);

  if (diffMin < 1) return "только что";
  if (diffMin < 60) return `${diffMin} мин. назад`;
  if (diffH < 24) return `${diffH} ч. назад`;
  if (diffD < 7) return `${diffD} д. назад`;
  return d.toLocaleDateString("ru-RU", { day: "numeric", month: "short" });
}

function pluralMessages(n: number) {
  if (n % 100 >= 11 && n % 100 <= 14) return `${n} сообщений`;
  const r = n % 10;
  if (r === 1) return `${n} сообщение`;
  if (r >= 2 && r <= 4) return `${n} сообщения`;
  return `${n} сообщений`;
}

const totalMessages = computed(() =>
  chats.value.reduce((s, c) => s + c.message_count, 0)
);
</script>

<template>
  <div class="w-full max-w-2xl px-4 flex flex-col gap-6">

    <!-- Profile header block -->
    <div
      class="rounded-3xl overflow-hidden relative"
      style="background: linear-gradient(135deg, rgba(124,58,237,0.18), rgba(236,72,153,0.18)); border: 1px solid rgba(255,255,255,0.07);"
    >
      <div class="px-6 py-5 flex items-center gap-4">
        <!-- Avatar placeholder -->
        <div
          class="w-16 h-16 rounded-2xl flex-shrink-0 flex items-center justify-center text-3xl font-bold text-white/80"
          style="background: linear-gradient(135deg, #7c3aed, #ec4899);"
        >
          {{ auth.user?.username?.[0]?.toUpperCase() ?? "?" }}
        </div>
        <div class="flex-1 min-w-0">
          <h2 class="text-white font-bold text-xl leading-tight">
            {{ auth.user?.username ?? "Аноним" }}
          </h2>
          <p class="text-white/40 text-sm mt-0.5">Профиль игрока</p>
        </div>
      </div>

      <!-- Stats row -->
      <div
        class="grid grid-cols-2 divide-x divide-white/5 border-t"
        style="border-color: rgba(255,255,255,0.06);"
      >
        <div class="px-6 py-3 flex flex-col items-center gap-0.5">
          <span class="text-white font-bold text-lg">{{ chats.length }}</span>
          <span class="text-white/35 text-xs">переписок</span>
        </div>
        <div class="px-6 py-3 flex flex-col items-center gap-0.5">
          <span class="text-white font-bold text-lg">{{ totalMessages }}</span>
          <span class="text-white/35 text-xs">сообщений</span>
        </div>
      </div>
    </div>

    <!-- Section title -->
    <div>
      <h3 class="text-white font-bold text-lg">Переписки</h3>
      <p class="text-white/35 text-xs mt-0.5">Все катгёрлс с которыми ты общался</p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-16">
      <div class="w-10 h-10 rounded-full border-4 border-white/20 border-t-white animate-spin" />
    </div>

    <!-- Empty -->
    <div
      v-else-if="!chats.length"
      class="flex flex-col items-center gap-4 py-16 text-center"
    >
      <div class="text-5xl">💬</div>
      <p class="text-white/50 text-sm">Ты ещё ни с кем не переписывался</p>
      <button
        class="px-6 py-2.5 rounded-2xl bg-white/10 text-white text-sm font-semibold hover:bg-white/20 transition border border-white/10"
        @click="router.push('/liked')"
      >
        Перейти к лайкам
      </button>
    </div>

    <!-- Chat list -->
    <div v-else class="flex flex-col gap-3">
      <div
        v-for="entry in chats"
        :key="entry.catgirl_id"
        class="group flex items-center gap-4 rounded-2xl px-4 py-3 cursor-pointer transition-all duration-200 hover:scale-[1.01]"
        style="background: rgba(255,255,255,0.04); border: 1px solid rgba(255,255,255,0.07);"
        @click="goToChat(entry)"
      >
        <!-- Avatar -->
        <div
          class="relative w-14 h-14 rounded-2xl flex-shrink-0 overflow-hidden"
          :style="{ background: entry.data?.colors.main ?? '#2d2d2d' }"
        >
          <img
            v-if="entry.data?.image.compressed?.url || entry.data?.image.original?.url"
            :src="entry.data?.image.compressed?.url || entry.data?.image.original?.url"
            :alt="entry.data?.anime.character ?? 'catgirl'"
            class="w-full h-full object-cover"
            loading="lazy"
          />
          <div
            v-else
            class="w-full h-full flex items-center justify-center text-2xl"
          >
            🐱
          </div>
        </div>

        <!-- Info -->
        <div class="flex-1 min-w-0">
          <!-- Name + time -->
          <div class="flex items-center justify-between gap-2">
            <span class="text-white font-semibold text-sm truncate">
              {{ entry.data?.anime.character || "Неизвестно" }}
            </span>
            <span class="text-white/30 text-xs flex-shrink-0">
              {{ formatTime(entry.last_message_at) }}
            </span>
          </div>

          <!-- Anime + tags -->
          <p v-if="entry.data?.anime.title" class="text-white/40 text-xs truncate mt-0.5">
            {{ entry.data.anime.title }}
          </p>

          <!-- Tags row -->
          <div v-if="entry.data?.tags?.length" class="flex flex-wrap gap-1 mt-1.5">
            <span
              v-for="tag in entry.data.tags.slice(0, 3)"
              :key="tag"
              class="text-[10px] font-medium px-1.5 py-0.5 rounded-full text-white/50"
              style="background: rgba(255,255,255,0.06);"
            >
              #{{ tag }}
            </span>
          </div>

          <!-- Last message preview + count -->
          <div class="flex items-center justify-between mt-1.5 gap-2">
            <p class="text-white/30 text-xs truncate italic">
              {{ entry.last_message || "..." }}
            </p>
            <span
              class="text-[10px] font-semibold px-2 py-0.5 rounded-full flex-shrink-0"
              style="background: rgba(236,72,153,0.18); color: rgba(244,114,182,0.85);"
            >
              {{ pluralMessages(entry.message_count) }}
            </span>
          </div>
        </div>

        <!-- Arrow -->
        <div class="flex-shrink-0 text-white/20 group-hover:text-white/60 transition">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
          </svg>
        </div>
      </div>
    </div>

  </div>
</template>
