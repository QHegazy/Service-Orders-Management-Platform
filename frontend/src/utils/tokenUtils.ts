export const isValidJWTFormat = (token: string): boolean => {
  if (!token || typeof token !== "string") {
    return false;
  }

  // JWT should have 3 parts separated by dots
  const parts = token.split(".");
  if (parts.length !== 3) {
    return false;
  }

  // Basic check - each part should not be empty
  return parts.every((part) => part.length > 0);
};

export const clearInvalidTokens = (): void => {
  const token = localStorage.getItem("access_token");
  if (token && !isValidJWTFormat(token)) {
    console.warn("Invalid token format detected, clearing storage");
    localStorage.removeItem("access_token");
    localStorage.removeItem("user");
    document.cookie =
      "access_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT";
  }
};
