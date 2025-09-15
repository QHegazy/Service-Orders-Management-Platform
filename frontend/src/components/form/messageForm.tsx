import React from "react";

interface MessageProps {
  message: string;
  type: "error" | "success";
}

const MessageComponent: React.FC<MessageProps> = ({ message, type }) => {
  const baseClasses = "p-4 rounded-md text-sm font-medium";
  const typeClasses =
    type === "error"
      ? "bg-red-50 text-red-800 border border-red-200"
      : "bg-green-50 text-green-800 border border-green-200";

  return <div className={`${baseClasses} ${typeClasses}`}>{message}</div>;
};

const ErrorMessage: React.FC<{ message: string }> = ({ message }) => {
  return <MessageComponent message={message} type="error" />;
};
const SuccessMessage: React.FC<{ message: string }> = ({ message }) => {
  return <MessageComponent message={message} type="success" />;
};

export { ErrorMessage, SuccessMessage };
