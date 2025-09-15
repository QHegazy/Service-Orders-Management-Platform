"use client";

interface Option {
  value: string;
  label: string;
}

interface TicketPropertySelectProps {
  value: string;
  onChange: (e: React.ChangeEvent<HTMLSelectElement>) => void;
  options: Option[];
  disabled?: boolean;
}

export default function TicketPropertySelect({
  value,
  onChange,
  options,
  disabled = false,
}: TicketPropertySelectProps) {
  const displayValue = options.find((opt) => opt.value === value)?.label || value;

  if (disabled) {
    return (
      <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-gray-100 text-gray-800">
        {displayValue}
      </span>
    );
  }

  return (
    <select
      value={value}
      onChange={onChange}
      className="block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
    >
      {options.map((option) => (
        <option key={option.value} value={option.value}>
          {option.label}
        </option>
      ))}
    </select>
  );
}
