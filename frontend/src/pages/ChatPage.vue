<script setup lang="ts">
import { ref, onMounted, nextTick } from "vue";
import { useRoute, useRouter } from "vue-router";
import { getChatHistory, sendChatMessage } from "@/api/backend";
import type { ChatMessage } from "@/api/backend";

const route     = useRoute();
const router    = useRouter();
const catgirlId = route.params.id as string;

// Аватарка и имя из query-параметров (переданы из LikedPage)
const avatarUrl  = ref((route.query.img as string) || "");
const girlName   = ref((route.query.name as string) || "Кошкодевочка");

const messages  = ref<ChatMessage[]>([]);
const input     = ref("");
const loading   = ref(true);
const sending   = ref(false);
const listRef   = ref<HTMLElement | null>(null);

function scrollBottom() {
  nextTick(() => {
    if (listRef.value) listRef.value.scrollTop = listRef.value.scrollHeight;
  });
}

onMounted(async () => {
  try {
    const { data } = await getChatHistory(catgirlId);
    messages.value = data ?? [];
  } catch {
    messages.value = [];
  } finally {
    loading.value = false;
    scrollBottom();
  }
});

async function send() {
  const text = input.value.trim();
  if (!text || sending.value) return;

  // Оптимистично добавляем сообщение пользователя
  const optimistic: ChatMessage = {
    content: text,
    from_user: true,
    created_at: new Date().toISOString(),
  };
  messages.value.push(optimistic);
  input.value = "";
  sending.value = true;
  scrollBottom();

  try {
    const { data } = await sendChatMessage(catgirlId, text);
    // Заменяем оптимистичное сообщение на реальное + добавляем ответ бота
    messages.value.splice(messages.value.length - 1, 1, data.user_message);
    messages.value.push(data.bot_message);
    scrollBottom();
  } catch {
    // Если ошибка — убираем оптимистичное сообщение
    messages.value.pop();
  } finally {
    sending.value = false;
  }
}

function formatTime(iso: string) {
  return new Date(iso).toLocaleTimeString("ru", { hour: "2-digit", minute: "2-digit" });
}
</script>

