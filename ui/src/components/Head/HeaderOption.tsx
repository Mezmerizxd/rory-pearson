import React from "react";
import { Button } from "@components/Elements/Button";
import { useNavigate } from "react-router-dom";

export const HeaderOption = ({
  name,
  href,
  onClick,
  customUrl,
}: {
  name: string;
  href?: string;
  onClick?: () => void;
  customUrl?: boolean;
}) => {
  const navigate = useNavigate();

  return (
    <Button
      size="sm"
      variant="primary"
      onClick={() => {
        if (customUrl) {
          window.open(href, "_blank");
          return;
        }
        navigate(href);
        onClick();
      }}
    >
      {name}
    </Button>
  );
};
