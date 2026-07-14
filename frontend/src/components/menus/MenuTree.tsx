"use client";

import { useEffect } from "react";
import { useMenuStore } from "@/store/menu-store";
import { MenuItem } from "@/lib/types";
import MenuNode from "./MenuNode";
import SearchBar from "./SearchBar";
import { ChevronsDownUp, ChevronsUpDown, Loader2, AlertCircle, FolderOpen } from "lucide-react";

export default function MenuTree() {
  const {
    tree,
    searchQuery,
    isLoading,
    error,
    fetchTree,
    expandAll,
    collapseAll,
  } = useMenuStore();

  useEffect(() => {
    fetchTree();
  }, [fetchTree]);

  const filteredTree = searchQuery
    ? tree
        .map((item) => filterTree(item, searchQuery))
        .filter((item): item is MenuItem => item !== null)
    : tree;

  if (error) {
    return (
      <div className="flex items-center justify-center p-12">
        <div className="text-center max-w-sm">
          <AlertCircle className="w-12 h-12 text-red-400 mx-auto mb-4" />
          <p className="text-red-600 font-medium mb-2">Failed to load menu tree</p>
          <p className="text-sm text-gray-500 mb-4">{error}</p>
          <button
            onClick={() => fetchTree()}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-sm font-medium"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center gap-3">
        <div className="flex-1">
          <SearchBar />
        </div>
        <button
          onClick={expandAll}
          className="flex items-center gap-1.5 px-4 py-2 text-sm text-white bg-gray-900 rounded-full hover:bg-gray-800 transition-colors"
          title="Expand All"
        >
          <ChevronsDownUp className="w-4 h-4" />
          <span className="hidden sm:inline">Expand All</span>
        </button>
        <button
          onClick={collapseAll}
          className="flex items-center gap-1.5 px-4 py-2 text-sm text-gray-700 bg-white border border-gray-200 rounded-full hover:bg-gray-50 transition-colors"
          title="Collapse All"
        >
          <ChevronsUpDown className="w-4 h-4" />
          <span className="hidden sm:inline">Collapse All</span>
        </button>
      </div>

      <div className="py-1">
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <Loader2 className="w-8 h-8 text-blue-500 animate-spin" />
          </div>
        ) : filteredTree.length === 0 ? (
          <div className="text-center py-12">
            <FolderOpen className="w-12 h-12 text-gray-300 mx-auto mb-3" />
            <p className="text-gray-500">
              {searchQuery
                ? "No menus match your search"
                : "Belum ada menu. Tambahkan menu pertama Anda."}
            </p>
          </div>
        ) : (
          <div className="space-y-0.5">
            {filteredTree
              .sort((a, b) => a.name.localeCompare(b.name))
              .map((item) => (
                <MenuNode key={item.id} item={item} />
              ))}
          </div>
        )}
      </div>
    </div>
  );
}

function filterTree(item: MenuItem, query: string): MenuItem | null {
  const q = query.toLowerCase();
  const nameMatch = item.name.toLowerCase().includes(q);
  const filteredChildren = item.children
    ? (item.children as MenuItem[])
        .map((child) => filterTree(child, q))
        .filter((child): child is MenuItem => child !== null)
    : [];

  if (nameMatch || filteredChildren.length > 0) {
    return { ...item, children: filteredChildren };
  }
  return null;
}
