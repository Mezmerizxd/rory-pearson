import clsx from "clsx";
import React from "react";
import { Textarea as MTTextarea } from "@material-tailwind/react";

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

export type TextAreaProps = {
  label: string;
  onChange?: (e: React.ChangeEvent<HTMLTextAreaElement>) => void;
  className?: string;
};

export const TextArea = ({
  label,
  onChange,
  className = "",
}: TextAreaProps) => {
  return (
    <>
      <MTTextarea
        className="font-fira text-t-light"
        variant="outlined"
        color="white"
        label="Outlined"
        onChange={onChange}
        onPointerEnterCapture={undefined}
        onPointerLeaveCapture={undefined}
      />
    </>
  );
};
