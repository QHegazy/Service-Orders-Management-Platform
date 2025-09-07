import Link from "next/link";
import { useState } from "react";

interface CustomerFormProps {
  setUserType: (userType: "business" | "customer" | null) => void;
}

export const CustomerForm = ({ setUserType }: CustomerFormProps) => {
  const [customerFormData, setCustomerFormData] = useState({
    firstName: "",
    lastName: "",
    username: "",
    email: "",
  });
  const [isLoading, setIsLoading] = useState(false);

  const handleCustomerChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setCustomerFormData({
      ...customerFormData,
      [e.target.name]: e.target.value,
    });
  };

  const handleCustomerSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const payload = {
        first_name: customerFormData.firstName,
        last_name: customerFormData.lastName,
        username: customerFormData.username,
        email: customerFormData.email,
      };
      console.log("Customer signup attempt:", payload);
      // Example: await createCustomer(payload);
    } catch (error) {
      console.error("Customer signup error:", error);
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
        <form className="mt-8 space-y-6" onSubmit={handleCustomerSubmit}>
          <div className="rounded-md shadow-sm -space-y-px">
            <div>
              <input
                id="firstName"
                name="firstName"
                type="text"
                required
                className="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                placeholder="First name"
                value={customerFormData.firstName}
                onChange={handleCustomerChange}
                maxLength={100}
              />
            </div>
            <div>
              <input
                id="lastName"
                name="lastName"
                type="text"
                required
                className="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                placeholder="Last name"
                value={customerFormData.lastName}
                onChange={handleCustomerChange}
                maxLength={100}
              />
            </div>
            <div>
              <input
                id="username"
                name="username"
                type="text"
                required
                className="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                placeholder="Username (3-50 characters)"
                value={customerFormData.username}
                onChange={handleCustomerChange}
                minLength={3}
                maxLength={50}
              />
            </div>
            <div>
              <input
                id="email"
                name="email"
                type="email"
                className="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                placeholder="Email address (optional)"
                value={customerFormData.email}
                onChange={handleCustomerChange}
                maxLength={100}
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
                ? "Creating customer account..."
                : "Create Customer Account"}
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
