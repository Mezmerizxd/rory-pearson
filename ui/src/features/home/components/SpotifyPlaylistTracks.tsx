import React, { useEffect, useState } from "react";
import { Button, Dialog, DialogTitle } from "@components/Elements";
import { useTracks } from "../api/spotify";

interface SpotifyPlaylistTracksProps {
  isOpen: {
    isOpen: boolean;
    playlist: SpotifyPlaylistData;
  };
  onClose: () => void;
  onConfirm: (post: CreatePostData) => void;
}

export const SpotifyPlaylistTracks = ({
  isOpen,
  onClose,
  onConfirm,
}: SpotifyPlaylistTracksProps) => {
  const [error, setError] = useState<string>(null);

  const tracksMut = useTracks();

  useEffect(() => {
    tracksMut.mutateAsync(isOpen.playlist.id).catch((error) => {
      setError(error.message);
    });
  }, [isOpen]);

  return (
    isOpen.playlist !== null && (
      <Dialog className="p-5" isOpen={isOpen.isOpen} onClose={onClose}>
        <DialogTitle as="h2" className="text-lg font-medium">
          {isOpen.playlist.name}
        </DialogTitle>
        <div className="w-96 m-2 space-y-6 overflow-auto py-4">
          {tracksMut.isLoading && <p>Loading...</p>}

          {error && (
            <p className="bg-orange-500/20 text-orange-500 border border-orange-500/20">
              {error}
            </p>
          )}
        </div>
        <div className="flex justify-between">
          <Button
            variant="danger"
            size="sm"
            className="inline-flex justify-center"
            onClick={onClose}
          >
            Close
          </Button>
          <Button
            variant="success"
            size="sm"
            className="inline-flex justify-center"
          >
            Create
          </Button>
        </div>
      </Dialog>
    )
  );
};
