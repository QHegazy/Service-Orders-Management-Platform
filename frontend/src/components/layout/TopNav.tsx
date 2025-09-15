"use client";

import Link from "next/link";
import { useSelector, useDispatch } from "react-redux";
import { useEffect, useState } from "react";
import { RootState } from "@/store";
import { logout, loadFromStorage } from "@/store/authSlice";
import { useRouter } from "next/navigation";
import {
  Bars3Icon,
  BellIcon,
  ChevronDownIcon,
  UserCircleIcon,
  Cog6ToothIcon,
  ArrowRightOnRectangleIcon,
} from "@heroicons/react/24/outline";

interface TopNavProps {
  onMenuClick: () => void;
}

const buttonClass =
  "p-2 rounded-md text-gray-600 hover:text-gray-900 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500";
const dropdownItemClass =
  "flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100";

export function TopNav({ onMenuClick }: TopNavProps) {
  const { user, isAuthenticated } = useSelector(
    (state: RootState) => state.auth
  );
  const dispatch = useDispatch();
  const router = useRouter();
  const [dropdownOpen, setDropdownOpen] = useState(false);

  useEffect(() => {
    dispatch(loadFromStorage());
  }, [dispatch]);

  const handleLogout = () => {
    dispatch(logout());
    router.push("/");
    setDropdownOpen(false);
  };

  const getUserInitials = (name: string) => {
    const words = name.split(" ");
    return words.length === 1
      ? name.slice(0, 2).toUpperCase()
      : (words[0][0] + words[1][0]).toUpperCase();
  };

  const UserDropdown = () => (
    <div className="relative">
      <button
        onClick={() => setDropdownOpen(!dropdownOpen)}
        className={`flex items-center space-x-3 ${buttonClass}`}
      >
        <div className="w-8 h-8 bg-indigo-600 text-white flex items-center justify-center rounded-full text-sm font-medium">
          {getUserInitials(user!.username)}
        </div>
        <div className="hidden md:block text-left">
          <div className="text-sm font-medium text-gray-900">
            {user!.username}
          </div>
          <div className="text-xs text-gray-500">{user!.role}</div>
        </div>
        <ChevronDownIcon className="h-4 w-4" />
      </button>

      {dropdownOpen && (
        <div className="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg py-1 z-50 border border-gray-200">
          <div className="px-4 py-2 border-b border-gray-100">
            <div className="text-sm font-medium text-gray-900">
              {user!.username}
            </div>
            <div className="text-xs text-gray-500">{user!.role}</div>
          </div>
          <Link
            href="/profile"
            className={dropdownItemClass}
            onClick={() => setDropdownOpen(false)}
          >
            <UserCircleIcon className="h-4 w-4 mr-3" />
            Profile
          </Link>
          <Link
            href="/dashboard"
            className={dropdownItemClass}
            onClick={() => setDropdownOpen(false)}
          >
            <Cog6ToothIcon className="h-4 w-4 mr-3" />
            Dashboard
          </Link>
          <button
            onClick={handleLogout}
            className={`${dropdownItemClass} w-full`}
          >
            <ArrowRightOnRectangleIcon className="h-4 w-4 mr-3" />
            Sign out
          </button>
        </div>
      )}
    </div>
  );

  return (
    <nav className="bg-white shadow-sm border-b border-gray-200 fixed top-0 left-0 right-0 z-30">
      <div className="px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <div className="flex items-center">
            {isAuthenticated && (
              <button
                onClick={onMenuClick}
                className={`${buttonClass} lg:hidden`}
              >
                <Bars3Icon className="h-6 w-6" />
              </button>
            )}
            <div className="flex-shrink-0 flex items-center ml-4 lg:ml-0">
              <Link href="/" className="text-xl font-bold text-gray-900">
                <span className="text-indigo-600">Order</span>Management
              </Link>
            </div>
          </div>

          <div className="flex items-center space-x-4">
            {isAuthenticated && user ? (
              <>
                <button
                  onClick={onMenuClick}
                  className={`hidden lg:block ${buttonClass}`}
                >
                  <Bars3Icon className="h-5 w-5" />
                </button>
                <button className={buttonClass}>
                  <BellIcon className="h-5 w-5" />
                </button>
                <UserDropdown />
              </>
            ) : (
              <div className="flex items-center space-x-4">
                <Link
                  href="/login"
                  className="text-gray-600 hover:text-gray-900 font-medium transition-colors duration-200"
                >
                  Sign in
                </Link>
                <Link
                  href="/signup"
                  className="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700 transition-colors duration-200 text-sm font-medium"
                >
                  Sign up
                </Link>
              </div>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
}
