import clsx from "clsx";
import * as React from "react";

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

export type LabelProps = React.HTMLAttributes<HTMLDivElement> & {
  variant?: keyof typeof variants;
  size?: keyof typeof sizes;
} & IconProps;

export const Label = React.forwardRef<HTMLDivElement, LabelProps>(
  (
    {
      className = "",
      variant = "primary",
      size = "md",
      startIcon,
      endIcon,
      ...props
    },
    ref
  ) => {
    return (
      <div
        ref={ref}
        className={clsx(
          "inline-flex justify-center items-center rounded-md font-medium",
          variants[variant],
          sizes[size],
          className
        )}
        {...props}
      >
        <div className="flex justify-between items-center space-x-2">
          {startIcon}
          <span>{props.children}</span>
          {endIcon}
        </div>
      </div>
    );
  }
);

Label.displayName = "Label";
