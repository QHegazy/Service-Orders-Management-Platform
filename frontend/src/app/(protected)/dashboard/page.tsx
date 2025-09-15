"use client";

import { useEffect } from "react";
import { useSelector, useDispatch } from "react-redux";
import { RootState } from "@/store";
import { loadFromStorage } from "@/store/authSlice";
import { useRoleGuard } from "@/hooks/useRoleGuard";
import {
  ChartBarIcon,
  UsersIcon,
  DocumentTextIcon,
  TicketIcon,
} from "@heroicons/react/24/outline";
import { CustomerLayout } from "@/components/layout/customerLayout";
import { AdminLayout } from "@/components/layout/adminLayout";
import { TechnicianLayout } from "@/components/layout/technicianLayout";

export default function Dashboard() {
  const adminAccess = useRoleGuard("Admin");
  const technicianAccess = useRoleGuard("Technician");
  const customerAccess = useRoleGuard("Customer");
  const dispatch = useDispatch();

  useEffect(() => {
    dispatch(loadFromStorage());
  }, [dispatch]);



  return (
    <div className="space-y-6">
      {adminAccess && <AdminLayout />}
      {technicianAccess && <TechnicianLayout />}
      {customerAccess && <CustomerLayout />}
      
    </div>
  );
}