<template>
  <div class="w-full max-w-sm px-4 flex flex-col" style="height: calc(100vh - 80px); padding-top: 16px;">

    <!-- Header -->
    <div class="flex items-center gap-3 mb-4">
      <button
        class="w-8 h-8 rounded-full bg-white/5 flex items-center justify-center text-white/50 hover:text-white transition flex-shrink-0"
        @click="router.push('/liked')"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
        </svg>
      </button>

      <!-- Аватарка -->
      <div class="relative flex-shrink-0">
        <div
          class="w-10 h-10 rounded-full overflow-hidden border-2 flex items-center justify-center text-xl"
          :style="avatarUrl ? 'border-color: rgba(236,72,153,0.5)' : 'border-color: rgba(255,255,255,0.1)'"
        >
          <img
            v-if="avatarUrl"
            :src="avatarUrl"
            class="w-full h-full object-cover"
            alt="avatar"
          />
          <span v-else>🐱</span>
        </div>
        <!-- Онлайн-индикатор -->
        <div class="absolute bottom-0 right-0 w-2.5 h-2.5 rounded-full bg-green-400 border-2 border-black" />
      </div>

      <!-- Имя и статус -->
      <div class="flex-1 min-w-0">
        <p class="text-white font-semibold text-sm truncate">{{ girlName || catgirlId }}</p>
        <p class="text-green-400 text-xs">онлайн</p>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex-1 flex items-center justify-center">
      <div class="w-8 h-8 rounded-full border-4 border-white/20 border-t-white animate-spin" />
    </div>

    <!-- Messages -->
    <div
      v-else
      ref="listRef"
      class="flex-1 overflow-y-auto flex flex-col gap-2 pr-1 pt-2 scrollbar-none"
      style="scrollbar-width: none;"
    >
      <!-- Empty state -->
      <div v-if="!messages.length && !sending" class="flex flex-col items-center justify-center h-full gap-2 text-center">
        <div class="text-4xl">🐱</div>
        <p class="text-white/30 text-sm">Напиши что-нибудь, ня~</p>
      </div>

      <!-- Message bubbles -->
      <div
        v-for="(msg, i) in messages"
        :key="i"
        class="flex items-end gap-2"
        :class="msg.from_user ? 'justify-end' : 'justify-start'"
      >
        <!-- Аватарка бота слева -->
        <div
          v-if="!msg.from_user"
          class="w-7 h-7 rounded-full overflow-hidden flex-shrink-0 flex items-center justify-center text-sm mb-0.5"
          style="background: rgba(255,255,255,0.05); border: 1px solid rgba(255,255,255,0.1)"
        >
          <img v-if="avatarUrl" :src="avatarUrl" class="w-full h-full object-cover" alt="" />
          <span v-else>🐱</span>
        </div>

        <div
          class="max-w-[75%] px-4 py-2.5 rounded-2xl text-sm leading-relaxed"
          :class="msg.from_user ? 'rounded-br-sm text-white' : 'rounded-bl-sm text-white'"
          :style="msg.from_user
            ? 'background: linear-gradient(135deg, #7c3aed, #ec4899);'
            : 'background: rgba(255,255,255,0.08); border: 1px solid rgba(255,255,255,0.1);'"
        >
          <p>{{ msg.content }}</p>
          <p class="text-[10px] mt-1 opacity-50 text-right">{{ formatTime(msg.created_at) }}</p>
        </div>
      </div>

      <!-- Typing indicator (пока ждём ответ) -->
      <div v-if="sending" class="flex justify-start items-end gap-2">
        <div
          class="w-7 h-7 rounded-full overflow-hidden flex-shrink-0 flex items-center justify-center text-sm"
          style="background: rgba(255,255,255,0.05); border: 1px solid rgba(255,255,255,0.1)"
        >
          <img v-if="avatarUrl" :src="avatarUrl" class="w-full h-full object-cover" alt="" />
          <span v-else>🐱</span>
        </div>
        <div class="px-4 py-3 rounded-2xl rounded-bl-sm flex gap-1 items-center" style="background: rgba(255,255,255,0.08); border: 1px solid rgba(255,255,255,0.1);">
          <span class="w-1.5 h-1.5 rounded-full bg-white/50 animate-bounce" style="animation-delay: 0ms" />
          <span class="w-1.5 h-1.5 rounded-full bg-white/50 animate-bounce" style="animation-delay: 150ms" />
          <span class="w-1.5 h-1.5 rounded-full bg-white/50 animate-bounce" style="animation-delay: 300ms" />
        </div>
      </div>
    </div>

    <!-- Input -->
    <div class="mt-3 flex gap-2 items-end">
      <input
        v-model="input"
        type="text"
        placeholder="Напиши сообщение..."
        class="flex-1 px-4 py-3 rounded-2xl text-sm text-white placeholder-white/30 outline-none focus:ring-2 focus:ring-pink-400/50 transition"
        style="background: rgba(255,255,255,0.06); border: 1px solid rgba(255,255,255,0.1);"
        :disabled="sending"
        @keyup.enter="send"
      />
      <button
        class="w-12 h-12 rounded-2xl flex items-center justify-center transition active:scale-95 flex-shrink-0"
        :class="(input.trim() && !sending) ? 'opacity-100' : 'opacity-40'"
        style="background: linear-gradient(135deg, #7c3aed, #ec4899);"
        :disabled="sending"
        @click="send"
      >
        <svg v-if="!sending" xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
        </svg>
        <div v-else class="w-4 h-4 border-2 border-white/40 border-t-white rounded-full animate-spin" />
      </button>
    </div>
  </div>
</template>
