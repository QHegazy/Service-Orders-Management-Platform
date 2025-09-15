import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { UserClaims } from "@/types/user";
import { loginUser, refreshToken, logoutUser } from "./authThunks";
import { isValidJWTFormat } from "@/utils/tokenUtils";

type AuthState = {
  user: UserClaims | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
};

const initialState: AuthState = {
  user: null,
  token: null,
  isAuthenticated: false,
  isLoading: false,
  error: null,
};

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    login: (
      state,
      action: PayloadAction<{ user: UserClaims; token: string }>
    ) => {
      state.user = action.payload.user;
      state.token = action.payload.token;
      state.isAuthenticated = true;
      localStorage.setItem("access_token", action.payload.token);
      localStorage.setItem("user", JSON.stringify(action.payload.user));
      document.cookie = `access_token=${action.payload.token}; path=/; max-age=86400; SameSite=Lax`;
    },
    logout: (state) => {
      state.user = null;
      state.token = null;
      state.isAuthenticated = false;
      state.error = null;
      localStorage.removeItem("access_token");
      localStorage.removeItem("user");
      document.cookie =
        "access_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT";
    },
    loadFromStorage: (state) => {
      const token = localStorage.getItem("access_token");
      const userData = localStorage.getItem("user");

      // Validate token format before using it
      if (token && userData && isValidJWTFormat(token)) {
        try {
          const parsedUser = JSON.parse(userData);
          state.token = token;
          state.user = parsedUser;
          state.isAuthenticated = true;
          document.cookie = `access_token=${token}; path=/; max-age=86400; SameSite=Lax`;
        } catch (error) {
          // Clear invalid data
          localStorage.removeItem("access_token");
          localStorage.removeItem("user");
          document.cookie =
            "access_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT";
        }
      } else {
        // Clear invalid tokens
        localStorage.removeItem("access_token");
        localStorage.removeItem("user");
        document.cookie =
          "access_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT";
      }
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(loginUser.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(loginUser.fulfilled, (state, action) => {
        state.isLoading = false;
        state.user = action.payload.user;
        state.token = action.payload.token;
        state.isAuthenticated = true;
        state.error = null;
        localStorage.setItem("access_token", action.payload.token);
        localStorage.setItem("user", JSON.stringify(action.payload.user));
        document.cookie = `access_token=${action.payload.token}; path=/; max-age=86400; SameSite=Lax`;
      })
      .addCase(loginUser.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Refresh token cases
      .addCase(refreshToken.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(refreshToken.fulfilled, (state, action) => {
        state.isLoading = false;
        state.user = action.payload.user;
        state.token = action.payload.token;
        state.isAuthenticated = true;
        state.error = null;
        localStorage.setItem("access_token", action.payload.token);
        localStorage.setItem("user", JSON.stringify(action.payload.user));
        document.cookie = `access_token=${action.payload.token}; path=/; max-age=86400; SameSite=Lax`;
      })
      .addCase(refreshToken.rejected, (state) => {
        state.isLoading = false;
        state.user = null;
        state.token = null;
        state.isAuthenticated = false;
        state.error = "Session expired. Please login again.";
        localStorage.removeItem("access_token");
        localStorage.removeItem("user");
        document.cookie =
          "access_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT";
      })
      // Logout cases
      .addCase(logoutUser.fulfilled, (state) => {
        state.user = null;
        state.token = null;
        state.isAuthenticated = false;
        state.error = null;
        localStorage.removeItem("access_token");
        localStorage.removeItem("user");
        document.cookie =
          "access_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT";
      });
  },
});

export const { login, logout, loadFromStorage } = authSlice.actions;
export default authSlice.reducer;
