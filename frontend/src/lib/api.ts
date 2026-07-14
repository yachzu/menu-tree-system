import {
  MenuItem,
  CreateMenuRequest,
  UpdateMenuRequest,
  MoveMenuRequest,
  ReorderMenuRequest,
} from "./types";

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api";

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    path: string,
    options?: RequestInit & { signal?: AbortSignal },
  ): Promise<T> {
    const res = await fetch(`${this.baseUrl}${path}`, {
      headers: { "Content-Type": "application/json" },
      ...options,
    });

    if (!res.ok) {
      const error = await res.json().catch(() => ({
        error: "unknown_error",
        message: res.statusText,
      }));
      throw new Error(error.message || "Request failed");
    }

    if (res.status === 204) return undefined as T;

    const json = await res.json();
    return json.data ?? json;
  }

  async getTree(signal?: AbortSignal): Promise<MenuItem[]> {
    return this.request<MenuItem[]>("/menus", { signal });
  }

  async getById(id: string, signal?: AbortSignal): Promise<MenuItem> {
    return this.request<MenuItem>(`/menus/${id}`, { signal });
  }

  async create(
    data: CreateMenuRequest,
    signal?: AbortSignal,
  ): Promise<MenuItem> {
    return this.request<MenuItem>("/menus", {
      method: "POST",
      body: JSON.stringify(data),
      signal,
    });
  }

  async update(
    id: string,
    data: UpdateMenuRequest,
    signal?: AbortSignal,
  ): Promise<MenuItem> {
    return this.request<MenuItem>(`/menus/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
      signal,
    });
  }

  async delete(id: string, signal?: AbortSignal): Promise<void> {
    return this.request<void>(`/menus/${id}`, {
      method: "DELETE",
      signal,
    });
  }

  async move(
    id: string,
    data: MoveMenuRequest,
    signal?: AbortSignal,
  ): Promise<MenuItem> {
    return this.request<MenuItem>(`/menus/${id}/move`, {
      method: "PATCH",
      body: JSON.stringify(data),
      signal,
    });
  }

  async reorder(
    id: string,
    data: ReorderMenuRequest,
    signal?: AbortSignal,
  ): Promise<MenuItem> {
    return this.request<MenuItem>(`/menus/${id}/reorder`, {
      method: "PATCH",
      body: JSON.stringify(data),
      signal,
    });
  }
}

export const api = new ApiClient(API_BASE);
