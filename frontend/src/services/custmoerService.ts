import { ApiResponse } from "@/types/response";
import { apiRequest, apiRequestWithToken } from "./api";
import { CustomerSignUp } from "@/types/user";

class CustomerService {
  async signup(customerData: CustomerSignUp): Promise<ApiResponse<string>> {
    return apiRequest<ApiResponse<string>>("/customer", {
      method: "POST",
      body: JSON.stringify(customerData),
    });
  }

  async updateCustomer(
    customerData: Partial<CustomerSignUp>
  ): Promise<ApiResponse<string>> {
    return apiRequestWithToken<ApiResponse<string>>("customer", {
      method: "PUT",
      body: JSON.stringify(customerData),
    });
  }

  async deleteCustomer(): Promise<ApiResponse<string>> {
    return apiRequestWithToken<ApiResponse<string>>("/customer", {
      method: "DELETE",
    });
  }
}

const customerService = new CustomerService();

export default customerService;
