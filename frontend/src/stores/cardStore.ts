import { defineStore } from "pinia";
import { computed, ref } from "vue";

export const useCardStore = defineStore("card", () => {
  const mainColor = ref("#0d0d0d");
  const palette = ref<string[]>([]);
  const currentIndex = ref(0);

  // Берём второй и третий цвет из palette для blob-ов фона
  const blobColor1 = computed(() => palette.value[1] ?? mainColor.value);
  const blobColor2 = computed(() => palette.value[3] ?? mainColor.value);
  const blobColor3 = computed(() => palette.value[5] ?? mainColor.value);

  function setCurrentCard(
    color: string,
    colorPalette: string[],
    currentIndexCard: number,
  ) {
    mainColor.value = color;
    palette.value = colorPalette;
    currentIndex.value = currentIndexCard;
  }

  return {
    mainColor,
    palette,
    blobColor1,
    blobColor2,
    blobColor3,
    setCurrentCard,
    currentIndex,
  };
});
