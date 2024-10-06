import React, { useEffect, useState } from "react";
import { ContentLayout } from "../../../components/Layout";
import { Button, Spinner } from "@components/Elements";

import { useAuth, openSpotifyLogin, usePlaylists, usePlayer } from "../api";
import { Player } from "../components/Player";
import clsx from "clsx";

export const Dashboard = () => {
  const auth = useAuth();
  const playlists = usePlaylists();
  const player = usePlayer();

  useEffect(() => {
    connect();
  }, []);

  const connect = async () => {
    const data = await auth.mutateAsync(null);
    console.log(data);
    if (data && data.isValid) {
      await playlists.mutateAsync(data);
      await player.mutateAsync(data);
    }
  };

  return (
    <ContentLayout title="Dashboard">
      {auth.isLoading && (
        <div
          className={clsx(
            "flex flex-col items-center justify-center space-y-2"
          )}
        >
          <Spinner />
          <p>Loading...</p>
        </div>
      )}

      {auth.data && auth.data.isValid && (
        <div className="space-y-5">
          <h1>Hello, {auth.data?.user?.display_name}</h1>
          <p className={clsx("text-green-500")}>
            You are logged into your spotify account.
          </p>

          {playlists.data && <p>You have {playlists.data.total} playlist/s</p>}

          <div className="space-y-2">
            {player.data && <Player fakeTime data={player.data} size="sm" />}
          </div>
        </div>
      )}

      {auth.data && !auth.data.isValid && (
        <div className="space-y-2">
          <p className={clsx("text-red-500")}>
            You are not logged into your spotify account
          </p>
          <Button
            size="sm"
            onClick={() => openSpotifyLogin(auth.data.login.url)}
          >
            Login into Spotify
          </Button>
        </div>
      )}
    </ContentLayout>
  );
};
