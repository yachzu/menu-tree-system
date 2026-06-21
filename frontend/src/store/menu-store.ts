"use client";

import { create } from "zustand";
import { MenuItem } from "@/lib/types";
import { api } from "@/lib/api";

interface MenuState {
  menus: MenuItem[];
  tree: MenuItem[];
  searchQuery: string;
  expandedIds: Set<string>;
  isLoading: boolean;
  error: string | null;
  selectedMenu: MenuItem | null;

  fetchTree: () => Promise<void>;
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

export const useMenuStore = create<MenuState>((set, get) => ({
  menus: [],
  tree: [],
  searchQuery: "",
  expandedIds: new Set<string>(),
  isLoading: false,
  error: null,
  selectedMenu: null,

  fetchTree: async () => {
    set({ isLoading: true, error: null });
    try {
      const tree = await api.getTree();
      const menus = flattenTree(tree);
      set({
        tree,
        menus,
        isLoading: false,
        expandedIds: new Set(menus.map((m) => m.id)),
      });
    } catch (err) {
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
    const expanded = new Set(get().expandedIds);
    if (expanded.has(id)) {
      expanded.delete(id);
    } else {
      expanded.add(id);
    }
    set({ expandedIds: expanded });
  },

  expandAll: () => {
    set({ expandedIds: new Set(get().menus.map((m) => m.id)) });
  },

  collapseAll: () => {
    set({ expandedIds: new Set() });
  },

  setSelectedMenu: (menu) => set({ selectedMenu: menu }),
}));
