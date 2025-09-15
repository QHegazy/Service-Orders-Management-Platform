type TicketStatus = "OPEN" | "IN_PROGRESS" | "RESOLVED" | "CLOSED";
type TicketPriority = "LOW" | "MEDIUM" | "HIGH" | "URGENT";

interface Ticket {
  id?: string;
  tenant_id: string;
  customer_id?: string;
  assigned_to?: string | null;
  title: string;
  description?: string | null;
  status?: TicketStatus;
  priority: TicketPriority;
  created_at?: string;
  updated_at?: string;
}

interface Comment {
  id: string;
  ticket_id: string;
  user_id?: string;
  username?: string;
  comment: string;
  created_at: string;
}

interface WSMessage {
  type:
    | "AUTH"
    | "SUBSCRIBE_TICKET"
    | "COMMENT"
    | "NEW_COMMENT"
    | "ERROR"
    | "CONNECTION";
  token?: string;
  ticketId?: string;
  comment?: string;
  data?: Comment;
  error?: string;
  message?: string;
  connectionId?: string;
}

export type { Ticket, Comment, WSMessage, TicketStatus, TicketPriority };
