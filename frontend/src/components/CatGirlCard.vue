<script setup lang="ts">
import { computed, ref, watch, watchEffect } from "vue";
import { getCatGirls } from "@/api/api";
import type { NekosImageData } from "@/api/type";
import { useCardStore } from "@/stores/cardStore";
import { addLiked } from "@/api/backend";
import { useQuery } from "@tanstack/vue-query";

const cardStore = useCardStore();

// ─── state ───────────────────────────────────────────────────────────────────
const images = ref<NekosImageData[]>([]);
const loading = ref(true);
const error = ref<string | null>(null);
const currentIndex = ref(0);
const finished = ref(false);

currentIndex.value = cardStore.currentIndex;
// ─── preload ──────────────────────────────────────────────────────────────────
const PRELOAD_AHEAD = 6;
const preloaded = new Set<string>();

function preloadImages(fromIndex: number) {
  for (
    let i = fromIndex;
    i < Math.min(fromIndex + PRELOAD_AHEAD, images.value.length);
    i++
  ) {
    const img = images.value[i];
    if (!img) continue;

    const url = img.image.compressed?.url || img.image.original.url;
    if (preloaded.has(url)) continue;
    preloaded.add(url);

    // Первая следующая — через <link rel="preload"> (высший приоритет браузера)
    if (i === fromIndex) {
      const link = document.createElement("link");
      link.rel = "preload";
      link.as = "image";
      link.fetchPriority = "high";
      link.href = url;
      document.head.appendChild(link);
    } else {
      // Остальные — через Image() в фоне
      const el = new window.Image();
      el.src = url;
    }
  }
}

// Трекаем загружена ли текущая картинка
const imgLoaded = ref(false);

function onImgLoad() {
  imgLoaded.value = true;
}
function onImgError() {
  imgLoaded.value = true; // скрываем скелетон даже при ошибке
}

// drag
const isDragging = ref(false);
const dragStartX = ref(0);
const dragOffsetX = ref(0); // px, horizontal
const isFlying = ref(false); // анимация вылета идёт
const flyDir = ref<"left" | "right" | null>(null);

// ─── порог ───────────────────────────────────────────────────────────────────
const THRESHOLD = 120; // px — после этого засчитывается

// ─── вычисляемые ─────────────────────────────────────────────────────────────
const current = computed(() => images.value[currentIndex.value] ?? null);
const bgColor = computed(() => current.value?.colors?.main ?? "#1a1a2e");

// Обновляем глобальный цвет фона при смене карточки + сбрасываем флаг загрузки
watch(
  current,
  (card) => {
    if (card) {
      cardStore.setCurrentCard(
        card.colors.main,
        card.colors.palette,
        currentIndex.value,
      );
      imgLoaded.value = false;
    }
  },
  { immediate: true },
);

// прогресс от 0 до 1 (насколько далеко от порога)
const dragProgress = computed(() =>
  Math.min(Math.abs(dragOffsetX.value) / THRESHOLD, 1),
);

const showLabel = computed<"like" | "nope" | null>(() => {
  if (isFlying.value) return flyDir.value === "right" ? "like" : "nope";
  if (Math.abs(dragOffsetX.value) > 60)
    return dragOffsetX.value > 0 ? "like" : "nope";
  return null;
});

// transform карточки во время перетаскивания
const cardTransform = computed(() => {
  if (isFlying.value) return ""; // управляется CSS-анимацией

  if (!isDragging.value && dragOffsetX.value === 0) return "";

  const x = dragOffsetX.value;
  const rotate = x * 0.07; // наклон пропорционален отводу
  // опускаем вниз по квадратичной кривой — max ~80px при пороге
  const y = dragProgress.value ** 1.5 * 80;
  // немного уменьшаем масштаб чем дальше
  const scale = 1 - dragProgress.value * 0.04;

  return `translateX(${x}px) translateY(${y}px) rotate(${rotate}deg) scale(${scale})`;
});

const cardStyle = computed(() => {
  if (isFlying.value) return {};
  return {
    transform: cardTransform.value,
    transition: isDragging.value
      ? "none"
      : "transform 0.35s cubic-bezier(0.34, 1.56, 0.64, 1)",
  };
});

// тень и оверлей меняются по прогрессу
const cardShadow = computed(() => {
  const p = dragProgress.value;
  const color =
    dragOffsetX.value > 0
      ? `rgba(74, 222, 128, ${p * 0.4})`
      : `rgba(248, 113, 113, ${p * 0.4})`;
  return `0 25px 60px -10px rgba(0,0,0,0.5), 0 0 0 3px ${color}`;
});

