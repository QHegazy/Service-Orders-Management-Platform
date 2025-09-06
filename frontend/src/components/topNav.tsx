export function TopNav() {
  return (
    <nav className="bg-white shadow-lg border-b border-gray-200">
      <div className="container mx-auto px-6 py-4">
        <div className="flex justify-between items-center">
          <div className="text-gray-800 text-xl font-bold">
            Order Management
          </div>
          <ul className="flex space-x-8">
            <li>
              <a
                href="/"
                className="text-gray-600 hover:text-blue-600 font-medium transition-colors duration-200"
              >
                Home
              </a>
            </li>
            <li>
              <a
                href="/signup"
                className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors duration-200"
              >
                Sign Up
              </a>
            </li>
            <li>
              <a
                href="/login"
                className="text-gray-600 hover:text-blue-600 font-medium transition-colors duration-200"
              >
                Login
              </a>
            </li>
          </ul>
        </div>
      </div>
    </nav>
  );
}
