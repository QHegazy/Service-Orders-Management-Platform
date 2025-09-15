import { apiRequestWithToken } from "./api";
import { TenantSignUp, TenantUpdate, Tenant } from "@/types/tenant";
import { ApiResponse } from "@/types/response";

interface AddUserToTenant {
  user_id: string;
  tenant_id: string;
}
class TenantService {
  async createTenant(tenantData: TenantSignUp): Promise<ApiResponse<string>> {
    return apiRequestWithToken<ApiResponse<string>>("/tenant", {
      method: "POST",
      body: JSON.stringify(tenantData),
    });
  }

  async updateTenant(tenantData: TenantUpdate): Promise<ApiResponse<string>> {
    return apiRequestWithToken<ApiResponse<string>>("/tenant", {
      method: "PUT",
      body: JSON.stringify(tenantData),
    });
  }

  async deleteTenant(tenantId: string): Promise<ApiResponse<string>> {
    return apiRequestWithToken<ApiResponse<string>>(`/tenant/${tenantId}`, {
      method: "DELETE",
    });
  }

  async getUserTenants(): Promise<Tenant[]> {
    try {
      const response = await apiRequestWithToken<ApiResponse<Tenant[]>>("/tenant/user", {
        method: "GET",
      });
      return response.data;
    } catch (err) {
      console.error("Service error:", err);
      throw new Error("Failed to fetch user tenants");
    }
  }
  async getAllTenants(
    limit: number = 50,
    offset: number = 0
  ): Promise<Tenant[]> {
    try {
      const response = await apiRequestWithToken<ApiResponse<Tenant[]>>(`/tenant/all?limit=${limit}&offset=${offset}`, {
        method: "GET",
      });
      return response.data;
    } catch (err) {
      console.error("Service error:", err);
      throw new Error("Failed to fetch all tenants");
    }
  }
  async addTechnicianToTenant( payload :AddUserToTenant): Promise<string> {
    try {
      const response = await apiRequestWithToken<ApiResponse<string>>(
        `/tenant/add-to`,
        {
          method: "POST",
          body: JSON.stringify(payload),
        }
      );
      return response.data;
    } catch (err) {
      console.error("Service error:", err);
      throw new Error("Failed to add technician to tenant");
    }
  
  }  
}
const tenantService = new TenantService();

export { tenantService };
export default tenantService;
