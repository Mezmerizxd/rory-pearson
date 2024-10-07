import { useMutation } from "react-query";
import { MutationConfig } from "../../../util/react-query";
import { GetHost } from "../../../util/host";
import { Authenticate } from "./auth";

export const usePlaylists = (config?: MutationConfig<typeof playlists>) => {
  return useMutation({
    mutationFn: playlists,
    ...config,
    onError: (error: any) => {
      console.error("Error getting Spotify playlists", error);
    },
    onSuccess: (data) => {
      console.log("Successfully fetched Spotify playlists", data);
    },
  });
};

export const usePlaylistToYoutube = (
  config?: MutationConfig<typeof playlistToYoutube>
) => {
  return useMutation({
    mutationFn: playlistToYoutube,
    ...config,
    onError: (error: any) => {
      console.error("Error converting Spotify playlist to YouTube", error);
    },
    onSuccess: (data) => {
      console.log("Successfully converted Spotify playlist to YouTube", data);
    },
  });
};

// APIS ###############################################################

const playlists = async (auth: Authenticate): Promise<SpotifyPlaylistData> => {
  const response = await fetch(
    `${GetHost()}/api/spotify/playlists?state=${auth.login.state}`,
    {
      method: "GET",
      headers: { "Content-Type": "application/json" },
    }
  );

  if (!response.ok) {
    throw new Error(`Error getting Spotify playlists: ${response.statusText}`);
  }

  const json = await response.json();
  if (json.error) {
    throw new Error(json.error);
  }

  return json;
};

const playlistToYoutube = async ({
  auth,
  playlistId,
}: {
  auth: Authenticate;
  playlistId: string;
}): Promise<YoutubeVideoData[]> => {
  const response = await fetch(
    `${GetHost()}/api/spotify/playlist-to-youtube?state=${auth.login.state}`,
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ playlistId }),
    }
  );

  if (!response.ok) {
    throw new Error(
      `Error converting Spotify playlist to YouTube: ${response.statusText}`
    );
  }

  const json = await response.json();
  if (json.error) {
    throw new Error(json.error);
  }

  return json;
};
