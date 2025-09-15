interface UserSignUp {
  username: string;
  email: string;
  password: string;
}

interface CustomerSignUp {
  first_name: string;
  last_name: string;
  username: string;
  email?: string;
  password: string;
}

interface Login {
  Username: string;
  Password: string;
}

interface UserClaims {
  id: string;
  username: string;
  role: string;
  tenants: string[];
}

interface User {
  id: string;
  username: string;
  email: string;
  role: string;
  is_active: boolean;
  created_at: string;
}

export type { User, CustomerSignUp, Login, UserSignUp, UserClaims };
