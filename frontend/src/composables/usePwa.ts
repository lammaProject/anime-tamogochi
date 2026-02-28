import { ref, onMounted, onBeforeUnmount } from "vue";

interface BeforeInstallPromptEvent extends Event {
  prompt(): Promise<void>;
  userChoice: Promise<{ outcome: "accepted" | "dismissed" }>;
}

const deferredPrompt = ref<BeforeInstallPromptEvent | null>(null);
const isInstallable = ref(false);
const isInstalled = ref(false);

export function usePwa() {
  function onBeforeInstall(e: Event) {
    e.preventDefault();
    deferredPrompt.value = e as BeforeInstallPromptEvent;
    isInstallable.value = true;
  }

  function onAppInstalled() {
    isInstalled.value = true;
    isInstallable.value = false;
    deferredPrompt.value = null;
  }

  async function install() {
    if (!deferredPrompt.value) return;
    await deferredPrompt.value.prompt();
    const { outcome } = await deferredPrompt.value.userChoice;
    if (outcome === "accepted") {
      isInstalled.value = true;
      isInstallable.value = false;
    }
    deferredPrompt.value = null;
  }

  function dismiss() {
    isInstallable.value = false;
  }

  onMounted(() => {
    // Check if already installed (standalone mode)
    if (
      window.matchMedia("(display-mode: standalone)").matches ||
      (navigator as any).standalone
    ) {
      isInstalled.value = true;
      return;
    }

    window.addEventListener("beforeinstallprompt", onBeforeInstall);
    window.addEventListener("appinstalled", onAppInstalled);
  });

  onBeforeUnmount(() => {
    window.removeEventListener("beforeinstallprompt", onBeforeInstall);
    window.removeEventListener("appinstalled", onAppInstalled);
  });

  return {
    isInstallable,
    isInstalled,
    install,
    dismiss,
  };
}
