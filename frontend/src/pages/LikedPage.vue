<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { getLiked, removeLiked } from "@/api/backend";
import type { NekosImageData } from "@/api/type";

const router = useRouter();
const items = ref<NekosImageData[]>([]);
const selected = ref<NekosImageData | null>(null);
const loading = ref(true);

async function removeFromLiked(id: string) {
  await removeLiked(id).catch(() => {});
  items.value = items.value.filter((i) => i.id !== id);
  if (selected.value?.id === id) selected.value = null;
}

onMounted(async () => {
  try {
    const { data } = await getLiked();
    items.value = data.map((i) => i.data);
  } catch {
    items.value = [];
  } finally {
    loading.value = false;
  }
});
</script>

<template>
  <div class="w-full max-w-2xl px-4 flex flex-col gap-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-white font-bold text-2xl">Мои катгёрлс</h2>
        <p class="text-white/40 text-sm mt-0.5">{{ items.length }} лайкнуто</p>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-16">
      <div
        class="w-10 h-10 rounded-full border-4 border-white/20 border-t-white animate-spin"
      />
    </div>

    <!-- Empty -->
    <div
      v-else-if="!items.length"
      class="flex flex-col items-center gap-4 py-20 text-center"
    >
      <div class="text-5xl">🐱</div>
      <p class="text-white/50 text-base">Ты ещё никого не лайкнул</p>
      <button
        class="px-6 py-2.5 rounded-2xl bg-white/10 text-white text-sm font-semibold hover:bg-white/20 transition border border-white/10"
        @click="router.push('/')"
      >
        Начать свайпать
      </button>
    </div>

    <!-- Grid -->
    <div v-else-if="items.length" class="grid grid-cols-2 gap-3">
      <div
        v-for="item in items"
        :key="item.id"
        class="group relative aspect-[3/4] rounded-2xl overflow-hidden cursor-pointer"
        :style="{ background: item.colors.main }"
        @click="selected = item"
      >
        <!-- skeleton -->
        <div class="absolute inset-0 animate-pulse bg-white/5" />
        <img
          :src="item.image.compressed?.url || item.image.original.url"
          :alt="item.anime.character ?? 'cat girl'"
          class="absolute inset-0 w-full h-full object-cover transition duration-300 group-hover:scale-105"
          draggable="false"
          loading="lazy"
        />
        <!-- overlay -->
        <div
          class="absolute inset-0 bg-gradient-to-t from-black/70 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition duration-300"
        />
        <!-- info on hover -->
        <div
          class="absolute bottom-0 left-0 right-0 p-3 translate-y-2 opacity-0 group-hover:translate-y-0 group-hover:opacity-100 transition duration-300"
        >
          <p class="text-white text-xs font-semibold truncate">
            {{ item.anime.character ?? "Неизвестно" }}
          </p>
          <p class="text-white/60 text-[10px] truncate">
            {{ item.anime.title ?? item.category }}
          </p>
        </div>
        <!-- remove btn -->
        <button
          class="absolute top-2 right-2 w-7 h-7 rounded-full bg-black/50 backdrop-blur flex items-center justify-center text-white/70 hover:text-red-400 hover:bg-black/70 transition opacity-0 group-hover:opacity-100 text-xs"
          @click.stop="removeFromLiked(item.id)"
          title="Убрать"
        >
          ✕
        </button>
      </div>
    </div>
  </div>

  <!-- Modal lightbox -->
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="selected"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        style="background: rgba(0, 0, 0, 0.85); backdrop-filter: blur(12px)"
        @click.self="selected = null"
      >
        <div
          class="relative w-full max-w-sm rounded-3xl overflow-hidden shadow-2xl animate-card-enter"
          :style="{ background: selected.colors.main }"
        >
          <img
            :src="
              selected!.image.compressed?.url || selected!.image.original.url
            "
            :alt="selected.anime.character ?? 'cat girl'"
            class="w-full aspect-[3/4] object-cover"
            loading="eager"
          />
          <div
            class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/10 to-transparent pointer-events-none"
          />

          <!-- close -->
          <button
            class="absolute top-4 right-4 w-9 h-9 rounded-full bg-black/50 backdrop-blur text-white flex items-center justify-center hover:bg-black/70 transition text-lg"
            @click="selected = null"
          >
            ✕
          </button>

          <!-- info -->
          <div class="absolute bottom-0 left-0 right-0 p-5">
            <h3 class="text-white font-bold text-xl">
              {{ selected.anime.character ?? "Неизвестно" }}
            </h3>
            <p class="text-white/60 text-sm mt-0.5">
              {{ selected.anime.title ?? selected.category }}
            </p>

            <div class="flex flex-wrap gap-1.5 mt-3">
              <span
                v-for="tag in selected.tags.slice(0, 5)"
                :key="tag"
                class="text-[10px] font-semibold px-2 py-0.5 rounded-full bg-white/10 text-white/70"
              >
                #{{ tag }}
              </span>
            </div>

            <!-- artist -->
            <div
              v-if="selected.attribution.artist.username"
              class="mt-3 flex items-center gap-2"
            >
              <span class="text-white/30 text-xs">Арт:</span>
              <a
                v-if="selected.attribution.artist.profile"
                :href="selected.attribution.artist.profile"
                target="_blank"
                class="text-pink-400 text-xs hover:underline"
              >
                @{{ selected.attribution.artist.username }}
              </a>
              <span v-else class="text-white/50 text-xs"
                >@{{ selected.attribution.artist.username }}</span
              >
            </div>

            <!-- actions -->
            <div class="mt-4 flex flex-col gap-2">
              <button
                class="w-full py-2.5 rounded-xl bg-pink-500/20 border border-pink-400/30 text-pink-300 text-sm font-semibold hover:bg-pink-500/30 transition flex items-center justify-center gap-2"
                @click="
                  router.push({
                    path: `/chat/${selected!.id}`,
                    query: {
                      img:
                        selected!.image.compressed?.url ||
                        selected!.image.original.url,
                      name: selected!.anime.character || '',
                    },
                  });
                  selected = null;
                "
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="w-4 h-4"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                  />
                </svg>
                Перейти в чат
              </button>
              <button
                class="w-full py-2.5 rounded-xl bg-red-500/20 border border-red-400/30 text-red-300 text-sm font-semibold hover:bg-red-500/30 transition"
                @click="removeFromLiked(selected!.id)"
              >
                Убрать из лайков
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
