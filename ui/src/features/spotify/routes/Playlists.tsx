import React, { useEffect, useState } from "react";
import { ContentLayout } from "../../../components/Layout";
import { Button, Spinner } from "@components/Elements";

import { useAuth, openSpotifyLogin, usePlaylists, usePlayer } from "../api";
import { PlaylistThumbnail } from "../components";
import clsx from "clsx";

export const Playlists = () => {
  const auth = useAuth();
  const playlists = usePlaylists();

  useEffect(() => {
    connect();
  }, []);

  const connect = async () => {
    const data = await auth.mutateAsync(null);

    if (data && !data.isValid) {
      window.location.href = "/spotify";
    } else {
      await playlists.mutateAsync(data);
    }
  };

  return (
    <ContentLayout title="Playlists">
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

      {auth.data && playlists.data && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {playlists.data.items.map((item) => (
            <PlaylistThumbnail
              key={item.id}
              data={item}
              options={[
                {
                  name: "View Playlist",
                  func: () => {
                    window.open(item.external_urls.spotify, "_blank");
                  },
                },
              ]}
            />
          ))}
        </div>
      )}
    </ContentLayout>
  );
};
