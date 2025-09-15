"use client";

import { useState } from "react";

interface CreateTicketModalProps {
  onClose: () => void;
  onTicketCreated: (newTicket: { title: string; description: string; priority: string }) => Promise<void>;
}

const priorityOptions = [
  { value: "LOW", label: "Low" },
  { value: "MEDIUM", label: "Medium" },
  { value: "HIGH", label: "High" },
  { value: "URGENT", label: "Urgent" },
  { value: "CRITICAL", label: "Critical" },
];

export default function CreateTicketModal({
  onClose,
  onTicketCreated,
}: CreateTicketModalProps) {
  const [newTicket, setNewTicket] = useState({
    title: "",
    description: "",
    priority: "MEDIUM",
  });
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleCreateTicket = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    try {
      await onTicketCreated(newTicket);
      onClose();
    } catch (error) {
      console.error("Failed to create ticket:", error);
      // Optionally, show an error message to the user
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white p-6 rounded-lg w-96">
        <h2 className="text-xl font-bold mb-4">Create New Ticket</h2>
        <form onSubmit={handleCreateTicket}>
          <div className="mb-4">
            <label className="block text-sm font-medium mb-2">Title</label>
            <input
              type="text"
              value={newTicket.title}
              onChange={(e) =>
                setNewTicket({ ...newTicket, title: e.target.value })
              }
              className="w-full border rounded px-3 py-2"
              required
              disabled={isSubmitting}
            />
          </div>
          <div className="mb-4">
            <label className="block text-sm font-medium mb-2">
              Description
            </label>
            <textarea
              value={newTicket.description}
              onChange={(e) =>
                setNewTicket({ ...newTicket, description: e.target.value })
              }
              className="w-full border rounded px-3 py-2 h-24"
              required
              disabled={isSubmitting}
            />
          </div>
          <div className="mb-4">
            <label className="block text-sm font-medium mb-2">
              Priority
            </label>
            <select
              value={newTicket.priority}
              onChange={(e) =>
                setNewTicket({ ...newTicket, priority: e.target.value })
              }
              className="w-full border rounded px-3 py-2"
              disabled={isSubmitting}
            >
              {priorityOptions.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </div>
          <div className="flex gap-2 justify-end">
            <button
              type="button"
              onClick={onClose}
              className="bg-gray-500 text-white px-4 py-2 rounded hover:bg-gray-600 disabled:opacity-50"
              disabled={isSubmitting}
            >
              Cancel
            </button>
            <button
              type="submit"
              className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 disabled:opacity-50"
              disabled={isSubmitting}
            >
              {isSubmitting ? "Creating..." : "Create"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
