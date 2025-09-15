export interface TenantSignUp {
  tenant_name: string;
  domain: string;
  email: string;
}

export interface TenantUpdate {
  tenant_name?: string;
  domain?: string;
  email?: string;
}

export interface Tenant {
  id: string;
  tenant_name: string;
  domain: string;
  email: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}
