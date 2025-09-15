"use client";

export function Footer() {
  return (
    <footer className="bg-white border-t border-gray-200">
      <div className="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
          {/* Company Info */}
          <div className="col-span-1 md:col-span-2">
            <div className="flex items-center">
              <span className="text-xl font-bold text-gray-900">
                <span className="text-indigo-600">Order</span>Management
              </span>
            </div>
            <p className="mt-4 text-gray-600 max-w-md">
              A comprehensive solution for managing orders, customers, and
              products efficiently. Streamline your business operations with our
              powerful tools.
            </p>
          </div>
        </div>
      </div>
    </footer>
  );
}
