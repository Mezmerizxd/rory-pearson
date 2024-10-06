import { useMutation } from "react-query";
import { MutationConfig } from "../../../util/react-query";
import { useNotificationStore } from "../../../stores/notifications";

// Helper function to determine the hostname
const getHostname = (): string => {
  return window.location.hostname === "localhost"
    ? "http://localhost:3000"
    : window.location.origin;
};

// Types for login and profile responses
export type LoginResponse = { connected: boolean; url: string; state: string };
export type ProfileResponse = {
  user: SpotifyUserData;
  playlists: SpotifyPlaylistData[];
};

// Utility: Get cached login from localStorage
const getCachedLogin = (): LoginResponse | null => {
  const storedLogin = localStorage.getItem("login");
  return storedLogin ? JSON.parse(storedLogin) : null;
};

// Utility: Save login to localStorage
const saveLogin = (login: LoginResponse) => {
  localStorage.setItem("login", JSON.stringify(login));
};

// Utility: Remove cached login
const removeCachedLogin = () => {
  localStorage.removeItem("login");
};

// LOGIN #################################################

// Validate cached login via server
const validateLogin = async (cachedLogin: LoginResponse): Promise<boolean> => {
  try {
    const response = await fetch(
      `${getHostname()}/api/spotify/validate?state=${cachedLogin.state}`,
      {
        method: "GET",
        headers: { "Content-Type": "application/json" },
      }
    );
    return response.ok;
  } catch (error) {
    console.error("Validation error: ", error);
    return false;
  }
};

// Perform login flow
export const login = async (): Promise<LoginResponse | null> => {
  const cachedLogin = getCachedLogin();
  if (cachedLogin) {
    const isValid = await validateLogin(cachedLogin);
    if (isValid) {
      return { connected: true, ...cachedLogin };
    }
    removeCachedLogin();
  }

  const response = await fetch(`${getHostname()}/api/spotify/login`, {
    method: "GET",
    headers: { "Content-Type": "application/json" },
  });

  if (!response.ok) {
    throw new Error(`Error connecting to Spotify: ${response.statusText}`);
  }

  const data: LoginResponse = await response.json();
  saveLogin(data); // Cache login for later use
  return { connected: false, ...data };
};

// React Query hook for login mutation
export const useLogin = (config?: MutationConfig<typeof login>) => {
  const { addNotification } = useNotificationStore();

  return useMutation({
    mutationFn: login,
    ...config,
    onError: (error: any) => {
      addNotification({
        type: "error",
        title: "Spotify",
        message: error?.message || "Unknown error",
      });
    },
    onSuccess: (data) => {
      if (data?.connected) {
        addNotification({
          type: "success",
          title: "Spotify",
          message: "Successfully connected to Spotify",
        });
      }
    },
  });
};

// PROFILE #################################################

// Fetch profile using the cached login
export const fetchProfile = async (): Promise<ProfileResponse | null> => {
  const cachedLogin = getCachedLogin();
  if (!cachedLogin) return null;

  const response = await fetch(
    `${getHostname()}/api/spotify/profile?state=${cachedLogin.state}`,
    {
      method: "GET",
      headers: { "Content-Type": "application/json" },
    }
  );

  if (!response.ok) {
    throw new Error(`Error fetching Spotify profile: ${response.statusText}`);
  }

  return await response.json();
};

// React Query hook for profile fetching
export const useProfile = (config?: MutationConfig<typeof fetchProfile>) => {
  const { addNotification } = useNotificationStore();

  return useMutation({
    mutationFn: fetchProfile,
    ...config,
    onError: (error: any) => {
      addNotification({
        type: "error",
        title: "Spotify",
        message: error?.message || "Failed to fetch profile",
      });
    },
    onSuccess: (data) => {
      if (data) {
        addNotification({
          type: "success",
          title: "Spotify",
          message: "Spotify profile fetched successfully",
        });
      }
    },
  });
};

// TRACKS #####################################################

// Fetch tracks for a given playlist
export const fetchTracks = async (
  playlistId: string
): Promise<SpotifyPlaylistItemData[]> => {
  const cachedLogin = getCachedLogin();
  if (!cachedLogin) return [];

  const response = await fetch(
    `${getHostname()}/api/spotify/tracks?state=${cachedLogin.state}&playlistId=${playlistId}`,
    {
      method: "GET",
      headers: { "Content-Type": "application/json" },
    }
  );

  if (!response.ok) {
    throw new Error(`Error fetching Spotify tracks: ${response.statusText}`);
  }

  return await response.json();
};

// React Query hook for fetching tracks
export const useTracks = (config?: MutationConfig<typeof fetchTracks>) => {
  const { addNotification } = useNotificationStore();

  return useMutation({
    mutationFn: fetchTracks,
    ...config,
    onError: (error: any) => {
      addNotification({
        type: "error",
        title: "Spotify",
        message: error?.message || "Failed to fetch tracks",
      });
    },
    onSuccess: (data) => {
      if (data) {
        addNotification({
          type: "success",
          title: "Spotify",
          message: "Tracks fetched successfully",
        });
      }
    },
  });
};

// DISCONNECT #################################################

// Disconnect from Spotify and clear local cache
export const disconnect = async (): Promise<void> => {
  const cachedLogin = getCachedLogin();
  if (!cachedLogin) return;

  const response = await fetch(
    `${getHostname()}/api/spotify/disconnect?state=${cachedLogin.state}`,
    {
      method: "GET",
      headers: { "Content-Type": "application/json" },
    }
  );

  if (!response.ok) {
    throw new Error(`Error disconnecting from Spotify: ${response.statusText}`);
  }

  removeCachedLogin();
};
