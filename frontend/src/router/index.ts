import { createRouter, createWebHashHistory } from "vue-router";
import { useAuthStore } from "@/stores/authStore";
import HomePage from "@/pages/HomePage.vue";
import LikedPage from "@/pages/LikedPage.vue";
import LoginPage from "@/pages/LoginPage.vue";
import ChatPage from "@/pages/ChatPage.vue";
import ProfilePage from "@/pages/ProfilePage.vue";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: "/", component: HomePage, meta: { auth: true } },
    { path: "/liked", component: LikedPage, meta: { auth: true } },
    { path: "/profile", component: ProfilePage, meta: { auth: true } },
    { path: "/chat/:id", component: ChatPage, meta: { auth: true } },
    { path: "/login", component: LoginPage, meta: { auth: false } },
  ],
});

router.beforeEach(async (to) => {
  const auth = useAuthStore();

  if (auth.loading) await auth.init();

  if (to.meta.auth && !auth.isLoggedIn) return "/login";
  if (to.path === "/login" && auth.isLoggedIn) return "/";
});

export default router;
