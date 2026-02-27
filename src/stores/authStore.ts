import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { getMe, login as apiLogin, register as apiRegister, logout as apiLogout } from "@/api/backend";

export const useAuthStore = defineStore("auth", () => {
  const user = ref<{ id: number; username: string } | null>(null);
  const loading = ref(true);

  const isLoggedIn = computed(() => !!user.value);

  async function init() {
    loading.value = true;
    try {
      const { data } = await getMe();
      user.value = data;
    } catch {
      user.value = null;
    } finally {
      loading.value = false;
    }
  }

  async function login(username: string, password: string) {
    const { data } = await apiLogin(username, password);
    user.value = data;
  }

  async function register(username: string, password: string) {
    const { data } = await apiRegister(username, password);
    user.value = data;
  }

  async function logout() {
    await apiLogout();
    user.value = null;
  }

  return { user, loading, isLoggedIn, init, login, register, logout };
});
