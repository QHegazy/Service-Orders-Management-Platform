"use client";

import Link from "next/link";
import { useState } from "react";
import { useRouter } from "next/navigation";
import customerService from "@/services/custmoerService";
import { ErrorMessage, SuccessMessage } from "./messageForm";

interface CustomerFormProps {
  setUserType: (userType: "business" | "customer" | null) => void;
}

export const CustomerForm = ({ setUserType }: CustomerFormProps) => {
  const router = useRouter();
  const [customerFormData, setCustomerFormData] = useState({
    firstName: "",
    lastName: "",
    username: "",
    email: "",
    password: "",
  });
  const [isLoading, setIsLoading] = useState(false);
  const [message, setMessage] = useState<{
    text: string;
    type: "error" | "success";
  } | null>(null);

  const handleCustomerChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    setCustomerFormData({
      ...customerFormData,
      [e.target.name]: e.target.value,
    });
  };

  const handleCustomerSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setMessage(null);
    setIsLoading(true);

    try {
      const payload = {
        first_name: customerFormData.firstName,
        last_name: customerFormData.lastName,
        username: customerFormData.username,
        email: customerFormData.email || undefined,
        password: customerFormData.password,
      };

      console.log("Attempting customer signup with payload:", payload);
      const response = await customerService.signup(payload);
      console.log("Customer signup success:", response);

      setMessage({
        text: "Customer account created successfully! Redirecting to login...",
        type: "success",
      });

      setTimeout(() => {
        console.log("Redirecting to login...");
        router.push("/login");
      }, 1000);
    } catch (error: any) {
      console.error("Customer signup error:", error);
      setMessage({
        text:
          error?.message ||
          "Failed to create customer account. Please try again.",
        type: "error",
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full space-y-8">
        <div>
          <button
            onClick={() => setUserType(null)}
            className="text-indigo-600 hover:text-indigo-500 text-sm"
          >
            ← Back to account type selection
          </button>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Create Customer Account
          </h2>
        </div>

        {message && (
          <div className="mt-4">
            {message.type === "error" ? (
              <ErrorMessage message={message.text} />
            ) : (
              <SuccessMessage message={message.text} />
            )}
          </div>
        )}

        <form className="mt-8 space-y-6" onSubmit={handleCustomerSubmit}>
          <div className="rounded-md shadow-sm -space-y-px">
            <input
              id="firstName"
              name="firstName"
              type="text"
              required
              placeholder="First name"
              value={customerFormData.firstName}
              onChange={handleCustomerChange}
              className="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              maxLength={100}
            />
            <input
              id="lastName"
              name="lastName"
              type="text"
              required
              placeholder="Last name"
              value={customerFormData.lastName}
              onChange={handleCustomerChange}
              className="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              maxLength={100}
            />
            <input
              id="username"
              name="username"
              type="text"
              required
              placeholder="Username (3-50 characters)"
              value={customerFormData.username}
              onChange={handleCustomerChange}
              className="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              minLength={3}
              maxLength={50}
            />
            <input
              id="email"
              name="email"
              type="email"
              required
              placeholder="Email"
              value={customerFormData.email}
              onChange={handleCustomerChange}
              className="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              maxLength={100}
            />
            <input
              id="password"
              name="password"
              type="password"
              required
              placeholder="Password (6-100 characters)"
              value={customerFormData.password}
              onChange={handleCustomerChange}
              className="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              minLength={6}
              maxLength={100}
            />
          </div>

          <button
            type="submit"
            disabled={isLoading}
            className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
          >
            {isLoading ? "Creating customer account..." : "Create Customer Account"}
          </button>

          <div className="text-center">
            <Link href="/login" className="text-indigo-600 hover:text-indigo-500">
              Already have an account? Sign in
            </Link>
          </div>
        </form>
      </div>
    </div>
  );
};