import React from "react";
import { useRoutes } from "react-router-dom";
import { lazyImport } from "../util/lazyImport";

const { Landing } = lazyImport(() => import("../features/home"), "Landing");
const { Components } = lazyImport(
  () => import("../features/home"),
  "Components"
);
const { Board } = lazyImport(() => import("../features/home"), "Board");

export const AppRoutes = () => {
  const commonRoutes = [
    { path: "/", element: <Landing /> },
    {
      path: "/components",
      element: <Components />,
    },
    {
      path: "/board",
      element: <Board />,
    },
  ];

  const element = useRoutes([...commonRoutes]);

  return <div className="text-t-dark bg-background">{element}</div>;
};