// ─── API liked ────────────────────────────────────────────────────────────────
function saveLike(item: NekosImageData) {
  // fire-and-forget, не блокируем UI
  addLiked(item.id, item).catch(() => {});
}

// ─── свайп ────────────────────────────────────────────────────────────────────
function commitSwipe(dir: "left" | "right") {
  if (isFlying.value || !current.value) return;

  if (dir === "right") {
    saveLike(current.value);
    likedCount.value++;
  }

  flyDir.value = dir;
  isFlying.value = true;

  setTimeout(() => {
    currentIndex.value++;
    isFlying.value = false;
    flyDir.value = null;
    dragOffsetX.value = 0;
    isDragging.value = false;

    if (currentIndex.value >= images.value.length) {
      finished.value = true;
    } else {
      // Предзагружаем следующие N после текущей
      preloadImages(currentIndex.value + 1);
    }
  }, 420);
}

// ─── drag / pointer ──────────────────────────────────────────────────────────
function onPointerDown(e: PointerEvent) {
  if (isFlying.value) return;
  isDragging.value = true;
  dragStartX.value = e.clientX;
  dragOffsetX.value = 0;
  (e.currentTarget as HTMLElement).setPointerCapture(e.pointerId);
}

function onPointerMove(e: PointerEvent) {
  if (!isDragging.value || isFlying.value) return;
  e.preventDefault(); // блокируем скролл страницы на мобилке
  dragOffsetX.value = e.clientX - dragStartX.value;
}

function onPointerUp() {
  if (!isDragging.value || isFlying.value) return;
  isDragging.value = false;

  if (dragOffsetX.value > THRESHOLD) {
    commitSwipe("right");
  } else if (dragOffsetX.value < -THRESHOLD) {
    commitSwipe("left");
  } else {
    // snap back
    dragOffsetX.value = 0;
  }
}

const { data, refetch } = useQuery({
  queryKey: ["cards"],
  queryFn: getCatGirls,
});

watchEffect(() => {
  if (data.value?.images) {
    images.value = data.value.images;
    preloadImages(0);
    loading.value = false;
  }
});

const likedCount = ref(0); // обновляется в процессе свайпов
</script>

