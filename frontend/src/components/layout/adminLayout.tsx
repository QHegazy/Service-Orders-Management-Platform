import { TenantForm } from "../form/tenantForm";

export function AdminLayout() {

  return <div>Admin Layout

               
            <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
              <h3 className="text-lg font-medium text-gray-900 mb-4">
                Admin Tools
              </h3>
              <TenantForm />
            </div>
         

  </div>;
}
