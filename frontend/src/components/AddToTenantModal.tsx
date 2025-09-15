import React, { useState } from "react";
import { Tenant } from "@/types/tenant";

interface AddToTenantModalProps {
  technicianId: string | null;
  tenants: Tenant[];
  onSuccess: (tenantId: string) => void;
  onCancel: () => void;
}

export const AddToTenantModal: React.FC<AddToTenantModalProps> = ({
  technicianId,
  tenants,
  onSuccess,
  onCancel,
}) => {
  const [selectedTenant, setSelectedTenant] = useState<string>("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (selectedTenant) {
      onSuccess(selectedTenant);
    } else {
      alert("Please select a tenant.");
    }
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow-lg w-full max-w-md text-blue-600">
      <h2 className="text-xl font-bold mb-4">Add Technician to Tenant</h2>
      {technicianId ? (
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label htmlFor="tenant-select" className="block text-gray-700 text-sm font-bold mb-2">
              Select Tenant:
            </label>
            <select
              id="tenant-select"
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              value={selectedTenant}
              onChange={(e) => setSelectedTenant(e.target.value)}
              required
            >
              <option value="">-- Select a Tenant --</option>
              {tenants.map((tenant) => (
                <option key={tenant.id} value={tenant.id}>
                  {tenant.tenant_name}
                </option>
              ))}
            </select>
          </div>
          <div className="flex justify-end">
            <button
              type="button"
              onClick={onCancel}
              className="bg-gray-300 text-gray-800 px-4 py-2 rounded mr-2 hover:bg-gray-400"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
            >
              Add
            </button>
          </div>
        </form>
      ) : (
        <p>No technician selected.</p>
      )}
    </div>
  );
};
