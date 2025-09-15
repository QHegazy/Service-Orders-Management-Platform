const API_URL = process.env.NEXT_PUBLIC_API_URL;

let isRefreshing = false;
let failedQueue: Array<{
  resolve: (value: any) => void;
  reject: (reason: any) => void;
}> = [];

const processQueue = (error: unknown, token: string | null = null) => {
  failedQueue.forEach(({ resolve, reject }) => {
    if (error) {
      reject(error);
    } else {
      resolve(token);
    }
  });

  failedQueue = [];
};

export async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_URL}${endpoint}`;
  console.log("API Request:", { url, method: options.method || "GET" });

  const res = await fetch(url, {
    ...options,
    credentials: "include",
    headers: {
      "Content-Type": "application/json",

      ...(options.headers || {}),
    },
  });

  if (!res.ok) {
    let errorMessage = `Request failed: ${res.status}`;
    try {
      const errorData = await res.json();
      errorMessage = errorData.message || errorData.error || errorMessage;
    } catch {
      // If we can't parse the error response, use the default message
    }
    throw new Error(errorMessage);
  }

  return res.json();
}

async function refreshAccessToken(): Promise<string> {
  const currentToken = localStorage.getItem("access_token");

  // Don't try to refresh with an empty or invalid token
  if (!currentToken || currentToken.trim() === "") {
    throw new Error("No token to refresh");
  }

  const response = await fetch(`${API_URL}/auth/refresh`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${currentToken}`,
    },
  });

  if (!response.ok) {
    throw new Error("Token refresh failed");
  }

  const data = await response.json();
  const newToken = data.data?.access_token;

  if (!newToken) {
    throw new Error("No access token in refresh response");
  }

  localStorage.setItem("access_token", newToken);
  localStorage.setItem(
    "user",
    JSON.stringify(
      data.data.user || JSON.parse(localStorage.getItem("user") || "{}")
    )
  );
  // Update cookie
  document.cookie = `access_token=${newToken}; path=/; max-age=86400; SameSite=Lax`;

  return newToken;
}

export async function apiRequestWithToken<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const token = localStorage.getItem("access_token");

  const makeRequest = async (authToken: string) => {
    // Don't send empty or invalid tokens
    if (!authToken || authToken.trim() === "") {
      throw new Error("No valid token available");
    }

    const res = await fetch(`${API_URL}${endpoint}`, {
      ...options,
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${authToken}`,
        ...(options.headers || {}),
      },
    });

    return res;
  };

  try {
    // Check if we have a valid token before making request
    if (!token || token.trim() === "") {
      // Redirect to login if no token
      localStorage.removeItem("access_token");
      localStorage.removeItem("user");
      window.location.href = "/login";
      throw new Error("No token available");
    }

    const res = await makeRequest(token);

    // If token is expired, try to refresh
    if (res.status === 401) {
      const responseData = await res.json();

      // Check if it's a token expiration error
      if (
        responseData.error?.includes("expired") ||
        responseData.error?.includes("Invalid token") ||
        responseData.error?.includes("malformed")
      ) {
        if (isRefreshing) {
          // If already refreshing, queue this request
          return new Promise((resolve, reject) => {
            failedQueue.push({ resolve, reject });
          }).then((newToken) => {
            return makeRequest(newToken as string).then((res) => {
              if (!res.ok) {
                throw new Error(`Request failed: ${res.status}`);
              }
              return res.json();
            });
          });
        }

        isRefreshing = true;

        try {
          const newToken = await refreshAccessToken();
          processQueue(null, newToken);

          // Retry original request with new token
          const retryRes = await makeRequest(newToken);
          if (!retryRes.ok) {
            throw new Error(`Request failed: ${retryRes.status}`);
          }
          return retryRes.json();
        } catch (refreshError) {
          processQueue(refreshError, null);

          // If refresh fails, redirect to login
          localStorage.removeItem("access_token");
          localStorage.removeItem("user");
          document.cookie =
            "access_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT";
          window.location.href = "/login";

          throw refreshError;
        } finally {
          isRefreshing = false;
        }
      }
    }

    if (!res.ok) {
      throw new Error(`Request failed: ${res.status}`);
    }

    return res.json();
  } catch (error) {
    // If it's a network error or other non-401 error, just throw it
    throw error;
  }
}
