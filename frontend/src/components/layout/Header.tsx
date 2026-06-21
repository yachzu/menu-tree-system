"use client";

import { useMemo } from "react";
import { Menu, Folder } from "lucide-react";
import { useMenuStore } from "@/store/menu-store";

interface HeaderProps {
  onMenuToggle: () => void;
}

export default function Header({ onMenuToggle }: HeaderProps) {
  const { selectedMenu, menus } = useMenuStore();

  const breadcrumbs = useMemo(() => {
    if (!selectedMenu) return [];
    const path: { id: string; name: string }[] = [];
    let current = selectedMenu;
    while (current) {
      path.unshift({ id: current.id, name: current.name });
      if (!current.parent_id) break;
      const parent = menus.find((m) => m.id === current.parent_id);
      if (!parent) break;
      current = parent;
    }
    return path;
  }, [selectedMenu, menus]);

  return (
    <header
      className="flex items-center justify-between bg-white px-4 md:px-8 lg:px-12"
      style={{ height: "84px" }}
    >
      <div className="flex items-center gap-2 min-w-0">
        <button
          onClick={onMenuToggle}
          className="p-2 -ml-2 shrink-0 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg lg:hidden"
        >
          <Menu className="w-5 h-5" />
        </button>

        <div className="flex items-center gap-1 text-sm overflow-x-auto whitespace-nowrap scrollbar-thin">
          <Folder className="w-4 h-4 text-gray-400" fill="#9CA3AF" />
          {breadcrumbs.length === 0 ? (
              <>
                <span className="text-gray-400">/</span>
                <span className="text-gray-700 font-medium">Menus</span>
              </>
            ) : (
              breadcrumbs.map((crumb, i) => (
                <div key={crumb.id} className="flex items-center gap-1">
                  {i > 0 && <span className="text-gray-400">/</span>}
                  <span
                    className={
                      i === breadcrumbs.length - 1
                        ? "text-gray-700 font-medium"
                        : "text-gray-400"
                    }
                  >
                    {crumb.name}
                  </span>
                </div>
              ))
            )}
        </div>
      </div>
    </header>
  );
}
