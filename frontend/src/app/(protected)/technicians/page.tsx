"use client";

import { useState, useEffect } from "react";
import { userService } from "@/services/userService";
import { User } from "@/types/user";
import { TechnicianForm } from "@/components/form/TechnicianForm";
import { tenantService } from "@/services/tenantService"; // Import tenantService
import { Tenant } from "@/types/tenant"; // Import Tenant type
import { AddToTenantModal } from "@/components/AddToTenantModal"; // Import AddToTenantModal

export default function TechniciansPage() {
  const [technicians, setTechnicians] = useState<User[]>([]);
  const [showForm, setShowForm] = useState(false);
  const [loading, setLoading] = useState(false);
  const [showAddToTenantModal, setShowAddToTenantModal] = useState(false);
  const [selectedTechnicianId, setSelectedTechnicianId] = useState<string | null>(null);
  const [userTenants, setUserTenants] = useState<Tenant[]>([]); // State to store user's tenants

  useEffect(() => {
    fetchTechnicians();
    fetchUserTenants(); // Fetch user's tenants when component mounts
  }, []);

  const fetchUserTenants = async () => {
    try {
      const tenants = await tenantService.getUserTenants();
      setUserTenants(tenants);
    } catch (error) {
      console.error("Failed to fetch user tenants:", error);
    }
  };

  const handleAddToTenantClick = (technicianId: string) => {
    setSelectedTechnicianId(technicianId);
    setShowAddToTenantModal(true);
  };

  const handleAddToTenantSuccess = async (tenantId: string) => {
    if (selectedTechnicianId) {
      try {
        await tenantService.addTechnicianToTenant({
          user_id: selectedTechnicianId,
          tenant_id: tenantId,
        });
        alert("Technician added to tenant successfully!");
        setShowAddToTenantModal(false);
        setSelectedTechnicianId(null);
      } catch (error) {
        console.error("Failed to add technician to tenant:", error);
        alert("Failed to add technician to tenant.");
      }
    }
  };

  const fetchTechnicians = async () => {
    try {
      setLoading(true);
      const response = await userService.getTechnicians();
      setTechnicians(response);
    } catch (error) {
      console.error("Failed to fetch technicians:", error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTechnicians();
  }, []);

  const handleDelete = async (id: string) => {
    if (!confirm("Are you sure you want to delete this technician?")) return;

    try {
      await userService.deleteUser(id);
      fetchTechnicians(); // Refresh the list
    } catch (error) {
      console.error("Failed to delete technician:", error);
    }
  };

  const handleFormSuccess = () => {
    setShowForm(false);
    fetchTechnicians(); // Refresh the list
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Technicians</h1>
        <button
          onClick={() => setShowForm(true)}
          className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
        >
          Add Technician
        </button>
      </div>
      <div className="bg-white rounded-lg shadow">
        <table className="w-full">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-black uppercase">
                Username
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-black uppercase">
                Email
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-black uppercase">
                Status
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-black uppercase">
                Created
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-black uppercase">
                Actions
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200 text-amber-400">
            {technicians.map((tech) => (
              <tr key={tech.id}>
                <td className="px-6 py-4 whitespace-nowrap  text-blue-500">
                  {tech.username}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-blue-500">
                  {tech.email}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span
                    className={`px-2 py-1 text-xs rounded ${
                      tech.is_active
                        ? "bg-green-100 text-green-800"
                        : "bg-red-100 text-red-800"
                    }`}
                  >
                    {tech.is_active ? "Active" : "Inactive"}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  {new Date(tech.created_at).toLocaleDateString()}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <button
                    onClick={() => handleDelete(tech.id)}
                    className="text-red-600 hover:text-red-900 mr-2"
                  >
                    Delete
                  </button>
                  <button
                    onClick={() => handleAddToTenantClick(tech.id)}
                    className="text-blue-600 hover:text-blue-900"
                  >
                    Add to Tenant
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {loading && <div className="text-center py-4">Loading...</div>}
      </div>

      {/* Modal */}
      {showForm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <TechnicianForm
            onSuccess={handleFormSuccess}
            onCancel={() => setShowForm(false)}
          />
        </div>
      )}

      {/* Add to Tenant Modal */}
      {showAddToTenantModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <AddToTenantModal
            technicianId={selectedTechnicianId}
            tenants={userTenants}
            onSuccess={handleAddToTenantSuccess}
            onCancel={() => setShowAddToTenantModal(false)}
          />
        </div>
      )}
    </div>
  );
}
