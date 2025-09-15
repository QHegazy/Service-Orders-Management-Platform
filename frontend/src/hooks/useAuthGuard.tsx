"use client";

import { useEffect, useState } from "react";
import { useRouter, usePathname } from "next/navigation";
import { useSelector } from "react-redux";
import { RootState } from "@/store";

export const useAuthGuard = () => {
  const router = useRouter();
  const pathname = usePathname();
  const isAuthenticated = useSelector(
    (state: RootState) => state.auth.isAuthenticated
  );

  const [hydrated, setHydrated] = useState(false);

  useEffect(() => {
    setHydrated(true);
  }, []);

  useEffect(() => {
    if (!hydrated) return;

    const protectedRoutes = [
      "/dashboard",
      "/technicians",
      "/tenants",
      "/tickets",
      "/profile",
    ];
    if (
      !isAuthenticated &&
      protectedRoutes.some((route) => pathname.startsWith(route))
    ) {
      router.replace("/login");
    }

    const publicAuthRoutes = ["/login", "/signup", "/"];
    if (isAuthenticated && publicAuthRoutes.includes(pathname)) {
      router.replace("/dashboard");
    }
  }, [hydrated, isAuthenticated, pathname, router]);
};
