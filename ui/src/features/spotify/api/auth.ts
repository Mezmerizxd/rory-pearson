import { useMutation } from "react-query";
import { MutationConfig } from "../../../util/react-query";
import { GetHost } from "../../../util/host";

type LoginResponse = {
  state: string;
  url: string;
};

export type Authenticate = {
  isValid: boolean;
  login: LoginResponse;
  user?: SpotifyUserData;
};

export const useAuth = (config?: MutationConfig<typeof authenticate>) => {
  return useMutation({
    mutationFn: authenticate,
    ...config,
    onError: (error: any) => {
      console.info("Error authenticating with Spotify", error);
    },
    onSuccess: (data) => {
      console.info("Successfully authenticated with Spotify", data);
    },
  });
};

export const openSpotifyLogin = (url: string) => {
  window.open(url, "_self");
};

// APIS ###############################################################

async function authenticate(): Promise<Authenticate> {
  try {
    const login = getLoginState();

    if (login) {
      const isValid = await validate(login.state);

      if (isValid) {
        const user = await profile();
        return { isValid: isValid, login, user };
      }

      setLoginState(null);
    }

    const newLogin = await llllooooggggiiiinnn();
    setLoginState(newLogin);

    console.log("newLogin", newLogin);
    return { isValid: false, login: newLogin };
  } catch (error) {
    console.error("Error during authentication:", error);
    throw error;
  }
}

async function llllooooggggiiiinnn(): Promise<LoginResponse> {
  try {
    const response = await fetch(`${GetHost()}/api/spotify/login`, {
      method: "GET",
      headers: { "Content-Type": "application/json" },
    });

    if (!response.ok) {
      throw new Error(`Error connecting to Spotify: ${response.statusText}`);
    }

    const json = await response.json();
    if (!json) {
      throw new Error("Error connecting to Spotify: No response");
    }

    return json;
  } catch (error) {
    console.error("Error in login process:", error);
    throw error;
  }
}

async function validate(state: string): Promise<boolean> {
  try {
    const response = await fetch(
      `${GetHost()}/api/spotify/validate?state=${state}`,
      {
        method: "GET",
        headers: { "Content-Type": "application/json" },
      }
    );

    if (!response.ok) {
      console.warn(
        `Validation failed for state ${state}: ${response.statusText}`
      );
    }

    return response.ok;
  } catch (error) {
    console.error("Error during validation:", error);
    return false;
  }
}

async function profile(): Promise<SpotifyUserData> {
  try {
    const response = await fetch(
      `${GetHost()}/api/spotify/profile?state=${state}`,
      {
        method: "GET",
        headers: { "Content-Type": "application/json" },
      }
    );

    if (!response.ok) {
      throw new Error(`Error getting Spotify profile: ${response.statusText}`);
    }

    const json = await response.json();
    if (!json) {
      throw new Error("Error getting Spotify profile: No response");
    }

    return json;
  } catch (error) {
    console.error("Error retrieving profile:", error);
    throw error;
  }
}

// LOCAL STORAGE STUFF #################################################

const localStorageKey = "spotify_login";
let state: string;

export function getLoginState(): LoginResponse | null {
  const local = localStorage.getItem(localStorageKey);
  if (!local || local === "null") {
    localStorage.removeItem(localStorageKey);
    return null;
  }

  const data = JSON.parse(local);

  // Check none of the objects are null
  if (data) {
    Object.keys(data).forEach((key) => {
      if (data[key] === null) {
        return null;
      }
    });
  }

  state = data.state;

  return data as LoginResponse;
}

function setLoginState(login: LoginResponse) {
  localStorage.setItem(localStorageKey, JSON.stringify(login));
}
