"use client";

import logoSTK from "@/assets/logo-stk.png";
import logoSTKIcon from "@/assets/logo-stk-icon.png";
import menuOpen from "@/assets/menu_open.png";
import Image from "next/image";
import { Folder, LayoutGrid, ListTree, X } from "lucide-react";
import { useMenuStore } from "@/store/menu-store";
import { MenuItem } from "@/lib/types";

interface SidebarProps {
  isOpen: boolean;
  onClose: () => void;
  minimized?: boolean;
  onToggleMinimize?: () => void;
}

export default function Sidebar({
  isOpen,
  onClose,
  minimized = false,
  onToggleMinimize,
}: SidebarProps) {
  const { selectedMenu, menus, setSelectedMenu } = useMenuStore();

  const depth2Menus = menus.filter((m) => m.depth === 2);

  return (
    <>
      {isOpen && (
        <div
          className="fixed inset-0 bg-black/30 z-20 lg:hidden"
          onClick={onClose}
        />
      )}

      <aside
        className={`fixed top-6 left-6 z-30 bg-primary flex flex-col transition-all duration-200 rounded-[24px] ${
          minimized ? "w-20" : "w-60"
        } ${isOpen ? "translate-x-0" : "-translate-x-[calc(100%+24px)] lg:translate-x-0"}`}
        style={{ height: "calc(100vh - 48px)" }}
      >
        <div
          className={`flex border-b border-white/10 ${
            minimized
              ? "flex-col items-center px-3 py-4 gap-3"
              : "items-center justify-between px-8 py-[30px]"
          }`}
          style={minimized ? {} : { height: "89.67px" }}
        >
          {!minimized ? (
            <>
              <Image
                src={logoSTK}
                alt="Logo"
                style={{ width: "70px", height: "29.67px" }}
              />
              <button
                onClick={onClose}
                className="flex lg:hidden items-center justify-center text-white/60 hover:text-white transition-colors"
                style={{ width: "24px", height: "24px" }}
              >
                <X className="w-5 h-5" />
              </button>
              <button
                onClick={onToggleMinimize}
                className="hidden lg:flex items-center justify-center text-white/60 hover:text-white transition-colors"
                style={{ width: "24px", height: "24px" }}
              >
                <Image src={menuOpen} alt="Toggle" width={24} height={24} />
              </button>
            </>
          ) : (
            <div className="flex flex-col items-center gap-3 w-full">
              <Image
                src={logoSTKIcon}
                alt="Logo"
                className="mx-auto mt-2.5"
                style={{ width: "28px", height: "28px", objectFit: "contain" }}
              />
            </div>
          )}
        </div>

        <nav className="flex-1 overflow-y-auto py-4 px-3 space-y-4 scrollbar-thin">
          {depth2Menus.length > 0 ? (
            depth2Menus.map((item) => (
              <SidebarGroup
                key={item.id}
                group={item}
                allMenus={menus}
                selectedId={selectedMenu?.id ?? null}
                onSelect={setSelectedMenu}
                minimized={minimized}
              />
            ))
          ) : (
            <div className="text-center text-white/40 text-xs py-8">
              <ListTree className="w-6 h-6 mx-auto mb-2 opacity-50" />
              {!minimized && <p>No menus yet</p>}
            </div>
          )}
        </nav>

        {!minimized && (
          <div className="px-5 py-4 border-t border-white/10">
            <p className="text-white/40 text-xs">STK Menu Tree System</p>
            <p className="text-white/30 text-[10px]">v1.0.0</p>
          </div>
        )}

        {minimized && (
          <div className="px-auto py-5 border-t border-white/10">
            <button
              onClick={onToggleMinimize}
              className="hidden lg:flex items-center justify-center w-6 h-6 mx-auto text-white/60 hover:text-white transition-colors rotate-180"
            >
              <Image src={menuOpen} alt="Toggle" width={22} height={22} />
            </button>
          </div>
        )}
      </aside>
    </>
  );
}

function findDepth3Children(items: MenuItem[], depth2Id: string): MenuItem[] {
  return items.filter((m) => m.parent_id === depth2Id && m.depth === 3);
}

function SidebarGroup({
  group,
  allMenus,
  selectedId,
  onSelect,
  minimized,
}: {
  group: MenuItem;
  allMenus: MenuItem[];
  selectedId: string | null;
  onSelect: (menu: MenuItem | null) => void;
  minimized: boolean;
}) {
  const children = findDepth3Children(allMenus, group.id);
  const isSelected = selectedId === group.id;

  return (
    <div className="bg-white/10 rounded-xl p-2 space-y-0.5">
      <button
        onClick={() => onSelect(group)}
        className={`flex items-center gap-3 px-3 py-2.5 text-sm w-full text-left font-medium rounded-lg transition-colors ${
          minimized ? "justify-center" : ""
        } ${isSelected ? "bg-white text-black" : "text-white hover:bg-white/10"}`}
        title={minimized ? group.name : undefined}
      >
        <Folder
          className={`w-4 h-4 shrink-0 ${isSelected ? "text-[#0051af]" : "text-white"}`}
          fill={isSelected ? "#0051af" : "none"}
        />
        {!minimized && <span className="truncate">{group.name}</span>}
      </button>
      {children.length > 0 && !minimized && (
        <div className="space-y-0.5">
          {children.map((child) => {
            const isChildSelected = selectedId === child.id;
            return (
              <button
                key={child.id}
                onClick={() => onSelect(child)}
                className={`flex items-center gap-3 px-3 py-2 rounded-lg text-sm w-full text-left transition-colors ${
                  isChildSelected
                    ? "bg-white text-black font-medium"
                    : "text-white hover:bg-white/10"
                }`}
              >
                <LayoutGrid
                  className={`w-4 h-4 shrink-0 ${
                    isChildSelected ? "text-[#0051af]" : "text-white"
                  }`}
                  fill={isChildSelected ? "#0051af" : "none"}
                />
                <span className="truncate">{child.name}</span>
              </button>
            );
          })}
        </div>
      )}
    </div>
  );
}
