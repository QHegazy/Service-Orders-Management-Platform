import { useState, useEffect } from "react";
import tenantService from "@/services/tenantService";
import { Tenant } from "@/types/tenant";
import TicketForm from "@/components/form/ticketForm";

// Modal wrapper for the TicketForm
function TicketFormModal({
  isOpen,
  onClose,
  tenant,
}: {
  isOpen: boolean;
  onClose: () => void;
  tenant: Tenant | null;
}) {
  if (!isOpen || !tenant) return null;

  return (
    <div className="fixed inset-0  backdrop-blur-md bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl max-w-md w-full mx-4 max-h-[90vh] overflow-y-auto">
        <div className="p-6">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-xl font-semibold text-gray-900">
              Create Ticket for {tenant.tenant_name}
            </h2>
            <button
              onClick={onClose}
              className="text-gray-400 hover:text-gray-600 text-2xl"
            >
              ×
            </button>
          </div>

          <div className="ticket-form-wrapper">
            <TicketForm tenantId={tenant.id} />
          </div>

          <div className="mt-4 pt-4 border-t">
            <button
              onClick={onClose}
              className="w-full bg-gray-300 text-gray-700 py-2 px-4 rounded-md hover:bg-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500"
            >
              Close
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export function CustomerLayout() {
  const [tenants, setTenants] = useState<Tenant[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedTenant, setSelectedTenant] = useState<Tenant | null>(null);

  useEffect(() => {
    const fetchTenants = async () => {
      try {
        const response = await tenantService.getAllTenants();
        setTenants(response);
      } catch (err) {
        setError("Failed to fetch tenants");
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchTenants();
  }, []);

  // Generate avatar initials from tenant name
  const getInitials = (name: string) => {
    return name
      .split(" ")
      .map((word) => word.charAt(0))
      .join("")
      .toUpperCase()
      .slice(0, 2);
  };

  // Generate a consistent color based on tenant name
  const getAvatarColor = (name: string) => {
    const colors = [
      "bg-blue-500",
      "bg-green-500",
      "bg-purple-500",
      "bg-pink-500",
      "bg-indigo-500",
      "bg-yellow-500",
      "bg-red-500",
      "bg-teal-500",
    ];
    const index = name.length % colors.length;
    return colors[index];
  };

  const handleTenantClick = (tenant: Tenant) => {
    setSelectedTenant(tenant);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setSelectedTenant(null);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-700">
          {error}
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Tenants</h1>
          <p className="mt-2 text-gray-600">Manage your tenant organizations</p>
        </div>

        {tenants.length === 0 ? (
          <div className="text-center py-12">
            <div className="text-gray-500 text-lg">No tenants found</div>
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
            {tenants.map((tenant) => (
              <div
                key={tenant.id}
                onClick={() => handleTenantClick(tenant)}
                className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow duration-200 cursor-pointer"
              >
                <div className="flex flex-col items-center text-center">
                  {/* Avatar */}
                  <div
                    className={`w-16 h-16 rounded-full flex items-center justify-center text-white font-semibold text-lg mb-4 ${getAvatarColor(
                      tenant.tenant_name
                    )}`}
                  >
                    {getInitials(tenant.tenant_name)}
                  </div>

                  {/* Tenant Name */}
                  <h3 className="text-lg font-medium text-gray-900 mb-2 line-clamp-2">
                    {tenant.tenant_name}
                  </h3>

                  <div className="text-sm text-gray-500 truncate max-w-[200px]">
                    {tenant.email}
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Ticket Form Modal */}
        <TicketFormModal
          isOpen={isModalOpen}
          onClose={handleCloseModal}
          tenant={selectedTenant}
        />
      </div>
    </div>
  );
}