<template>
  <!-- Loading -->
  <div v-if="loading" class="flex flex-col items-center gap-4">
    <div
      class="w-16 h-16 rounded-full border-4 border-white/20 border-t-white animate-spin"
    />
    <p class="text-white/60 text-sm font-medium tracking-widest uppercase">
      Загрузка...
    </p>
  </div>

  <!-- Error -->
  <div v-else-if="error" class="text-red-400 font-semibold text-lg">
    {{ error }}
  </div>

  <!-- Finished -->
  <div
    v-else-if="finished"
    class="flex flex-col items-center gap-6 text-center px-6"
  >
    <div class="text-6xl">✨</div>
    <h2 class="text-white text-3xl font-bold">Всё!</h2>
    <p class="text-white/60 text-base">
      Ты лайкнул
      <span class="text-pink-400 font-bold">{{ likedCount }}</span> катгёрлс
    </p>
    <button
      class="mt-2 px-8 py-3 rounded-2xl bg-white/10 backdrop-blur text-white font-semibold hover:bg-white/20 transition border border-white/20"
      @click="
        async () => {
          await refetch();
          finished = false;
          currentIndex = 0;
        }
      "
    >
      Сначала
    </button>
  </div>

  <!-- Card deck -->
  <div
    v-else
    class="relative flex flex-col items-center gap-6 select-none w-full max-w-sm px-4"
  >
    <!-- Counter -->
    <div
      class="flex items-center gap-2 text-white/50 text-xs font-medium tracking-widest uppercase"
    >
      <span>{{ currentIndex + 1 }}</span>
      <span>/</span>
      <span>{{ images.length }}</span>
    </div>

    <!-- Tinder card -->
    <div
      :key="currentIndex"
      class="relative w-full aspect-[3/4] rounded-3xl overflow-hidden cursor-grab active:cursor-grabbing"
      style="touch-action: none"
      :class="{
        'animate-swipe-left': isFlying && flyDir === 'left',
        'animate-swipe-right': isFlying && flyDir === 'right',
        'animate-card-enter': !isFlying,
      }"
      :style="{
        ...cardStyle,
        background: bgColor,
        boxShadow: cardShadow,
      }"
      @pointerdown="onPointerDown"
      @pointermove="onPointerMove"
      @pointerup="onPointerUp"
      @pointercancel="onPointerUp"
    >
      <!-- Skeleton — виден пока картинка не загрузилась -->
      <div
        class="absolute inset-0 transition-opacity duration-300"
        :class="imgLoaded ? 'opacity-0' : 'opacity-100'"
        :style="{
          background: `linear-gradient(135deg, ${bgColor}cc, ${bgColor}66)`,
        }"
      >
        <div class="absolute inset-0 animate-pulse bg-white/5" />
        <!-- Иконка-заглушка по центру -->
        <div
          class="absolute inset-0 flex items-center justify-center text-5xl opacity-20 animate-pulse"
        >
          🐱
        </div>
      </div>

      <!-- Image — используем compressed (лёгкий), preload стартует заранее -->
      <img
        :src="current!.image.compressed?.url || current!.image.original.url"
        :alt="current!.anime.character ?? 'cat girl'"
        class="absolute inset-0 w-full h-full object-cover pointer-events-none transition-opacity duration-300"
        :class="imgLoaded ? 'opacity-100' : 'opacity-0'"
        draggable="false"
        loading="eager"
        fetchpriority="high"
        @load="onImgLoad"
        @error="onImgError"
      />

      <!-- Gradient overlay -->
      <div
        class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/10 to-transparent pointer-events-none"
      />

      <!-- Colored tint overlay based on drag direction -->
      <div
        v-if="dragProgress > 0.1"
        class="absolute inset-0 pointer-events-none transition-opacity"
        :style="{
          background:
            dragOffsetX > 0
              ? `rgba(74, 222, 128, ${dragProgress * 0.25})`
              : `rgba(248, 113, 113, ${dragProgress * 0.25})`,
        }"
      />

      <!-- LIKE label -->
      <div
        v-if="showLabel === 'like'"
        class="animate-pop-like absolute top-8 left-6 border-4 border-green-400 text-green-400 font-black text-3xl px-4 py-1 rounded-xl rotate-[-12deg] tracking-widest pointer-events-none"
        :style="{ opacity: Math.min(dragProgress * 2, 1) }"
      >
        LIKE
      </div>

      <!-- NOPE label -->
      <div
        v-if="showLabel === 'nope'"
        class="animate-pop-nope absolute top-8 right-6 border-4 border-red-400 text-red-400 font-black text-3xl px-4 py-1 rounded-xl rotate-[12deg] tracking-widest pointer-events-none"
        :style="{ opacity: Math.min(dragProgress * 2, 1) }"
      >
        NOPE
      </div>

      <!-- Card info -->
      <div class="absolute bottom-0 left-0 right-0 p-5 pointer-events-none">
        <div class="flex items-end justify-between">
          <div>
            <h3 class="text-white font-bold text-xl leading-tight">
              {{ current!.anime.character ?? "Неизвестно" }}
            </h3>
            <p class="text-white/60 text-sm mt-0.5">
              {{ current!.anime.title ?? current!.category }}
            </p>
          </div>
        </div>
        <div v-if="current!.tags.length" class="flex flex-wrap gap-1.5 mt-3">
          <span
            v-for="tag in current!.tags.slice(0, 4)"
            :key="tag"
            class="text-[10px] font-semibold px-2 py-0.5 rounded-full bg-white/10 text-white/70 backdrop-blur-sm"
          >
            #{{ tag }}
          </span>
        </div>
      </div>
    </div>

    <!-- Progress bar -->
    <div class="w-full h-1 rounded-full bg-white/5 overflow-hidden">
      <div
        class="h-full rounded-full transition-all duration-300"
        :style="{
          width: `${(currentIndex / images.length) * 100}%`,
          background: 'linear-gradient(90deg, #7c3aed, #ec4899)',
        }"
      />
    </div>

    <!-- Buttons -->
    <div class="flex items-center gap-6">
      <button
        class="group w-16 h-16 rounded-full bg-white/5 border border-white/10 backdrop-blur flex items-center justify-center text-2xl hover:bg-red-500/20 hover:border-red-400/50 hover:scale-110 transition-all duration-200 active:scale-95"
        @click="commitSwipe('left')"
      >
        <span class="group-hover:scale-125 transition-transform">❌</span>
      </button>

      <button
        class="group w-16 h-16 rounded-full bg-white/5 border border-white/10 backdrop-blur flex items-center justify-center text-2xl hover:bg-green-500/20 hover:border-green-400/50 hover:scale-110 transition-all duration-200 active:scale-95"
        @click="commitSwipe('right')"
      >
        <span class="group-hover:scale-125 transition-transform">💚</span>
      </button>
    </div>

    <p class="text-white/25 text-xs tracking-wider">
      ← тяни карточку или нажми кнопки →
    </p>
  </div>
</template>
