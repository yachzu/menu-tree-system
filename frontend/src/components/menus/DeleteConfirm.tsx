"use client";

import { useState } from "react";
import { MenuItem } from "@/lib/types";
import { useMenuStore } from "@/store/menu-store";
import { AlertTriangle, Loader2 } from "lucide-react";

interface DeleteConfirmProps {
  item: MenuItem;
  onClose: () => void;
}

export default function DeleteConfirm({ item, onClose }: DeleteConfirmProps) {
  const [isDeleting, setIsDeleting] = useState(false);
  const { deleteMenu } = useMenuStore();
  const hasChildren = item.children && item.children.length > 0;

  const handleDelete = async () => {
    setIsDeleting(true);
    await deleteMenu(item.id);
    onClose();
  };

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50" onClick={onClose}>
      <div
        className="bg-white rounded-xl shadow-xl w-full max-w-sm mx-4 p-6"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-start gap-4 mb-4">
          <div className="w-10 h-10 rounded-full bg-red-100 flex items-center justify-center shrink-0">
            <AlertTriangle className="w-5 h-5 text-red-500" />
          </div>
          <div>
            <h3 className="text-lg font-semibold text-gray-900">Delete Menu</h3>
            <p className="text-sm text-gray-500 mt-0.5">
              {hasChildren
                ? "This will also delete all sub-menus"
                : "This action cannot be undone"}
            </p>
          </div>
        </div>

        <p className="text-sm text-gray-700 mb-6">
          Are you sure you want to delete <strong>{item.name}</strong>
          {hasChildren && " and all its children"}?
        </p>

        <div className="flex justify-end gap-3">
          <button
            onClick={onClose}
            disabled={isDeleting}
            className="px-4 py-2 text-sm text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors disabled:opacity-50"
          >
            Cancel
          </button>
          <button
            onClick={handleDelete}
            disabled={isDeleting}
            className="flex items-center gap-2 px-4 py-2 text-sm text-white bg-red-500 rounded-lg hover:bg-red-600 transition-colors disabled:opacity-60"
          >
            {isDeleting && <Loader2 className="w-4 h-4 animate-spin" />}
            Delete
          </button>
        </div>
      </div>
    </div>
  );
}
