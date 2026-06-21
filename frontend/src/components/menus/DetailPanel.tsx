"use client";

import { useState } from "react";
import { useMenuStore } from "@/store/menu-store";
import { Save, Loader2, X } from "lucide-react";

export default function DetailPanel() {
  const { selectedMenu, menus, updateMenu, setSelectedMenu } = useMenuStore();

  const parentMenu = selectedMenu
    ? menus.find((m) => m.id === selectedMenu.parent_id)
    : null;

  if (!selectedMenu) {
    return (
      <div className="flex items-center justify-center h-full text-gray-400 text-sm p-6">
        Select a menu item to view details
      </div>
    );
  }

  return (
    <div className="p-5 space-y-5" key={selectedMenu.id}>
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-semibold text-gray-900 uppercase tracking-wide">Details</h3>
        <button
          onClick={() => setSelectedMenu(null)}
          className="p-1 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-100"
        >
          <X className="w-4 h-4" />
        </button>
      </div>

      <DetailForm
        selectedMenu={selectedMenu}
        parentName={parentMenu?.name ?? null}
        updateMenu={updateMenu}
        setSelectedMenu={setSelectedMenu}
      />
    </div>
  );
}

function DetailForm({
  selectedMenu,
  parentName,
  updateMenu,
  setSelectedMenu,
}: {
  selectedMenu: NonNullable<ReturnType<typeof useMenuStore.getState>["selectedMenu"]>;
  parentName: string | null;
  updateMenu: ReturnType<typeof useMenuStore.getState>["updateMenu"];
  setSelectedMenu: ReturnType<typeof useMenuStore.getState>["setSelectedMenu"];
}) {
  const [name, setName] = useState(selectedMenu.name);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

  const handleSave = async () => {
    if (!name.trim()) {
      setError("Name is required");
      return;
    }
    setSaving(true);
    setError("");
    try {
      await updateMenu(selectedMenu.id, name.trim());
      setSelectedMenu({ ...selectedMenu, name: name.trim() });
    } catch {
      setError("Failed to update menu");
    } finally {
      setSaving(false);
    }
  };

  return (
    <>
      <div className="space-y-4">
        <div>
          <label className="block text-xs font-medium text-gray-500 mb-1">Menu ID</label>
          <div className="px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-xs text-gray-500 font-mono truncate">
            {selectedMenu.id}
          </div>
        </div>

        <div>
          <label className="block text-xs font-medium text-gray-500 mb-1">Depth</label>
          <div className="px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm text-gray-700">
            {selectedMenu.depth}
          </div>
        </div>

        <div>
          <label className="block text-xs font-medium text-gray-500 mb-1">Parent Data</label>
          <div className="px-3 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm text-gray-700">
            {parentName || "(Root)"}
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Name</label>
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="w-full px-3 py-2 bg-gray-50 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
            maxLength={255}
          />
          {error && <p className="text-red-500 text-xs mt-1">{error}</p>}
        </div>
      </div>

      <button
        onClick={handleSave}
        disabled={saving}
        className="flex items-center justify-center gap-2 w-full px-4 py-2.5 text-sm text-white bg-[#0051af] rounded-full hover:bg-[#00408e] transition-colors disabled:opacity-60"
      >
        {saving ? (
          <Loader2 className="w-4 h-4 animate-spin" />
        ) : (
          <Save className="w-4 h-4" />
        )}
        Save
      </button>
    </>
  );
}
