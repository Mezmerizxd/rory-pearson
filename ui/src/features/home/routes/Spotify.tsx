import React, { useEffect, useState } from "react";
import { Button } from "@components/Elements";
import {
  useLogin,
  useProfile,
  disconnect as disconnectSpotify,
  useTracks,
} from "../api/spotify";
import HomeLayout from "@components/Layout/HomeLayout"; // Use non-lazy import for simplicity
import { SpotifyPlaylistTracks } from "../components/SpotifyPlaylistTracks";

export const Spotify = () => {
  const loginMut = useLogin();
  const profileMut = useProfile();
  const [isConnected, setIsConnected] = useState(false); // Track connection status
  const [spotifyPlaylistTracksDialog, setSpotifyPlaylistTracksDialog] =
    React.useState<{
      isOpen: boolean;
      playlist: SpotifyPlaylistData;
    }>({ isOpen: false, playlist: null });

  useEffect(() => {
    const checkLoginStatus = async () => {
      if (!isConnected) {
        try {
          const loginData = await loginMut.mutateAsync(null);
          if (loginData?.connected) {
            setIsConnected(true);
            await profileMut.mutateAsync(null); // Only fetch profile if connected
          }
        } catch (error) {
          console.error("Error checking login status:", error);
        }
      }
    };

    checkLoginStatus(); // Only fire this once on mount
  }, []);

  const handleConnect = async () => {
    if (loginMut.data?.url) {
      window.location.href = loginMut.data.url; // Redirect to Spotify's login
    }
  };

  const handleDisconnect = async () => {
    try {
      await disconnectSpotify();
      localStorage.removeItem("login");
      loginMut.reset();
      profileMut.reset();
      setIsConnected(false); // Reset connection status
    } catch (error) {
      console.error("Error during disconnect:", error);
    }
  };

  return (
    <HomeLayout title="Spotify">
      <SpotifyPlaylistTracks
        key={"spotify-playlist-tracks"}
        isOpen={spotifyPlaylistTracksDialog}
        onClose={() =>
          setSpotifyPlaylistTracksDialog({ isOpen: false, playlist: null })
        }
        onConfirm={() => {}}
      />

      {isConnected ? (
        <div className="p-4 sm:p-6 lg:p-10 space-y-5">
          <p className="bg-green-500/20 text-green-500 border border-green-500/20 text-center p-2 rounded">
            Connected
          </p>

          {profileMut.data?.user && (
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 p-4 bg-background-dark rounded-md">
              <UserInfoSection user={profileMut.data.user} />
              <div className="flex justify-center items-center">
                <Button variant="danger" onClick={handleDisconnect}>
                  Disconnect
                </Button>
              </div>
            </div>
          )}

          {profileMut.data?.playlists && (
            <div className="p-4 space-y-2">
              <h2 className="text-xl font-bold text-white">Playlists</h2>
              <div className="grid grid-cols-2 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                {profileMut.data.playlists.map((playlist) => (
                  <PlaylistCard
                    key={playlist.id}
                    playlist={playlist}
                    setDialog={setSpotifyPlaylistTracksDialog}
                  />
                ))}
              </div>
            </div>
          )}
        </div>
      ) : (
        <div className="w-screen h-screen flex justify-center items-center">
          <Button onClick={handleConnect}>Connect to Spotify</Button>
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

const PlaylistCard = ({
  playlist,
  setDialog,
}: {
  playlist: SpotifyPlaylistData;
  setDialog: (dialog: {
    isOpen: boolean;
    playlist: SpotifyPlaylistData;
  }) => void;
}) => {
  const {
    images = [],
    name,
    external_urls,
    description,
    owner,
    tracks,
  } = playlist;

  const handleSelectPlaylist = async () => {
    try {
      // await tracksMut.mutateAsync(playlistId);
      setDialog({ isOpen: true, playlist });
    } catch (error) {
      console.error("Error selecting playlist:", error);
    }
  };

  return (
    <div
      className="flex flex-col justify-between items-start p-4 bg-background-dark rounded-lg transition-transform duration-200 hover:scale-105"
      onClick={() => handleSelectPlaylist()}
    >
      {images !== null && images.length > 0 ? (
        <img
          src={images[0].url}
          alt={name}
          className="w-full h-64 object-cover rounded-lg"
        />
      ) : (
        <div className="w-full h-32 bg-gray-700 rounded-lg flex items-center justify-center text-gray-300">
          No Image Available
        </div>
      )}
      <h1 className="mt-2">
        <a className="text-t-light font-bold hover:underline">
          {name || "Unnamed Playlist"}
        </a>
      </h1>
      <p className="mt-2 text-t-light text-left">
        {description || "No Description Available"}
      </p>
      <p className="mt-1 text-sm text-gray-400">
        Owner: {owner.display_name || "Unknown"}
      </p>
      <p className="mt-1 text-sm text-gray-400">
        Tracks: {tracks.total || "0"}
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
