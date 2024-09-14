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

type IconProps =
  | { startIcon: React.ReactElement; endIcon?: never }
  | { endIcon: React.ReactElement; startIcon?: never }
  | { endIcon?: undefined; startIcon?: undefined };

export type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: keyof typeof variants;
  size?: keyof typeof sizes;
  isLoading?: boolean;
} & IconProps;

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  (
    {
      type = "button",
      className = "",
      variant = "primary",
      size = "md",
      isLoading = false,
      startIcon,
      endIcon,
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
          {!isLoading && startIcon}
          <span>{props.children}</span> {!isLoading && endIcon}
        </div>
      </button>
    );
  }
);

Button.displayName = "Button";
