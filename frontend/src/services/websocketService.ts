import { WebSocketMessage } from "@/types/websocket"; // Add this import

class WebSocketService {
  private socket: WebSocket | null = null;
  private onMessageCallback: ((message: WebSocketMessage) => void) | null = null; // Changed from any to WebSocketMessage

  connect(ticketId: string) {
    if (this.socket) {
      this.disconnect();
    }
    const token = localStorage.getItem("access_token");
    if (!token) {
      console.error("No token found in localStorage");
      return;
    }
    
    const url = `ws://localhost:8080/v1/ws/ticket/${ticketId}?token=${token}`;

    this.socket = new WebSocket(url);

    this.socket.onopen = () => {
      console.log("WebSocket connected");
    };

    this.socket.onmessage = (event) => {
      const message = JSON.parse(event.data);
      if (this.onMessageCallback) {
        this.onMessageCallback(message);
      }
    };

    this.socket.onclose = () => {
      console.log("WebSocket disconnected");
      this.socket = null;
    };

    this.socket.onerror = (error) => {
      console.error("WebSocket error:", error);
    };
  }

  disconnect() {
    if (this.socket) {
      this.socket.close();
    }
  }

  onMessage(callback: (message: WebSocketMessage) => void) { // Changed from any to WebSocketMessage
    this.onMessageCallback = callback;
  }

  sendMessage(message: { content: string }) {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify(message));
    } else {
      console.error("WebSocket is not connected.");
    }
  }
}

const webSocketService = new WebSocketService();
export default webSocketService;
