"use client";

import { useState } from "react";
import Link from "next/link";
import { UserForm } from "@/components/userForm";
import { CustomerForm } from "@/components/customerForm";

export default function SignupPage() {
  const [userType, setUserType] = useState<"business" | "customer" | null>(
    null
  );

  if (!userType) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="max-w-md w-full space-y-8">
          <div>
            <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
              Choose Account Type
            </h2>
            <p className="mt-2 text-center text-sm text-gray-600">
              Select how you want to use our platform
            </p>
          </div>
          <div className="space-y-4">
            <button
              onClick={() => setUserType("business")}
              className="w-full flex flex-col items-center p-6 border-2 border-gray-300 rounded-lg hover:border-indigo-500 hover:bg-indigo-50 transition-colors"
            >
              <div className="text-2xl mb-2">🏢</div>
              <h3 className="text-lg font-medium text-gray-900">
                Business Owner
              </h3>
              <p className="text-sm text-gray-500 text-center">
                Manage your business operations and technicians
              </p>
            </button>
            <button
              onClick={() => setUserType("customer")}
              className="w-full flex flex-col items-center p-6 border-2 border-gray-300 rounded-lg hover:border-indigo-500 hover:bg-indigo-50 transition-colors"
            >
              <div className="text-2xl mb-2">👤</div>
              <h3 className="text-lg font-medium text-gray-900">Customer</h3>
              <p className="text-sm text-gray-500 text-center">
                Book services and manage your appointments
              </p>
            </button>
          </div>
          <div className="text-center">
            <Link
              href="/login"
              className="text-indigo-600 hover:text-indigo-500"
            >
              Already have an account? Sign in
            </Link>
          </div>
        </div>
      </div>
    );
  }

  if (userType === "business") {
    return <UserForm setUserType={setUserType} />;
  }

  if (userType === "customer") {
    return <CustomerForm setUserType={setUserType} />;
  }

  return null;
}
