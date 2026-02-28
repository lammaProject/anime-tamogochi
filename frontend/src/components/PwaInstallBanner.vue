<script setup lang="ts">
import { usePwa } from "@/composables/usePwa";

const { isInstallable, install, dismiss } = usePwa();
</script>

<template>
  <Transition
    enter-active-class="transition duration-300 ease-out"
    enter-from-class="translate-y-full opacity-0"
    enter-to-class="translate-y-0 opacity-100"
    leave-active-class="transition duration-200 ease-in"
    leave-from-class="translate-y-0 opacity-100"
    leave-to-class="translate-y-full opacity-0"
  >
    <div
      v-if="isInstallable"
      class="fixed bottom-4 left-4 right-4 z-50 mx-auto max-w-md rounded-2xl border border-white/10 bg-[#1a1a2e]/95 backdrop-blur-xl p-4 shadow-2xl"
    >
      <div class="flex items-center gap-3">
        <img src="/icon.svg" alt="nekoswipe" class="w-10 h-10 rounded-xl" />
        <div class="flex-1 min-w-0">
          <p class="text-white text-sm font-semibold">
            Установить nekoswipe
          </p>
          <p class="text-white/40 text-xs">
            Добавить на главный экран
          </p>
        </div>
        <div class="flex gap-2 shrink-0">
          <button
            @click="dismiss"
            class="px-3 py-1.5 rounded-xl text-xs text-white/40 hover:text-white/60 transition"
          >
            Потом
          </button>
          <button
            @click="install"
            class="px-4 py-1.5 rounded-xl bg-pink-500 text-white text-xs font-semibold hover:bg-pink-400 transition"
          >
            Установить
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>
