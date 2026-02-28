import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "./App.vue";
import router from "./router/index.ts";
import "./style.css";
import { VueQueryPlugin } from "@tanstack/vue-query";

createApp(App)
  .use(createPinia())
  .use(router)
  .use(VueQueryPlugin, {
    queryClientConfig: {
      defaultOptions: {
        queries: { refetchOnWindowFocus: false, refetchOnMount: false },
      },
    },
  })
  .mount("#app");
