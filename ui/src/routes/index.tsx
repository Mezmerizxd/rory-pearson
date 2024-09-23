import React from "react";
import { useRoutes } from "react-router-dom";
import { lazyImport } from "../util/lazyImport";

const { Landing } = lazyImport(() => import("../features/home"), "Landing");
const { Components } = lazyImport(
  () => import("../features/home"),
  "Components"
);
const { Board } = lazyImport(() => import("../features/home"), "Board");
const { BackgroundRemover } = lazyImport(
  () => import("../features/home"),
  "BackgroundRemover"
);
const { ImageToIcon } = lazyImport(
  () => import("../features/home"),
  "ImageToIcon"
);

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
    {
      path: "/background-remover",
      element: <BackgroundRemover />,
    },
    {
      path: "/image-to-icon",
      element: <ImageToIcon />,
    },
  ];

  const element = useRoutes([...commonRoutes]);

  return <div className="text-t-dark bg-background">{element}</div>;
};
