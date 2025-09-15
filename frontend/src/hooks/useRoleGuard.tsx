"use client";

import { useSelector } from "react-redux";
import { RootState } from "@/store";

export const useRoleGuard = (role: string | string[]): boolean => {
  const userRole = useSelector((state: RootState) => state.auth.user?.role);

  if (!userRole) return false;

  if (Array.isArray(role)) {
    return role.includes(userRole);
  }

  return userRole === role;
};