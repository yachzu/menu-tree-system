"use client";

import { useState } from "react";
import { MenuItem } from "@/lib/types";
import { useMenuStore } from "@/store/menu-store";
import { Save, X, Loader2 } from "lucide-react";

interface MenuFormProps {
  parentId?: string | null;
  parentName?: string;
  depth?: number;
  editItem?: MenuItem;
  onClose: () => void;
}

export default function MenuForm({ parentId, parentName, depth, editItem, onClose }: MenuFormProps) {
  const isEdit = !!editItem;
  const [name, setName] = useState(editItem?.name || "");
  const [error, setError] = useState("");
  const [saving, setSaving] = useState(false);
  const { createMenu, updateMenu } = useMenuStore();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      setError("Name is required");
      return;
    }
    setError("");
    setSaving(true);

    try {
      if (isEdit && editItem) {
        await updateMenu(editItem.id, name.trim());
      } else {
        await createMenu(name.trim(), parentId ?? null);
      }
      onClose();
    } catch {
      setError("Failed to save menu");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50" onClick={onClose}>
      <div
        className="bg-white rounded-xl shadow-xl w-full max-w-md mx-4 p-6"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between mb-5">
          <h3 className="text-lg font-semibold text-gray-900">
            {isEdit ? "Edit Menu Item" : "Add Menu Item"}
          </h3>
          <button
            onClick={onClose}
            className="p-1 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-100"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          {!isEdit && parentName !== undefined && (
            <div>
              <label className="block text-xs font-medium text-gray-500 mb-1">Parent</label>
              <div className="px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm text-gray-700">
                {parentName || "(Root)"}
              </div>
            </div>
          )}

          {!isEdit && depth !== undefined && (
            <div>
              <label className="block text-xs font-medium text-gray-500 mb-1">Depth</label>
              <div className="px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm text-gray-700">
                {depth}
              </div>
            </div>
          )}

          {isEdit && editItem && (
            <>
              <div>
                <label className="block text-xs font-medium text-gray-500 mb-1">Menu ID</label>
                <div className="px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm text-gray-500 font-mono truncate">
                  {editItem.id}
                </div>
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-500 mb-1">Depth</label>
                <div className="px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm text-gray-700">
                  {editItem.depth}
                </div>
              </div>
              {editItem.parent_id && (
                <div>
                  <label className="block text-xs font-medium text-gray-500 mb-1">Parent ID</label>
                  <div className="px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm text-gray-500 font-mono truncate">
                    {editItem.parent_id}
                  </div>
                </div>
              )}
            </>
          )}

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Name</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Enter menu name"
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
              autoFocus
              maxLength={255}
            />
            {error && <p className="text-red-500 text-xs mt-1">{error}</p>}
          </div>

          <div className="flex justify-end gap-3 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-sm text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
              disabled={saving}
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={saving}
              className="flex items-center gap-2 px-4 py-2 text-sm text-white bg-[#0051af] rounded-lg hover:bg-[#00408e] transition-colors disabled:opacity-60"
            >
              {saving ? (
                <Loader2 className="w-4 h-4 animate-spin" />
              ) : (
                <Save className="w-4 h-4" />
              )}
              {isEdit ? "Save" : "Add"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
