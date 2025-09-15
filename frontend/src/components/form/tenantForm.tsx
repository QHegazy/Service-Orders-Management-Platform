import { useState } from "react";
import { useRouter } from "next/navigation";
import tenantService from "@/services/tenantService";
import {
  CheckCircleIcon,
  ExclamationTriangleIcon,
} from "@heroicons/react/24/outline";

export const TenantForm = () => {
  const [tenantFormData, setTenantFormData] = useState({
    tenantName: "",
    domain: "",
    email: "",
  });
  const [isLoading, setIsLoading] = useState(false);
  const [message, setMessage] = useState<{
    type: "success" | "error";
    text: string;
  } | null>(null);
  const router = useRouter();

  const handleTenantChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setTenantFormData({
      ...tenantFormData,
      [e.target.name]: e.target.value,
    });
  };

  const handleTenantSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setMessage(null);

    try {
      const payload = {
        tenant_name: tenantFormData.tenantName,
        domain: tenantFormData.domain,
        email: tenantFormData.email,
      };

      const response = await tenantService.createTenant(payload);

      // Check if response has data (success case)
      if (response.data) {
        setMessage({
          type: "success",
          text: `Tenant created successfully! Tenant ID: ${
            response.data.tenant_id || "Generated"
          }`,
        });

        // Reset form
        setTenantFormData({
          tenantName: "",
          domain: "",
          email: "",
        });

        // Redirect after a short delay to show success message
        setTimeout(() => {
          router.push("/dashboard");
        }, 2000);
      }
    } catch (error: any) {
      console.error("Tenant creation error:", error);

      // Handle specific error cases
      let errorMessage = "Failed to create tenant. Please try again.";

      if (
        error.message &&
        error.message.includes("duplicate key value violates unique constraint")
      ) {
        errorMessage =
          "This domain is already taken. Please choose a different domain.";
      } else if (
        error.message &&
        error.message.includes("tenants_domain_key")
      ) {
        errorMessage =
          "This domain is already in use. Please choose a different domain.";
      } else if (error.response?.data?.error) {
        // Handle API error response
        const apiError = error.response.data.error;
        if (
          typeof apiError === "string" &&
          apiError.includes("duplicate key")
        ) {
          errorMessage =
            "This domain is already taken. Please choose a different domain.";
        } else if (typeof apiError === "object" && apiError.error) {
          if (apiError.error.includes("duplicate key")) {
            errorMessage =
              "This domain is already taken. Please choose a different domain.";
          } else {
            errorMessage = apiError.error;
          }
        }
      }

      setMessage({
        type: "error",
        text: errorMessage,
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full space-y-8">
        <div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Create Tenant
          </h2>
          <p className="mt-2 text-center text-sm text-gray-600">
            Set up a new tenant organization
          </p>
        </div>

        {/* Success/Error Message */}
        {message && (
          <div
            className={`rounded-md p-4 ${
              message.type === "success"
                ? "bg-green-50 border border-green-200"
                : "bg-red-50 border border-red-200"
            }`}
          >
            <div className="flex">
              <div className="flex-shrink-0">
                {message.type === "success" ? (
                  <CheckCircleIcon className="h-5 w-5 text-green-400" />
                ) : (
                  <ExclamationTriangleIcon className="h-5 w-5 text-red-400" />
                )}
              </div>
              <div className="ml-3">
                <p
                  className={`text-sm font-medium ${
                    message.type === "success"
                      ? "text-green-800"
                      : "text-red-800"
                  }`}
                >
                  {message.text}
                </p>
              </div>
            </div>
          </div>
        )}

        <form className="mt-8 space-y-6" onSubmit={handleTenantSubmit}>
          <div className="rounded-md shadow-sm -space-y-px">
            <div>
              <input
                id="tenantName"
                name="tenantName"
                type="text"
                required
                className="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                placeholder="Tenant name (1-100 characters)"
                value={tenantFormData.tenantName}
                onChange={handleTenantChange}
                minLength={1}
                maxLength={100}
              />
            </div>
            <div>
              <input
                id="domain"
                name="domain"
                type="text"
                required
                className="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                placeholder="Domain (3-100 characters)"
                value={tenantFormData.domain}
                onChange={handleTenantChange}
                minLength={3}
                maxLength={100}
              />
            </div>
            <div>
              <input
                id="email"
                name="email"
                type="email"
                required
                className="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                placeholder="Email address"
                value={tenantFormData.email}
                onChange={handleTenantChange}
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
              {isLoading ? "Creating tenant..." : "Create Tenant"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
