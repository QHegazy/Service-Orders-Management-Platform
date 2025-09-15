const freeProviders = [
  "gmail.com",
  "yahoo.com",
  "outlook.com",
  "hotmail.com",
  "live.com",
  "aol.com",
  "icloud.com",
  "protonmail.com",
];

/**
 * Check if an email belongs to a business domain
 * @param email string
 * @returns boolean
 */
export function isBusinessMail(email: string): boolean {
  if (!email || !email.includes("@")) {
    return false;
  }

  const domain = email.split("@")[1].toLowerCase();
  return !freeProviders.includes(domain);
}
