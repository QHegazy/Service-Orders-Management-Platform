"use client";

import { useState, useEffect } from "react";
import { tenantService } from "@/services/tenantService";
import { Tenant } from "@/types/tenant";
import Modal from "@/components/Modal"; // Import the Modal component
import { TenantForm } from "@/components/TenantForm"; // Import the TenantForm component

export default function TenantsPage() {
  const [tenants, setTenants] = useState<Tenant[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showCreateTenantModal, setShowCreateTenantModal] = useState(false); // State for modal visibility

  const fetchTenants = async () => {
    try {
      setLoading(true);
      setError(null);
      console.log("Fetching tenants...");
      const response = await tenantService.getUserTenants();
      console.log("Tenants response:", response);

      // Ensure we have data and it's an array
      if (response && Array.isArray(response)) {
        setTenants(response);
      } else {
        console.warn("Invalid response format:", response);
        setTenants([]);
      }
    } catch (error) {
      console.error("Failed to fetch tenants:", error);
      const errorMessage =
        error instanceof Error ? error.message : "Failed to fetch tenants";
      setError(errorMessage);
      setTenants([]); // Set empty array on error
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTenants();
  }, []);

  const handleCreateTenantSuccess = () => {
    fetchTenants(); // Refresh the list of tenants
    setShowCreateTenantModal(false); // Close the modal
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">My Tenants</h1>
        <button
          onClick={() => setShowCreateTenantModal(true)}
          className="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        >
          Create New Tenant
        </button>
      </div>

      <div className="bg-white rounded-lg shadow text-amber-400">
        <table className="w-full">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                Name
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                Domain
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                Email
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                Status
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                Created
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {tenants && tenants.length > 0
              ? tenants.map((tenant) => (
                  <tr key={tenant.id}>
                    <td className="px-6 py-4 whitespace-nowrap font-medium">
                      {tenant.tenant_name}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      {tenant.domain}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      {tenant.email}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span
                        className={`px-2 py-1 text-xs rounded ${
                          tenant.is_active
                            ? "bg-green-100 text-green-800"
                            : "bg-red-100 text-red-800"
                        }`}
                      >
                        {tenant.is_active ? "Active" : "Inactive"}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      {new Date(tenant.created_at).toLocaleDateString()}
                    </td>
                  </tr>
                ))
              : !loading && (
                  <tr>
                    <td
                      colSpan={5}
                      className="px-6 py-4 text-center text-gray-500"
                    >
                      No tenants found
                    </td>
                  </tr>
                )}
          </tbody>
        </table>
        {loading && <div className="text-center py-4">Loading...</div>}
      </div>

      {/* Create Tenant Modal */}
      <Modal
        isOpen={showCreateTenantModal}
        onClose={() => setShowCreateTenantModal(false)}
        title="Create New Tenant"
      >
        <TenantForm
          onSuccess={handleCreateTenantSuccess}
          onClose={() => setShowCreateTenantModal(false)}
        />
      </Modal>
    </div>
  );
}
