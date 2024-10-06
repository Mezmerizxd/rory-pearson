import clsx from "clsx";
import Marquee from "react-fast-marquee";
import * as React from "react";

const variants = {
  playing: "text-green-500",
  notPlaying: "text-red-500",
};

const sizes = {
  xs: "text-xs",
  sm: "text-sm",
  md: "text-md",
  lg: "text-lg",
};

export type PlayerProps = {
  data: SpotifyNowPlayingData;
  fakeTime?: boolean;
  variant?: keyof typeof variants;
  size?: keyof typeof sizes;
  className?: string;
};

export const Player: React.FC<PlayerProps> = ({
  data,
  fakeTime = false,
  variant = data.item ? "playing" : "notPlaying",
  size = "md",
  className = "",
}) => {
  const [fakeProgress, setFakeProgress] = React.useState(data.progress_ms);

  React.useEffect(() => {
    if (fakeTime) {
      const interval = setInterval(() => {
        if (fakeProgress >= data.item.duration_ms) {
          clearInterval(interval);
          return;
        }

        setFakeProgress((prev) => {
          if (prev >= data.item.duration_ms) return 0;
          return prev + 1000;
        });
      }, 1000);
      return () => clearInterval(interval);
    }
  }, []);

  const imageSize = () => {
    if (size === "xs") return 60;
    if (size === "sm") return 80;
    if (size === "md") return 100;
    if (size === "lg") return 120;
  };

  return data && data.item ? (
    <div className="bg-background-dark border border-background-light rounded-md w-fit">
      <h1 className={clsx("text-accent-light font-bold px-2 pt-1")}>
        Now Playing
      </h1>
      <div
        className={clsx("flex p-2", sizes[size], variants[variant], className)}
      >
        <div>
          <img
            className="rounded-md object-cover flex-shrink-0"
            style={{ width: imageSize(), height: imageSize() }}
            src={data.item.album.images[0].url}
            alt={`${data.item.name} by ${data.item.artists.map((artist) => artist.name).join(", ")}`}
          />
        </div>

        <div className="grid grid-rows-2 justify-normal items-center px-2">
          <div>
            <Marquee gradient={false} speed={100}>
              <h1 className="text-accent-light font-bold">{data.item.name}</h1>
            </Marquee>
            <Marquee gradient={false} speed={100}>
              <p className="text-t-dark">
                {data.item.artists.map((artist) => artist.name).join(", ")} -{" "}
                {data.item.album.name}
              </p>
            </Marquee>
          </div>
          <p className="text-accent-light text-sm font-bold">
            {msToTime(fakeTime ? fakeProgress : data.progress_ms)} /{" "}
            {msToTime(data.item.duration_ms)}
          </p>
        </div>
      </div>
    </div>
  ) : (
    <div className="bg-background-dark border border-background-light rounded-md w-fit">
      <h1 className={clsx("text-accent-light font-bold px-2 pt-1")}>
        Now Playing
      </h1>

      <p className={clsx("text-t-dark text-sm p-2")}>
        Nothing is playing at the moment
      </p>
    </div>
  );
};

function msToTime(ms: number) {
  const seconds = Math.floor(ms / 1000);
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;
  return `${minutes}:${remainingSeconds < 10 ? "0" : ""}${remainingSeconds}`;
}
