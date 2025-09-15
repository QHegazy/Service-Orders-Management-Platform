import { apiRequestWithToken } from "./api";
import { Ticket, Comment } from "@/types/ticket";
import { ApiResponse } from "@/types/response";

class TicketService {
  async createTicket(ticketData: Ticket): Promise<ApiResponse<string>> {
    return apiRequestWithToken<ApiResponse<string>>("/ticket", {
      method: "POST",
      body: JSON.stringify(ticketData),
    });
  }

  async getTicket(ticketId: string): Promise<ApiResponse<Ticket>> {
    return apiRequestWithToken<ApiResponse<Ticket>>(`/ticket/${ticketId}`, {
      method: "GET",
    });
  }

  async updateTicket(
    ticketId: string,
    ticketData: Partial<Ticket>
  ): Promise<ApiResponse<string>> {
    return apiRequestWithToken<ApiResponse<string>>(`/ticket/${ticketId}`, {
      method: "PUT",
      body: JSON.stringify(ticketData),
    });
  }

  async deleteTicket(ticketId: string): Promise<ApiResponse<string>> {
    return apiRequestWithToken<ApiResponse<string>>(`/ticket/${ticketId}`, {
      method: "DELETE",
    });
  }

  async getComments(
    ticketId: string,
    limit: number = 50,
    offset: number = 0
  ): Promise<ApiResponse<Comment[]>> {
    return apiRequestWithToken<ApiResponse<Comment[]>>(
      `/ticket/${ticketId}/comments?limit=${limit}&offset=${offset}`,
      {
        method: "GET",
      }
    );
  }
  async getCustomersTickets(
    limit: number = 50,
    offset: number = 0
  ): Promise<Ticket[]> {
    try {
      const response = await apiRequestWithToken<ApiResponse<{ tickets: Ticket[] }>>(
        `/ticket/customer?limit=${limit}&offset=${offset}`,
        {
          method: "GET",
        }
      );
      return response.data.tickets;
    } catch (error) {
      console.error(`Error fetching customer tickets:`, error);
      throw error;
    }
  }
  async getTechnicianTickets(
    limit: number = 50,
    offset: number = 0
  ): Promise<Ticket[]> {
    try {
      const response = await apiRequestWithToken<ApiResponse<{ tickets: Ticket[] }>>(
        `/ticket/technician?limit=${limit}&offset=${offset}`,
        {
          method: "GET",
        }
      );
      return response.data.tickets;
    } catch (error) {
      console.error(`Error fetching technician tickets:`, error);
      throw error;
    }
  }
  async getTenantTickets(
    limit: number = 50,
    offset: number = 0
  ): Promise<Ticket[]> {
    try {
      const response = await apiRequestWithToken<ApiResponse<{ tickets: Ticket[] }>>(
        `/ticket/tenant?limit=${limit}&offset=${offset}`,
        {
          method: "GET",
        }
      );
      return response.data.tickets;
    } catch (error) {
      console.error(`Error fetching tenant tickets:`, error);
      throw error;
    }
  }
  async getUserTickets(
    limit: number = 50,
    offset: number = 0
  ): Promise<Ticket[]> {
    try {
      const response = await apiRequestWithToken<ApiResponse<{ tickets: Ticket[] }>>(
        `/ticket/user?limit=${limit}&offset=${offset}`,
        {
          method: "GET",
        }
      );
      return response.data.tickets;
    } catch (error) {
      console.error(`Error fetching user tickets:`, error);
      throw error;
    }
  }
  
}

const ticketService = new TicketService();

export default ticketService;
