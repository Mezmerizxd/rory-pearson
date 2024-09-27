import React, { useEffect } from "react";
import { Button } from "@components/Elements";
import { useLogin, useProfile, Disconnect } from "../api/spotify";
import clsx from "clsx";

const HomeLayout = React.lazy(() => import("@components/Layout/HomeLayout"));

export const Spotify = () => {
  const loginMut = useLogin();
  const profileMut = useProfile();

  useEffect(() => {
    const checkLoginStatus = async () => {
      const data = await loginMut.mutateAsync(null);
      if (data.connected) {
        profileMut.mutateAsync(null);
      }
    };
    const timeoutId = setTimeout(checkLoginStatus, 1000);
    return () => clearTimeout(timeoutId);
  }, []);

  // Handle the Spotify login connection
  const connect = async () => {
    if (loginMut.data?.url) {
      window.location.href = loginMut.data.url;
    }
  };

  const disconnect = async () => {
    localStorage.removeItem("login");
    await Disconnect();
    loginMut.reset();
    profileMut.reset();
  };

  return (
    <HomeLayout title="Spotify">
      {loginMut.data?.connected ? (
        <div className="p-4 sm:p-6 lg:p-10 space-y-5">
          <p className="bg-green-500/20 text-green-500 border border-green-500/20 text-center p-2 rounded">
            Connected
          </p>

          {profileMut.data?.user && (
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 p-4 bg-background-dark rounded-md">
              <UserInfoSection user={profileMut.data.user} />
              <div className="flex justify-center items-center">
                <Button variant="danger" onClick={disconnect}>
                  Disconnect
                </Button>
              </div>
            </div>
          )}

          {profileMut.data?.playlists && (
            <div className="p-4 space-y-2">
              <h2 className="text-xl font-bold text-white">Playlists</h2>
              <div className="grid grid-cols-2 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                {profileMut.data.playlists.items.map((playlist) => (
                  <PlaylistCard key={playlist.id} playlist={playlist} />
                ))}
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

const UserInfoSection = ({ user }) => (
  <div className="flex flex-col space-y-2">
    <InfoRow label="Display Name" value={user.display_name} />
    <InfoRow label="Email" value={user.email} />
    <InfoRow label="Country" value={user.country} />
    <InfoRow label="Followers" value={user.followers.total} />
    <InfoRow label="Product" value={user.product} />
    <InfoRow label="Birthdate" value={user.birthdate} />
  </div>
);

const PlaylistCard = ({ playlist }) => {
  const images = playlist.images || [];
  return (
    <div className="flex flex-col justify-between items-start p-4 bg-background-dark rounded-lg transition-transform duration-200 hover:scale-105">
      {images.length > 0 ? (
        <img
          src={images[0].url}
          alt={playlist.name}
          className="w-full h-64 object-cover rounded-lg"
        />
      ) : (
        <div className="w-full h-32 bg-gray-700 rounded-lg flex items-center justify-center text-gray-300">
          No Image Available
        </div>
      )}
      <h1 className="mt-2">
        <a
          href={playlist.external_urls.spotify}
          className="text-t-light font-bold hover:underline"
        >
          {playlist.name || "Unnamed Playlist"}
        </a>
      </h1>
      <p className="mt-2 text-t-light text-left">
        {playlist.description || "No Description Available"}
      </p>
      <p className="mt-1 text-sm text-gray-400">
        Owner: {playlist.owner.display_name || "Unknown"}
      </p>
      <p className="mt-1 text-sm text-gray-400">
        Tracks: {playlist.tracks.total || "0"}
      </p>
    </div>
  );
};

const InfoRow = ({ label, value }) => (
  <div className="flex whitespace-nowrap">
    <p className="font-bold w-36 text-right text-t-light">{label}:</p>
    <p className="ml-2 text-t-dark">{value || "Not Available"}</p>
  </div>
);
