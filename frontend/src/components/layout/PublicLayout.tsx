"use client";

import { TopNav } from "./TopNav";
import { Footer } from "./Footer";

interface PublicLayoutProps {
  children: React.ReactNode;
}

export function PublicLayout({ children }: PublicLayoutProps) {
  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">
      <TopNav onMenuClick={() => {}} />
      <main className="pt-16 flex-1">{children}</main>
      <Footer />
    </div>
  );
}
