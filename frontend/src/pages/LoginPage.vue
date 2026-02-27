<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "@/stores/authStore";

const auth   = useAuthStore();
const router = useRouter();

const mode     = ref<"login" | "register">("login");
const username = ref("");
const password = ref("");
const error    = ref("");
const loading  = ref(false);

async function submit() {
  if (!username.value || !password.value) {
    error.value = "Заполни все поля";
    return;
  }
  loading.value = true;
  error.value   = "";
  try {
    if (mode.value === "login") {
      await auth.login(username.value, password.value);
    } else {
      await auth.register(username.value, password.value);
    }
    router.push("/");
  } catch (e: any) {
    error.value = e?.response?.data?.error ?? "Ошибка";
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div class="w-full max-w-sm px-4 flex flex-col gap-6">
    <!-- Logo -->
    <div class="text-center">
      <h1 class="text-white font-bold text-4xl tracking-tight">
        neko<span style="color: #ec4899">swipe</span>
      </h1>
      <p class="text-white/30 text-xs mt-1 tracking-[0.3em] uppercase">find your catgirl</p>
    </div>

    <!-- Card -->
    <div class="rounded-3xl p-6 flex flex-col gap-4" style="background: rgba(255,255,255,0.04); border: 1px solid rgba(255,255,255,0.08);">
      <!-- Tabs -->
      <div class="flex rounded-2xl overflow-hidden bg-white/5 p-1 gap-1">
        <button
          class="flex-1 py-2 rounded-xl text-sm font-semibold transition"
          :class="mode === 'login' ? 'bg-white/10 text-white' : 'text-white/40 hover:text-white/60'"
          @click="mode = 'login'; error = ''"
        >
          Войти
        </button>
        <button
          class="flex-1 py-2 rounded-xl text-sm font-semibold transition"
          :class="mode === 'register' ? 'bg-white/10 text-white' : 'text-white/40 hover:text-white/60'"
          @click="mode = 'register'; error = ''"
        >
          Регистрация
        </button>
      </div>

      <!-- Fields -->
      <div class="flex flex-col gap-3">
        <input
          v-model="username"
          type="text"
          placeholder="Имя пользователя"
          autocomplete="username"
          class="w-full px-4 py-3 rounded-2xl text-sm text-white placeholder-white/30 outline-none focus:ring-2 focus:ring-pink-400/50 transition"
          style="background: rgba(255,255,255,0.06); border: 1px solid rgba(255,255,255,0.1);"
          @keyup.enter="submit"
        />
        <input
          v-model="password"
          type="password"
          placeholder="Пароль"
          autocomplete="current-password"
          class="w-full px-4 py-3 rounded-2xl text-sm text-white placeholder-white/30 outline-none focus:ring-2 focus:ring-pink-400/50 transition"
          style="background: rgba(255,255,255,0.06); border: 1px solid rgba(255,255,255,0.1);"
          @keyup.enter="submit"
        />
      </div>

      <!-- Error -->
      <p v-if="error" class="text-red-400 text-xs text-center">{{ error }}</p>

      <!-- Submit -->
      <button
        class="w-full py-3 rounded-2xl font-semibold text-sm text-white transition active:scale-95 disabled:opacity-50"
        style="background: linear-gradient(135deg, #7c3aed, #ec4899);"
        :disabled="loading"
        @click="submit"
      >
        <span v-if="loading" class="inline-block w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
        <span v-else>{{ mode === 'login' ? 'Войти' : 'Создать аккаунт' }}</span>
      </button>
    </div>
  </div>
</template>
