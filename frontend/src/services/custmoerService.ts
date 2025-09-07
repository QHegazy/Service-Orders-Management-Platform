import { apiRequest } from "./api";

class CustomerService {
  async signup(username: string, email: string, password: string) {
    return apiRequest<{ message: string; userId: string }>("/customer/signup", {
      method: "POST",
      body: JSON.stringify({ username, email, password }),
    });
  }
}

const customerService = new CustomerService();

export default customerService;