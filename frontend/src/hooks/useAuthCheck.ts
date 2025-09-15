import { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { AppDispatch, RootState } from "@/store";
import { loadFromStorage } from "@/store/authSlice";
import { refreshToken } from "@/store/authThunks";
import { jwtDecode } from "jwt-decode";
import { isValidJWTFormat } from "@/utils/tokenUtils";

interface JwtPayload {
  exp: number;
}

export const useAuthCheck = () => {
  const dispatch = useDispatch<AppDispatch>();
  const { token, isAuthenticated } = useSelector(
    (state: RootState) => state.auth
  );

  useEffect(() => {
    // Load auth state from localStorage on app start
    dispatch(loadFromStorage());
  }, [dispatch]);

  useEffect(() => {
    if (token && isAuthenticated) {
      // Validate token format first
      if (!isValidJWTFormat(token)) {
        console.error("Invalid token format, logging out");
        dispatch({ type: "auth/logout" });
        return;
      }

      try {
        const decoded: JwtPayload = jwtDecode(token);
        const currentTime = Date.now() / 1000;

        // Check if token is already expired
        if (decoded.exp <= currentTime) {
          console.log("Token expired, attempting refresh...");
          dispatch(refreshToken());
        } else if (decoded.exp - currentTime < 300) {
          // Check if token expires in the next 5 minutes (300 seconds)
          console.log("Token expiring soon, refreshing...");
          dispatch(refreshToken());
        }
      } catch (error) {
        console.error("Error decoding token:", error);
        // If token is invalid, logout instead of trying to refresh
        dispatch({ type: "auth/logout" });
      }
    }
  }, [token, isAuthenticated, dispatch]);

  // Set up periodic token check (every 5 minutes)
  useEffect(() => {
    if (!isAuthenticated) return;

    const interval = setInterval(() => {
      if (token) {
        // Validate token format
        if (!isValidJWTFormat(token)) {
          console.error("Invalid token format detected, logging out");
          dispatch({ type: "auth/logout" });
          return;
        }

        try {
          const decoded: JwtPayload = jwtDecode(token);
          const currentTime = Date.now() / 1000;

          // Check if token is expired
          if (decoded.exp <= currentTime) {
            console.log("Token expired during check, attempting refresh...");
            dispatch(refreshToken());
          } else if (decoded.exp - currentTime < 300) {
            // Check if token expires in the next 5 minutes
            console.log("Token expiring soon, refreshing...");
            dispatch(refreshToken());
          }
        } catch (error) {
          console.error("Error checking token expiration:", error);
          dispatch({ type: "auth/logout" });
        }
      }
    }, 5 * 60 * 1000); // Check every 5 minutes

    return () => clearInterval(interval);
  }, [token, isAuthenticated, dispatch]);
};
