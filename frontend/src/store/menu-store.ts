"use client";

import { create } from "zustand";
import { MenuItem } from "@/lib/types";
import { api } from "@/lib/api";

interface MenuState {
  menus: MenuItem[];
  tree: MenuItem[];
  searchQuery: string;
  expandedIds: string[];
  isLoading: boolean;
  error: string | null;
  selectedMenu: MenuItem | null;

  fetchTree: (signal?: AbortSignal) => Promise<void>;
  createMenu: (name: string, parentId: string | null) => Promise<void>;
  updateMenu: (id: string, name: string) => Promise<void>;
  deleteMenu: (id: string) => Promise<void>;
  moveMenu: (id: string, newParentId: string | null, orderIndex: number) => Promise<void>;
  reorderMenu: (id: string, orderIndex: number) => Promise<void>;
  setSearchQuery: (query: string) => void;
  toggleExpand: (id: string) => void;
  expandAll: () => void;
  collapseAll: () => void;
  setSelectedMenu: (menu: MenuItem | null) => void;
}

function flattenTree(tree: MenuItem[]): MenuItem[] {
  const result: MenuItem[] = [];
  function walk(items: MenuItem[]) {
    for (const item of items) {
      result.push(item);
      if (item.children) walk(item.children);
    }
  }
  walk(tree);
  return result;
}

const expandedSet = (ids: string[]) => new Set(ids);

export const useMenuStore = create<MenuState>((set, get) => ({
  menus: [],
  tree: [],
  searchQuery: "",
  expandedIds: [],
  isLoading: false,
  error: null,
  selectedMenu: null,

  fetchTree: async (signal) => {
    set({ isLoading: true, error: null });
    try {
      const tree = await api.getTree(signal);
      const menus = flattenTree(tree);
      set(state => ({
        tree,
        menus,
        isLoading: false,
        expandedIds: state.expandedIds.length === 0 && state.menus.length === 0
          ? menus.filter((m) => m.depth < 2).map((m) => m.id)
          : state.expandedIds,
      }));
    } catch (err) {
      if (err instanceof DOMException && err.name === "AbortError") return;
      set({
        isLoading: false,
        error: err instanceof Error ? err.message : "Failed to fetch menu tree",
      });
    }
  },

  createMenu: async (name, parentId) => {
    set({ error: null });
    try {
      await api.create({ name, parent_id: parentId });
      await get().fetchTree();
    } catch (err) {
      set({
        error: err instanceof Error ? err.message : "Failed to create menu",
      });
    }
  },

  updateMenu: async (id, name) => {
    set({ error: null });
    try {
      await api.update(id, { name });
      await get().fetchTree();
    } catch (err) {
      set({
        error: err instanceof Error ? err.message : "Failed to update menu",
      });
    }
  },

  deleteMenu: async (id) => {
    set({ error: null });
    try {
      await api.delete(id);
      await get().fetchTree();
    } catch (err) {
      set({
        error: err instanceof Error ? err.message : "Failed to delete menu",
      });
    }
  },

  moveMenu: async (id, newParentId, orderIndex) => {
    set({ error: null });
    try {
      await api.move(id, { new_parent_id: newParentId, order_index: orderIndex });
      await get().fetchTree();
    } catch (err) {
      set({
        error: err instanceof Error ? err.message : "Failed to move menu",
      });
    }
  },

  reorderMenu: async (id, orderIndex) => {
    set({ error: null });
    try {
      await api.reorder(id, { order_index: orderIndex });
      await get().fetchTree();
    } catch (err) {
      set({
        error: err instanceof Error ? err.message : "Failed to reorder menu",
      });
    }
  },

  setSearchQuery: (query) => set({ searchQuery: query }),

  toggleExpand: (id) => {
    const expanded = [...get().expandedIds];
    const idx = expanded.indexOf(id);
    if (idx !== -1) {
      expanded.splice(idx, 1);
    } else {
      expanded.push(id);
    }
    set({ expandedIds: expanded });
  },

  expandAll: () => {
    set({ expandedIds: get().menus.map((m) => m.id) });
  },

  collapseAll: () => {
    set({ expandedIds: [] });
  },

  setSelectedMenu: (menu) => set({ selectedMenu: menu }),
}));

export { expandedSet };
