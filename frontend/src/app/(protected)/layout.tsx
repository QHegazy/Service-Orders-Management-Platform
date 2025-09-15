"use client";

import { ReactNode } from "react";
import { useAuthGuard } from "@/hooks/useAuthGuard";
import { AppLayout } from "@/components/layout";
import { useAuthCheck } from "@/hooks/useAuthCheck";

export default function ProtectedLayout({ children }: { children: ReactNode }) {
  useAuthGuard();
  useAuthCheck()

  return  <AppLayout>{children}</AppLayout>
  
}