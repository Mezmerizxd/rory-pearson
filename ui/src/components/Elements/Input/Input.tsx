import clsx from "clsx";
import React from "react";
import { Input as MTInput } from "@material-tailwind/react";

const sizes = {
  xs: "h-2 w-2",
  sm: "h-4 w-4",
  md: "h-6 w-6",
  lg: "h-8 w-8",
  xl: "h-10 w-10",
};

const variants = {
  primary: "text-accent-light",
};

export type InputProps = {
  label: string;
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
  className?: string;
};

export const Input = ({ label, onChange, className = "" }: InputProps) => {
  return (
    <>
      <MTInput
        className="font-fira text-t-light"
        color="white"
        variant="outlined"
        label={label}
        onChange={onChange}
        onPointerEnterCapture={undefined}
        onPointerLeaveCapture={undefined}
        crossOrigin={undefined}
      />
    </>
  );
};
