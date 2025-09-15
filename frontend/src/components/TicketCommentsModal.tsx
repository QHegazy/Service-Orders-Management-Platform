"use client";

import { useState, useEffect, useRef } from "react";
import ticketService from "@/services/ticketService";
import webSocketService from "@/services/websocketService";

interface Comment {
  id: string;
  content: string;
  username: string;
  userRole: string;
  createdAt: string;
}

interface TicketCommentsModalProps {
  ticketId: string;
  onClose: () => void;
}

export default function TicketCommentsModal({
  ticketId,
  onClose,
}: TicketCommentsModalProps) {
  const [comments, setComments] = useState<Comment[]>([]);
  const [newComment, setNewComment] = useState("");
  const [loading, setLoading] = useState(true);
  const commentsEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const fetchComments = async () => {
      try {
        setLoading(true);
        // Assuming a method to list comments by ticket ID exists
        // You might need to implement this in ticketService
        const response = await ticketService.getComments(ticketId, 0, 100); // Adjust page/size as needed
        setComments(response.data.comments); // Access comments from response.data
      } catch (error) {
        console.error("Failed to fetch comments:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchComments();

    // Connect to WebSocket
    const token = localStorage.getItem("access_token"); // Get token from localStorage
    if (token) {
      webSocketService.connect(ticketId, token);
      webSocketService.onMessage((message: any) => {
        // Assuming message is already parsed JSON
        setComments((prevComments) => [...prevComments, message]);
      });
    } else {
      console.error("No access token found for WebSocket connection.");
    }

    return () => {
      webSocketService.disconnect();
    };
  }, [ticketId]);

  useEffect(() => {
    commentsEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [comments]);

  const handleSendComment = (e: React.FormEvent) => {
    e.preventDefault();
    if (newComment.trim()) {
      webSocketService.sendMessage({ content: newComment });
      setNewComment("");
    }
  };

  return (
    <div className="fixed inset-0  backdrop-blur-md  bg-opacity-50 flex items-center justify-center z-50  text-cyan-500">
      <div className="bg-white p-6 rounded-lg w-1/2 max-h-[80vh] flex flex-col">
        <h2 className="text-xl font-bold mb-4">Comments for Ticket {ticketId}</h2>
        <div className="flex-grow overflow-y-auto border p-4 mb-4 rounded">
          {loading ? (
            <div className="text-center">Loading comments...</div>
          ) : comments.length === 0 ? (
            <div className="text-center text-gray-500">No comments yet.</div>
          ) : (
            comments.map((comment, index) => (
              <div key={index} className="mb-2 p-2 bg-gray-100 rounded">
                <p className="font-semibold">{comment.username} ({comment.userRole})</p>
                <p>{comment.content}</p>
                <p className="text-xs text-gray-500 text-right">{new Date(comment.createdAt).toLocaleString()}</p>
              </div>
            ))
          )}
          <div ref={commentsEndRef} />
        </div>
        <form onSubmit={handleSendComment} className="flex gap-2">
          <input
            type="text"
            value={newComment}
            onChange={(e) => setNewComment(e.target.value)}
          placeholder="Add a comment..."
            className="flex-grow border rounded px-3 py-2"
          />
          <button
            type="submit"
            className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
          >
            Send
          </button>
        </form>
        <div className="mt-4 flex justify-end">
          <button
            onClick={onClose}
            className="bg-gray-500 text-white px-4 py-2 rounded hover:bg-gray-600"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  );
}