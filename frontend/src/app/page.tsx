"use client";

import { useState } from "react";
import Sidebar from "@/components/layout/Sidebar";
import Header from "@/components/layout/Header";
import MenuTree from "@/components/menus/MenuTree";
import DetailPanel from "@/components/menus/DetailPanel";
import { ListTree } from "lucide-react";
import { useMenuStore } from "@/store/menu-store";

export default function Home() {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [sidebarMinimized, setSidebarMinimized] = useState(false);
  const { tree, selectedMenu } = useMenuStore();

  return (
    <div className="flex h-screen bg-white">
      <Sidebar
        isOpen={sidebarOpen}
        onClose={() => setSidebarOpen(false)}
        minimized={sidebarMinimized}
        onToggleMinimize={() => setSidebarMinimized(!sidebarMinimized)}
      />

      <div
        className={`flex-1 flex flex-col min-w-0 transition-all duration-200 ${sidebarMinimized ? "lg:pl-[108px]" : "lg:pl-[288px]"} lg:pt-6`}
      >
        <Header onMenuToggle={() => setSidebarOpen(!sidebarOpen)} />

        <div className="flex-1 flex flex-col lg:flex-row overflow-hidden">
          <div className="flex-1 overflow-y-auto">
            <div
              className="flex items-center px-4 md:px-8 lg:px-12"
              style={{ height: "84px" }}
            >
              <div className="flex items-center gap-3 lg:gap-4">
                <span className="inline-flex items-center justify-center rounded-full bg-[#0051af] w-10 h-10 lg:w-[52px] lg:h-[52px]">
                  <svg
                    className="w-[17px] h-[17px] lg:w-[22px] lg:h-[22px]"
                    viewBox="0 0 24 24"
                    fill="none"
                  >
                    <rect
                      x="3"
                      y="3"
                      width="7"
                      height="7"
                      rx="1.5"
                      fill="white"
                    />
                    <rect
                      x="14"
                      y="3"
                      width="7"
                      height="7"
                      rx="1.5"
                      fill="white"
                    />
                    <rect
                      x="3"
                      y="14"
                      width="7"
                      height="7"
                      rx="1.5"
                      fill="white"
                    />
                    <rect
                      x="14"
                      y="14"
                      width="7"
                      height="7"
                      rx="1.5"
                      fill="white"
                    />
                  </svg>
                </span>
                <h1
                  className="text-[#1A1A1A] text-xl lg:text-[32px] truncate"
                  style={{
                    fontWeight: 800,
                    lineHeight: "125%",
                    letterSpacing: "-0.04em",
                  }}
                >
                  {selectedMenu
                    ? selectedMenu.name
                    : tree.find((m) => m.parent_id === null)?.name || "Menus"}
                </h1>
              </div>
            </div>

            <div className="px-4 py-4 md:px-8 lg:px-12">
              <div className="w-full max-w-[349px]">
                <label
                  className="block mb-2 text-[#101828]"
                  style={{
                    fontWeight: 400,
                    fontSize: "14px",
                    lineHeight: "100%",
                    letterSpacing: "-0.02em",
                  }}
                >
                  Menu
                </label>
                <div
                  className="flex items-center bg-[#F9FAFB] w-full"
                  style={{
                    height: "52px",
                    minHeight: "52px",
                    padding: "14px 16px",
                    borderRadius: "16px",
                    gap: "16px",
                  }}
                >
                  <span
                    className="flex-1 text-[#101828]"
                    style={{
                      fontWeight: 400,
                      fontSize: "16px",
                      lineHeight: "100%",
                      letterSpacing: "-0.02em",
                    }}
                  >
                    system management
                  </span>
                  <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                    <path
                      d="M9 13L12 10L15 13"
                      stroke="#475467"
                      strokeWidth="1.5"
                      strokeLinecap="round"
                      strokeLinejoin="round"
                    />
                  </svg>
                </div>
              </div>
            </div>

            <div className="max-w-4xl px-4 md:px-8 lg:px-12">
              <MenuTree />

              {selectedMenu && (
                <div className="mt-6 pt-6 border-t border-gray-200 lg:hidden">
                  <DetailPanel />
                </div>
              )}
            </div>
          </div>

          <div className="hidden lg:block w-[400px] bg-white overflow-y-auto">
            {selectedMenu ? (
              <DetailPanel />
            ) : (
              <div className="flex flex-col items-center justify-center h-full text-gray-400 p-6">
                <ListTree className="w-12 h-12 mb-3 text-gray-300" />
                <p className="text-sm">Select a menu item to edit</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
