"use client";

import { useState, useEffect, useCallback } from "react";
import { useRouter } from "next/navigation";
import ticketService from "@/services/ticketService";
import TicketPropertySelect from "@/components/TicketPropertySelect";
import TicketCommentsModal from "@/components/TicketCommentsModal"; // Import the new modal
import { Selector } from "react-redux";

// The Ticket interface remains the same
interface Ticket {
  id: string;
  title: string;
  description: string;
  status: string;
  priority: string;
  customer_first_name: string;
  customer_last_name: string;
  assigned_username: string;
  created_at: string;
}

// Constants for select options are unchanged
const statusOptions = [
  { value: "OPEN", label: "Open" },
  { value: "IN_PROGRESS", label: "In Progress" },
  { value: "RESOLVED", label: "Resolved" },
  { value: "CLOSED", label: "Closed" },
  { value: "REOPENED", label: "Reopened" },
];

const priorityOptions = [
  { value: "LOW", label: "Low" },
  { value: "MEDIUM", label: "Medium" },
  { value: "HIGH", label: "High" },
  { value: "URGENT", label: "Urgent" },
  { value: "CRITICAL", label: "Critical" },
];

const tableHeaders = ["Title", "Customer", "Status", "Priority", "Assigned To", "Created", "Actions"];

export default function TicketsPage() {
  const router = useRouter();
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [loading, setLoading] = useState(true);
  const [userRole, setUserRole] = useState("");
  const [showCommentsModal, setShowCommentsModal] = useState(false); // New state for modal visibility
  const [selectedTicketId, setSelectedTicketId] = useState<string | null>(null); // New state for selected ticket ID

  const canEdit = userRole === "Admin" || userRole === "Technician";

  // Fetches tickets from the server
  const fetchTickets = useCallback(async () => {
    setLoading(true);
    try {
      const user = JSON.parse(localStorage.getItem("user") || "{}");
      const role = user.role || "";

      if (role === "Admin" || role === "Technician") {
        const response = await ticketService.getUserTickets(50, 0);
        setTickets(response);
      } else {
        const response = await ticketService.getCustomersTickets();
        setTickets(response);
      }
    } catch (error) {
      console.error("Failed to fetch tickets:", error);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    const user = JSON.parse(localStorage.getItem("user") || "{}");
    setUserRole(user.role || "");
    fetchTickets();
  }, [fetchTickets]);

  const handlePropertyChange = async (
    ticketId: string,
    property: "status" | "priority",
    value: string
  ) => {
    setTickets((prevTickets) =>
      prevTickets.map((ticket) =>
        ticket.id === ticketId ? { ...ticket, [property]: value } : ticket
      )
    );
  };

  const handleViewComments = (ticketId: string) => {
    setSelectedTicketId(ticketId);
    setShowCommentsModal(true);
  };

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-4">Tickets</h1>
      <div className="bg-white rounded-lg shadow overflow-x-auto">
        <table className="w-full">
          <thead className="bg-gray-50">
            <tr>
              {tableHeaders.map((header) => (
                <th
                  key={header}
                  className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  {header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {tickets.map((ticket) => (
              <tr key={ticket.id}>
                <td className="px-6 py-4 whitespace-nowrap font-medium text-gray-900">
                  {ticket.title}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {`${ticket.customer_first_name} ${ticket.customer_last_name}`}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <TicketPropertySelect
                    value={ticket.status}
                    onChange={(e) => handlePropertyChange(ticket.id, "status", e.target.value)}
                    options={statusOptions}
                    disabled={!canEdit}
                  />
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <TicketPropertySelect
                    value={ticket.priority}
                    onChange={(e) => handlePropertyChange(ticket.id, "priority", e.target.value)}
                    options={priorityOptions}
                    disabled={!canEdit}
                  />
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {ticket.assigned_username || "Unassigned"}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {new Date(ticket.created_at).toLocaleDateString()}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <button
                    onClick={() => handleViewComments(ticket.id)} // Changed onClick
                    className="text-indigo-600 hover:text-indigo-900"
                  >
                    View
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {loading && <div className="text-center py-4">Loading tickets...</div>}
        {!loading && tickets.length === 0 && (
          <div className="text-center py-4 text-gray-500">No tickets found.</div>
        )}
      </div>

      {showCommentsModal && selectedTicketId && (
        <TicketCommentsModal
          ticketId={selectedTicketId}
          onClose={() => setShowCommentsModal(false)}
        />
      )}
    </div>
  );
}
