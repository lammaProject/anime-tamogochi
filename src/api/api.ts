import axios from "axios";
import type { ResponseGetGirls } from "./type.ts";

const api = axios.create({
  baseURL: "https://api.nekosia.cat/api/v1",
});

export const getCatGirls = async () => {
  const { data } = await api.get<ResponseGetGirls>("/images/catgirl?count=20");
  return data;
};
