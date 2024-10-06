import { useMutation } from "react-query";
import { MutationConfig } from "../../../util/react-query";
import { GetHost } from "../../../util/host";
import { Authenticate } from "./auth";

export const usePlayer = (config?: MutationConfig<typeof testPlayer>) => {
  return useMutation({
    mutationFn: testPlayer,
    ...config,
    onError: (error: any) => {
      console.error("Error getting Spotify player", error);
    },
    onSuccess: (data) => {
      console.log("Successfully fetched Spotify player", data);
    },
  });
};

// APIS ###############################################################

const testPlayer = async (
  auth: Authenticate
): Promise<SpotifyNowPlayingData> => {
  const response = await fetch(
    `${GetHost()}/api/spotify/now-playing?state=${auth.login.state}`,
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
