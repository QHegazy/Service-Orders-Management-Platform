export function Footer() {
  return (
    <footer className="bg-gray-50 border-t border-gray-200 py-8 mt-16">
      <div className="container mx-auto px-6 text-center">
        <p className="text-gray-600">
          &copy; {new Date().getFullYear()} Order Management System. All rights
          reserved.
        </p>
      </div>
    </footer>
  );
}
