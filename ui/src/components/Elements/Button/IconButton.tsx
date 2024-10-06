import clsx from "clsx";
import * as React from "react";

import { Spinner } from "../Spinner";

const variants = {
  primary: "bg-accent-dark/20 text-t-light border border-accent-dark/20",
  secondary: "bg-accent-light/20 text-t-light border border-accent-light/20",
  danger: "bg-red-500/20 text-red-500 border border-red-500/20",
  success: "bg-green-500/20 text-green-500 border border-green-500/20",
  warning: "bg-orange-500/20 text-orange-500 border border-orange-500/20",
  info: "bg-info text-t-light border border-info",
};

const sizes = {
  xs: "py-1 px-2 text-xs",
  sm: "py-2 px-4 text-sm",
  md: "py-4 px-6 text-md",
  lg: "py-4 px-6 text-lg",
  xl: "py-4 px-6 text-xl",
};

export type IconButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  icon: React.ReactElement;
  variant?: keyof typeof variants;
  size?: keyof typeof sizes;
  isLoading?: boolean;
};

export const IconButton = React.forwardRef<HTMLButtonElement, IconButtonProps>(
  (
    {
      type = "button",
      className = "",
      icon,
      variant = "primary",
      size = "md",
      isLoading = false,
      ...props
    },
    ref
  ) => {
    return (
      <button
        ref={ref}
        type={type}
        className={clsx(
          "flex justify-center items-center disabled:opacity-70 disabled:cursor-not-allowed rounded-md font-medium focus:outline-none hover:opacity-80 duration-150",
          variants[variant],
          sizes[size],
          className
        )}
        {...props}
      >
        <div className="flex justify-between items-center space-x-2">
          {isLoading && (
            <Spinner size={size as any} className=" text-current" />
          )}
          {React.cloneElement(icon, {
            className: clsx(icon.props.className, isLoading && "hidden"),
          })}
        </div>
      </button>
    );
  }
);

IconButton.displayName = "IconButton";
