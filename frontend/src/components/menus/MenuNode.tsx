"use client";

import { useState } from "react";
import { MenuItem } from "@/lib/types";
import { useMenuStore } from "@/store/menu-store";
import {
  ChevronRight,
  ChevronDown,
  Plus,
  Pencil,
  Trash2,
  GripVertical,
} from "lucide-react";
import MenuForm from "./MenuForm";
import DeleteConfirm from "./DeleteConfirm";

interface MenuNodeProps {
  item: MenuItem;
  depth?: number;
}

export default function MenuNode({ item, depth = 0 }: MenuNodeProps) {
  const {
    expandedIds,
    toggleExpand,
    searchQuery,
    setSelectedMenu,
    selectedMenu,
    moveMenu,
  } = useMenuStore();
  const [showAddForm, setShowAddForm] = useState(false);
  const [showEditForm, setShowEditForm] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [isDragOver, setIsDragOver] = useState(false);
  const [isDragging, setIsDragging] = useState(false);

  const isExpanded = expandedIds.has(item.id);
  const hasChildren = item.children && item.children.length > 0;
  const isSelected = selectedMenu?.id === item.id;

  const matchesSearch = searchQuery
    ? item.name.toLowerCase().includes(searchQuery.toLowerCase())
    : false;

  const childMatchesSearch = searchQuery
    ? item.children?.some(
        (c) =>
          c.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
          c.children?.some((gc) =>
            gc.name.toLowerCase().includes(searchQuery.toLowerCase()),
          ),
      )
    : false;

  const shouldShow = !searchQuery || matchesSearch || childMatchesSearch;

  if (!shouldShow) return null;

  const handleNodeClick = () => {
    setSelectedMenu(item);
    if (hasChildren) {
      toggleExpand(item.id);
    }
  };

  const handleDragStart = (e: React.DragEvent<HTMLDivElement>) => {
    e.stopPropagation();
    e.dataTransfer.setData("text/plain", item.id);
    e.dataTransfer.effectAllowed = "move";
    setIsDragging(true);
  };

  const handleDragEnd = (e: React.DragEvent<HTMLDivElement>) => {
    e.stopPropagation();
    setIsDragging(false);
  };

  const handleDragOver = (e: React.DragEvent<HTMLDivElement>) => {
    e.stopPropagation();
    e.preventDefault();
    e.dataTransfer.dropEffect = "move";
    setIsDragOver(true);
  };

  const handleDragLeave = (e: React.DragEvent<HTMLDivElement>) => {
    e.stopPropagation();
    setIsDragOver(false);
  };

  const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
    e.stopPropagation();
    e.preventDefault();
    setIsDragOver(false);
    const draggedId = e.dataTransfer.getData("text/plain");
    if (!draggedId || draggedId === item.id) return;
    const tree = useMenuStore.getState().tree;
    if (isOrContains(tree, draggedId, item.id)) return;
    moveMenu(draggedId, item.id, 0);
  };

  const lineLeft = depth * 20 + 17;

  return (
    <>
      <div
        className={`group flex items-center gap-1 py-1.5 px-2 rounded-lg cursor-pointer transition-colors relative ${
          isDragOver ? "bg-blue-50 ring-2 ring-blue-400" : isSelected
          // ? "bg-blue-50 text-blue-700"
          // : "hover:bg-gray-50 hover:border"
        } ${isDragging ? "opacity-50" : ""}`}
        style={{ paddingLeft: `${depth * 20 + 8}px` }}
        onClick={handleNodeClick}
        draggable={true}
        onDragStart={handleDragStart}
        onDragEnd={handleDragEnd}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
      >
        {/*<GripVertical className="w-4 h-4 text-[#101828] opacity-0 group-hover:opacity-100 transition-opacity shrink-0" />*/}

        <button
          onClick={(e) => {
            e.stopPropagation();
            toggleExpand(item.id);
          }}
          className="w-5 h-5 flex items-center justify-center text-[#101828] shrink-0"
        >
          {hasChildren ? (
            isExpanded ? (
              <ChevronDown className="w-4 h-4" />
            ) : (
              <ChevronRight className="w-4 h-4" />
            )
          ) : (
            <span className="w-5 h-5" />
          )}
        </button>

        <div className="flex items-center">
          <span className="flex-1 text-sm truncate select-none text-[#333333]">
            {item.name}
          </span>
        </div>

        <div className="flex items-center mx-2">
          {isSelected && (
            <button
              onClick={(e) => {
                e.stopPropagation();
                setShowAddForm(true);
              }}
              className="flex items-center justify-center w-5 h-5 rounded-full bg-[#0051af] text-white hover:bg-[#00408e] transition-colors shrink-0"
              title="Add child"
            >
              <Plus className="w-3 h-3" />
            </button>
          )}
        </div>

        <div className="ml-auto hidden group-hover:flex items-center gap-0.5">
          {/*<button
            onClick={(e) => {
              e.stopPropagation();
              setShowAddForm(true);
            }}
            className="flex items-center justify-center w-5 h-5 rounded-full bg-[#0051af] text-white hover:bg-[#00408e] transition-colors shrink-0"
            title="Add child"
          >
            <Plus className="w-3 h-3" />
          </button>*/}
          {!isSelected && (
            <button
              onClick={(e) => {
                e.stopPropagation();
                setShowAddForm(true);
              }}
              className="p-1 text-gray-400 hover:text-blue-600 rounded hover:bg-blue-50"
              title="Add child"
            >
              <Plus className="w-3.5 h-3.5" />
            </button>
          )}
          <button
            onClick={(e) => {
              e.stopPropagation();
              setShowEditForm(true);
            }}
            className="p-1 text-gray-400 hover:text-green-600 rounded hover:bg-green-50"
            title="Edit"
          >
            <Pencil className="w-3.5 h-3.5" />
          </button>
          <button
            onClick={(e) => {
              e.stopPropagation();
              setShowDelete(true);
            }}
            className="p-1 text-gray-400 hover:text-red-600 rounded hover:bg-red-50"
            title="Delete"
          >
            <Trash2 className="w-3.5 h-3.5" />
          </button>
          <GripVertical className="w-4 h-4 text-[#101828] opacity-0 group-hover:opacity-100 transition-opacity shrink-0" />
        </div>
      </div>

      {isExpanded && hasChildren && (
        <div className="relative">
          {item.children
            .sort((a, b) => a.order_index - b.order_index)
            .map((child, index) => {
              const isLast = index === item.children.length - 1;
              return (
                <div key={child.id} className="relative">
                  <div
                    className="absolute top-0 w-px bg-[#98A2B3]"
                    style={{
                      left: `${lineLeft}px`,
                      height: isLast ? "16px" : "100%",
                    }}
                  />
                  <div
                    className="absolute h-px bg-[#98A2B3]"
                    style={{
                      left: `${lineLeft}px`,
                      top: "16px",
                      width: "10px",
                    }}
                  />
                  <MenuNode item={child} depth={depth + 1} />
                </div>
              );
            })}
        </div>
      )}

      {showAddForm && (
        <MenuForm
          parentId={item.id}
          parentName={item.name}
          depth={item.depth + 1}
          onClose={() => setShowAddForm(false)}
        />
      )}

      {showEditForm && (
        <MenuForm editItem={item} onClose={() => setShowEditForm(false)} />
      )}

      {showDelete && (
        <DeleteConfirm item={item} onClose={() => setShowDelete(false)} />
      )}
    </>
  );
}

function isOrContains(
  items: MenuItem[],
  ancestorId: string,
  targetId: string,
): boolean {
  for (const item of items) {
    if (item.id === ancestorId) {
      return containsTarget(item, targetId);
    }
    if (item.children) {
      const result = isOrContains(item.children, ancestorId, targetId);
      if (result) return true;
    }
  }
  return false;
}

function containsTarget(item: MenuItem, targetId: string): boolean {
  if (item.id === targetId) return true;
  if (!item.children) return false;
  for (const child of item.children) {
    if (containsTarget(child, targetId)) return true;
  }
  return false;
}
