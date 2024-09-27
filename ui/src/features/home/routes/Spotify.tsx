import React, { useEffect } from "react";
import { Button } from "@components/Elements";
import { useLogin, useProfile } from "../api/spotify";
import clsx from "clsx";

const HomeLayout = React.lazy(() => import("@components/Layout/HomeLayout"));

export const Spotify = () => {
  const loginMut = useLogin();
  const profileMut = useProfile();

  useEffect(() => {
    setTimeout(async () => {
      const data = await loginMut.mutateAsync(null);

      if (data.connected) {
        profileMut.mutateAsync(null);
      }
    });
  }, []);

  // Handle the Spotify login connection
  async function connect() {
    if (loginMut.data?.url) {
      window.location.href = loginMut.data.url;
    }
  }

  // Connected state (profile data)
  return (
    <HomeLayout title="Spotify">
      {loginMut.data?.connected ? (
        <div className="p-10">
          <p>Connected</p>
          {/* <pre>{JSON.stringify(profileMut.data, null, 2)}</pre> */}

          {/* I want a grid with 2 columns and a max height of 800px with scrolling Y */}
          <div className="grid grid-cols-2 max-h-800px">
            <div
              className={clsx(
                "p-2 text-left",
                "overflow-y-auto",
                "border",
                "border-gray-300"
              )}
            >
              <pre>{JSON.stringify(profileMut.data?.user, null, 2)}</pre>
            </div>
            <div
              className={clsx(
                "p-2",
                "overflow-y-auto",
                "border",
                "border-gray-300"
              )}
            >
              {profileMut.data?.playlists?.items?.map((playlist) => (
                <div key={playlist.id} className="p-2">
                  <p className="text-center">{playlist.name}</p>
                </div>
              ))}
            </div>
          </div>
        </div>
      ) : (
        <div className="w-screen h-screen flex justify-center items-center">
          <Button onClick={connect}>Connect to Spotify</Button>
        </div>
      )}
    </HomeLayout>
  );
};
