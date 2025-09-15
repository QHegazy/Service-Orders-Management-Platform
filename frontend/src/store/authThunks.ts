import { createAsyncThunk } from "@reduxjs/toolkit";
import { UserClamis } from "@/types/user";
import authService from "@/services/authService";
import { ApiResponse } from "@/types/response";
import { jwtDecode } from "jwt-decode";
import { isValidJWTFormat } from "@/utils/tokenUtils";

interface JwtPayload {
  Data: {
    id: string;
    username: string;
    role: string;
    belong?: string[];
    exp: number;
  };
}

interface LoginResponse {
  access_token: string;
}

export const loginUser = createAsyncThunk(
  "auth/loginUser",
  async (
    { email, password }: { email: string; password: string },
    { rejectWithValue }
  ) => {
    try {
      console.log("AuthThunk: Starting login process");
      const res = (await authService.login({
        email,
        password,
      })) as ApiResponse<LoginResponse>;

      console.log("AuthThunk: Received response:", res);

      if (!res.data?.access_token) {
        console.error("AuthThunk: No access token in response");
        return rejectWithValue(res.message || "Login failed");
      }

      // Validate token format before decoding
      const token = res.data.access_token;
      if (!isValidJWTFormat(token)) {
        return rejectWithValue("Invalid token format received");
      }

      const decoded: JwtPayload = jwtDecode(token);

      const user: UserClamis = {
        id: decoded.Data.id,
        username: decoded.Data.username,
        role: decoded.Data.role,
        tenants: decoded.Data.belong || [],
      };

      return {
        user,
        token,
      };
    } catch (err: unknown) {
      console.error("Login error:", err);
      const errorMessage = err instanceof Error ? err.message : "Login failed";
      return rejectWithValue(errorMessage);
    }
  }
);

export const logoutUser = createAsyncThunk(
  "auth/logoutUser",
  async (_, { rejectWithValue }) => {
    try {
      await authService.logout();
      return true;
    } catch (err: unknown) {
      console.error("Logout error:", err);
      const errorMessage = err instanceof Error ? err.message : "Logout failed";
      return rejectWithValue(errorMessage);
    }
  }
);

export const refreshToken = createAsyncThunk(
  "auth/refreshToken",
  async (_, { rejectWithValue }) => {
    try {
      const res =
        (await authService.refreshToken()) as ApiResponse<LoginResponse>;

      if (!res.data?.access_token) {
        return rejectWithValue(res.message || "Token refresh failed");
      }

      // Validate token format before decoding
      const token = res.data.access_token;
      if (!isValidJWTFormat(token)) {
        return rejectWithValue("Invalid token format received");
      }

      // Decode JWT to get user info
      const decoded: JwtPayload = jwtDecode(token);

      const user: UserClamis = {
        id: decoded.Data.id,
        username: decoded.Data.username,
        role: decoded.Data.role,
        tenants: decoded.Data.belong || [],
      };

      return {
        user,
        token,
      };
    } catch (err: unknown) {
      console.error("Token refresh error:", err);
      const errorMessage =
        err instanceof Error ? err.message : "Token refresh failed";
      return rejectWithValue(errorMessage);
    }
  }
);
