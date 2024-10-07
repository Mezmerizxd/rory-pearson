/**
 * @param tracks
 * @returns Array [{Song Name - Artist Name, Artist Name, etc...}]
 */
export function GetTrackNamesFromSpotifyTracks(
  tracks: SpotifyTrackData[]
): string[] {
  const trackNames: string[] = [];
  tracks.forEach((item) => {
    const artists = item.artists.map((artist) => artist.name).join(", ");
    trackNames.push(`${item.name} - ${artists}`);
  });

  return trackNames;
}
