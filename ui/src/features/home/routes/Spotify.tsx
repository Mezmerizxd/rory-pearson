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
        <div className="p-10 space-y-5">
          <p>Connected</p>

          {profileMut.data?.user && (
            <div
              className={clsx(
                "grid",
                "grid-cols-3",
                "gap-4",
                "p-4 bg-background-dark"
              )}
            >
              <div id="user_image">
                <img
                  src={profileMut.data.user.images[0].url}
                  alt={profileMut.data.user.display_name}
                  className="rounded-full"
                />
              </div>

              <div
                id="user_info"
                className="flex flex-col justify-center items-start space-y-2"
              >
                <InfoRow
                  label="Display Name"
                  value={profileMut.data.user.display_name}
                />
                <InfoRow label="Email" value={profileMut.data.user.email} />
                <InfoRow label="Country" value={profileMut.data.user.country} />
              </div>

              <div
                id="user_data"
                className="flex flex-col justify-center items-start space-y-2"
              >
                <InfoRow
                  label="Followers"
                  value={profileMut.data.user.followers.total}
                />
                <InfoRow label="Product" value={profileMut.data.user.product} />
                <InfoRow
                  label="Birthdate"
                  value={profileMut.data.user.birthdate}
                />
              </div>
            </div>
          )}

          {profileMut.data?.playlists && (
            <div className="p-4">
              <h2 className="text-xl font-bold text-white">Playlists</h2>
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                {profileMut.data.playlists.items.map((playlist) => {
                  const images = playlist.images || []; // Fallback to an empty array if images is null
                  return (
                    <div
                      key={playlist.id}
                      className="flex flex-col justify-center items-start p-4 bg-background-dark rounded-lg transition-transform duration-200 hover:scale-105"
                    >
                      {images.length > 0 ? (
                        <img
                          src={images[0].url}
                          alt={playlist.name}
                          className="w-full h-32 object-cover rounded-lg"
                        />
                      ) : (
                        <div className="w-full h-32 bg-gray-700 rounded-lg flex items-center justify-center text-gray-300">
                          No Image Available
                        </div>
                      )}
                      <h1 className="mt-2">
                        <a
                          href={playlist.external_urls.spotify}
                          className="text-t-light hover:underline"
                        >
                          {playlist.name || "Unnamed Playlist"}
                        </a>
                      </h1>
                      <p className="mt-2 text-t-light">
                        {playlist.description || "No Description Available"}
                      </p>
                      <p className="mt-2 text-sm text-gray-400">
                        Owner: {playlist.owner.display_name || "Unknown"}
                      </p>
                      <p className="mt-2 text-sm text-gray-400">
                        Tracks: {playlist.tracks.total || "0"}
                      </p>
                    </div>
                  );
                })}
              </div>
            </div>
          )}
        </div>
      ) : (
        <div className="w-screen h-screen flex justify-center items-center">
          <Button onClick={connect}>Connect to Spotify</Button>
        </div>
      )}
    </HomeLayout>
  );
};

const InfoRow = ({ label, value }) => (
  <div className="flex whitespace-nowrap">
    <p className="font-bold w-36 text-right text-t-light">{label}:</p>
    <p className="ml-2 text-t-dark">{value || "Not Available"}</p>
  </div>
);
