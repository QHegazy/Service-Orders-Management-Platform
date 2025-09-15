"use client";

import { useState } from "react";
import { TopNav } from "./TopNav";
import { SideNav } from "./SideNav";

interface AppLayoutProps {
  children: React.ReactNode;
}

export function AppLayout({ children }: AppLayoutProps) {
  const [sidebarOpen, setSidebarOpen] = useState(false);

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">
      <TopNav onMenuClick={() => setSidebarOpen(!sidebarOpen)} />

      <div className="flex  flex-1">
        <SideNav isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} />
        <main
          className={`flex-1 transition-all duration-300 ${
            sidebarOpen ? "lg:ml-64" : "ml-0"
          } flex flex-col`}
        >
          <div className="p-6 flex-1">{children}</div>
        </main>
      </div>
      {sidebarOpen && (
        <div
          className="fixed inset-0 bg-black bg-opacity-50 z-20 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}
    </div>
  );
}
