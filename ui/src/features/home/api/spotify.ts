import { useMutation } from "react-query";
import { MutationConfig } from "../../../util/react-query";
import { useNotificationStore } from "../../../stores/notifications";

// Helper function to determine the hostname
function getHostname() {
  if (window.location.hostname === "localhost") {
    return "http://localhost:3000";
  }
  return window.location.origin;
}

// Types for login and profile responses
export type LoginResponse = { connected: boolean; url: string; state: string };
export type ProfileResponse = { user: any; playlists: any };

// LOGIN #################################################

// Get cached login from localStorage
export const GetLogin = (): LoginResponse | null => {
  const storedLogin = localStorage.getItem("login");
  return storedLogin ? JSON.parse(storedLogin) : null;
};

// Login mutation: Handles the login process with validation
export const Login = async (): Promise<LoginResponse | null> => {
  const cachedLogin = GetLogin();
  if (cachedLogin) {
    const validateResponse = await fetch(
      `${getHostname()}/api/spotify/validate?state=${cachedLogin.state}`,
      {
        method: "GET",
        headers: { "Content-Type": "application/json" },
      }
    );

    // If cached login is invalid, remove it and reset login
    if (!validateResponse.ok) {
      localStorage.removeItem("login");
      return null;
    }

    return {
      connected: true,
      ...cachedLogin,
    }; // Return cached login if valid
  }

  const response = await fetch(`${getHostname()}/api/spotify/login`, {
    method: "GET",
    headers: { "Content-Type": "application/json" },
  });

  if (!response.ok) {
    throw new Error(`Error connecting to Spotify: ${response.statusText}`);
  }

  const data: LoginResponse = await response.json();
  localStorage.setItem("login", JSON.stringify(data)); // Cache login
  return {
    connected: false,
    ...data,
  };
};

// React Query hook for login mutation
export const useLogin = (config?: MutationConfig<typeof Login>) => {
  const { addNotification } = useNotificationStore();

  return useMutation({
    onError: (error: any) => {
      addNotification({
        type: "error",
        title: "Spotify",
        message: "You must be logged in to connect to Spotify",
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
    ...config,
    mutationFn: Login,
  });
};

// PROFILE #################################################

// Profile fetching logic
export const Profile = async (): Promise<ProfileResponse | null> => {
  const login = GetLogin();
  if (!login) return null; // Ensure login exists before fetching profile

  const response = await fetch(
    `${getHostname()}/api/spotify/profile?state=${login.state}`,
    {
      method: "GET",
      headers: { "Content-Type": "application/json" },
    }
  );

  if (!response.ok) {
    throw new Error(`Error getting Spotify profile: ${response.statusText}`);
  }

  return await response.json();
};

// React Query hook for profile mutation
export const useProfile = (config?: MutationConfig<typeof Profile>) => {
  const { addNotification } = useNotificationStore();

  return useMutation({
    onError: (error: any) => {
      addNotification({
        type: "error",
        title: "Spotify",
        message: error.message,
      });
    },
    onSuccess: (data) => {
      if (data) {
        addNotification({
          type: "success",
          title: "Spotify",
          message: "Successfully fetched Spotify profile",
        });
      }
    },
    ...config,
    mutationFn: Profile,
  });
};
