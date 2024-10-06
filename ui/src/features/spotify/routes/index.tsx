import * as React from "react";
import { lazyImport } from "../../../util/lazyImport";
import { DashLayout } from "@components/Layout";
import { Label, Spinner } from "@components/Elements";
import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "../api";
import { FaHome } from "react-icons/fa";
import { RiPlayListFill } from "react-icons/ri";

const { Dashboard } = lazyImport(() => import("./Dashboard"), "Dashboard");
const { Playlists } = lazyImport(() => import("./Playlists"), "Playlists");

const Component = () => {
  const auth = useAuth();

  React.useEffect(() => {
    if (auth.data === undefined || auth.data === null) {
      auth.mutate(null);
    }
  }, [auth.data]);

  return (
    <DashLayout
      navigation={[
        { name: "Dashboard", to: "/spotify", icon: <FaHome /> },
        {
          name: "Playlists",
          to: "/spotify/playlists",
          icon: <RiPlayListFill />,
        },
      ]}
      userNavigation={[
        {
          name: "Profile",
          to: ".",
          onClick: () => {
            window.open(auth.data?.user?.external_urls.spotify, "_blank");
          },
        },
      ]}
      userElement={
        <div>
          {auth.data && auth.data.user && (
            <Label size="sm">{auth.data.user.display_name}</Label>
          )}
        </div>
      }
    >
      <React.Suspense
        fallback={
          <div className="h-full w-full flex items-center justify-center">
            <Spinner size="xl" />
          </div>
        }
      >
        <Outlet />
      </React.Suspense>
    </DashLayout>
  );
};

export const Routes = {
  path: "/spotify",
  element: <Component />,
  children: [
    { index: true, element: <Dashboard /> },
    { path: "playlists", element: <Playlists /> },
    { path: "*", element: <Navigate to="." /> },
  ],
};
