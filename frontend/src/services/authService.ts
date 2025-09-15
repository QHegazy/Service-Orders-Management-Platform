import { ApiResponse } from "@/types/response";
import { apiRequest, apiRequestWithToken } from "./api";

type loginForm = {
  email: string;
  password: string;
};
class AuthService {
  async login(loginData: loginForm): Promise<ApiResponse<unknown>> {
    console.log("AuthService: Attempting login with:", {
      email: loginData.email,
    });
    try {
      const result = await apiRequest<ApiResponse<unknown>>(`/auth/login`, {
        method: "POST",
        body: JSON.stringify(loginData),
      });
      console.log("AuthService: Login successful:", result);
      return result;
    } catch (error) {
      console.error("AuthService: Login failed:", error);
      throw error;
    }
  }
  async logout(): Promise<ApiResponse<unknown>> {
    localStorage.removeItem("access_token");
    return apiRequestWithToken<ApiResponse<unknown>>("/auth/logout", {
      method: "POST",
    });
  }
  async refreshToken(): Promise<ApiResponse<unknown>> {
    const res = await apiRequestWithToken<ApiResponse<unknown>>(
      "/auth/refresh",
      {
        method: "POST",
      }
    );
    return res;
  }
}

const authService = new AuthService();

export default authService;
