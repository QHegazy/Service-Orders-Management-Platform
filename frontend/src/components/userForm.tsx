import Link from "next/link";
import { useState } from "react";

interface UserFormProps {
  setUserType: (userType: "business" | "customer" | null) => void;
}

export const UserForm = ({ setUserType }: UserFormProps) => {
  const [isLoading, setIsLoading] = useState(false);
  const [businessFormData, setBusinessFormData] = useState({
    username: "",
    email: "",
    password: "",
    confirmPassword: "",
  });
  const handleBusinessChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    setBusinessFormData({
      ...businessFormData,
      [e.target.name]: e.target.value,
    });
  };
  const handleBusinessSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (businessFormData.password !== businessFormData.confirmPassword) {
      alert("Passwords do not match");
      return;
    }

    setIsLoading(true);

    try {
      const payload = {
        username: businessFormData.username,
        email: businessFormData.email,
        password: businessFormData.password,
      };
      console.log("Business signup attempt:", payload);
      // Example: await createUser(payload);
    } catch (error) {
      console.error("Business signup error:", error);
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
            Create Business Account
          </h2>
        </div>
        <form className="mt-8 space-y-6" onSubmit={handleBusinessSubmit}>
          <div className="rounded-md shadow-sm -space-y-px">
            <div>
              <input
                id="username"
                name="username"
                type="text"
                required
                className="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                placeholder="Username (3-50 characters)"
                value={businessFormData.username}
                onChange={handleBusinessChange}
                minLength={3}
                maxLength={50}
              />
            </div>
            <div>
              <input
                id="email"
                name="email"
                type="email"
                required
                className="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                placeholder="Email address"
                value={businessFormData.email}
                onChange={handleBusinessChange}
                maxLength={100}
              />
            </div>
            <div>
              <input
                id="password"
                name="password"
                type="password"
                required
                className="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                placeholder="Password (6-100 characters)"
                value={businessFormData.password}
                onChange={handleBusinessChange}
                minLength={6}
                maxLength={100}
              />
            </div>
            <div>
              <input
                id="confirmPassword"
                name="confirmPassword"
                type="password"
                required
                className="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                placeholder="Confirm password"
                value={businessFormData.confirmPassword}
                onChange={handleBusinessChange}
              />
            </div>
          </div>

          <div>
            <button
              type="submit"
              disabled={isLoading}
              className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
            >
              {isLoading
                ? "Creating business account..."
                : "Create Business Account"}
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
        </form>
      </div>
    </div>
  );
};
