import { apiRequest, apiRequestWithToken } from "./api";
import { UserSignUp, Login, User } from "@/types/user";
import { ApiResponse } from "@/types/response";

interface CreateTechnicianData {
  username: string;
  email: string;
  password: string;
}

class UserService {

  

  async signup(userData: UserSignUp): Promise<ApiResponse<string>> {
    return apiRequest<ApiResponse<string>>("/user", {
      method: "POST",
      body: JSON.stringify(userData),
    });
  }

  async updateUser(
    userData: Partial<UserSignUp>
  ): Promise<ApiResponse<string>> {
    return apiRequestWithToken<ApiResponse<string>>("/user", {
      method: "PUT",
      body: JSON.stringify(userData),
    });
  }

  async deleteUser(userId?: string): Promise<ApiResponse<string>> {
    const url = userId ? `/user/${userId}` : "/user";
    return apiRequestWithToken<ApiResponse<string>>(url, {
      method: "DELETE",
    });
  }

  async createTechnician(
    userData: CreateTechnicianData
  ): Promise<ApiResponse<string>> {
    return apiRequestWithToken<ApiResponse<string>>("/user/technician", {
      method: "POST",
      body: JSON.stringify(userData),
    });
  }

async getTechnicians(): Promise<User[]> {
  try {
    const response = await apiRequestWithToken<ApiResponse<User[]>>(
      "/user/technicians",
      { method: "GET" }
    );

    return response.data

  } catch (err) {
    console.error("Service error:", err);
    throw new Error("Failed to fetch technicians");
  }
}
  async updateProfile(
    userData: Partial<UserSignUp>
  ): Promise<ApiResponse<string>> {
    return apiRequestWithToken<ApiResponse<string>>("/user", {
      method: "PUT",
      body: JSON.stringify(userData),
    });
  }
}
const userService = new UserService();

export { userService };
export default userService;
