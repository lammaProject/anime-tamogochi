import axios from "axios";

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL ?? "",
  withCredentials: true, // куки
});

// Auth
export const register = (username: string, password: string) =>
  api.post("/api/register", { username, password });

export const login = (username: string, password: string) =>
  api.post("/api/login", { username, password });

export const logout = () => api.post("/api/logout");

export const getMe = () => api.get<{ id: number; username: string }>("/api/me");

// Liked
export const getLiked = () => api.get<LikedItem[]>("/api/liked");

export const addLiked = (catgirlId: string, data: object) =>
  api.post("/api/liked", { catgirl_id: catgirlId, data });

export const removeLiked = (catgirlId: string) =>
  api.delete(`/api/liked/${catgirlId}`);

export interface LikedItem {
  catgirl_id: string;
  data: import("./type").NekosImageData;
}

// Chat
export interface ChatMessage {
  content: string;
  from_user: boolean;
  created_at: string;
}

export const getChatHistory = (catgirlId: string) =>
  api.get<ChatMessage[]>(`/api/chat/${catgirlId}`);

export const sendChatMessage = (catgirlId: string, content: string) =>
  api.post<{ user_message: ChatMessage; bot_message: ChatMessage }>(`/api/chat/${catgirlId}`, { content });
