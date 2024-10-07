import styled from "styled-components";
import clsx from "clsx";
import * as React from "react";
import { IconButton } from "@components/Elements";
import { MdEdit } from "react-icons/md";

const sizes = {
  sm: "text-sm",
  md: "text-md",
  lg: "text-lg",
};

type PlaylistOption = {
  name: string;
  func: (playlist: SpotifyPlaylistItemData) => void;
};

export type PlaylistThumbnailProps = {
  data: SpotifyPlaylistItemData;
  options?: PlaylistOption[];
  size?: keyof typeof sizes;
  className?: string;
};

export const PlaylistThumbnail: React.FC<PlaylistThumbnailProps> = ({
  data,
  options = [],
  size = "md",
  className = "",
}) => {
  const [showOptions, setShowOptions] = React.useState(false);
  const optionsRef = React.useRef<HTMLDivElement>(null);

  function handleOptions() {
    setShowOptions((prev) => !prev);
  }

  // Close popup if clicking outside of it
  React.useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (
        optionsRef.current &&
        !optionsRef.current.contains(event.target as Node)
      ) {
        setShowOptions(false);
      }
    }

    // Add event listener to the document
    document.addEventListener("mousedown", handleClickOutside);

    // Cleanup the event listener on component unmount
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  return (
    data && (
      <div
        className={clsx(
          "w-full max-h-[500px] p-2 space-y-2 flex flex-col",
          "bg-background-dark border border-background-light rounded-md w-fit",
          sizes[size],
          className
        )}
      >
        <div className="flex justify-between items-center">
          <h1 className={clsx("text-accent-light font-bold py-1")}>Playlist</h1>

          {/* Only render options button if options are available */}
          {options.length > 0 && (
            <div className="relative" ref={optionsRef}>
              <IconButton size="xs" icon={<MdEdit />} onClick={handleOptions} />
              {showOptions && (
                <div className="absolute top-full right-0 p-1 bg-background border border-background-light rounded shadow-lg z-10 w-fit">
                  {options.map((option, index) => (
                    <div
                      key={index}
                      className={clsx(
                        "cursor-pointer px-4 py-2 hover:bg-background-dark",
                        "text-sm whitespace-nowrap"
                        // stop the wrods from breaking
                      )}
                      onClick={() => option.func(data)}
                    >
                      {option.name}
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}
        </div>

        <div className="flex-grow flex flex-col">
          <div className="mb-2">
            {data.images?.length > 0 ? (
              <img
                key={data.images[0].url}
                src={data.images[0].url}
                width={data.images[0].width}
                height={data.images[0].height}
                alt="Playlist Thumbnail"
                className="w-full h-48 object-cover rounded-md"
              />
            ) : (
              <div className="w-full h-48 bg-background-light rounded-md">
                <div className="flex items-center justify-center h-full">
                  <p className="text-accent-light">No Image</p>
                </div>
              </div>
            )}
          </div>

          <div className="flex-grow">
            <h1 className="font-bold text-accent-light">{data.name}</h1>

            <DescriptionContainer
              dangerouslySetInnerHTML={{ __html: data.description }}
            />
          </div>
        </div>

        <div
          id="footer"
          className="flex justify-between mt-auto text-sm font-bold"
        >
          <p>
            {data.public === true ? (
              <span className="text-t-light/80">Public</span>
            ) : (
              <span className="text-t-dark/80">Private</span>
            )}
          </p>
          <p>Tracks: {data.tracks.total}</p>
        </div>
      </div>
    )
  );
};

const DescriptionContainer = styled.div`
  overflow-y: auto;
  overflow-x: hidden;
  max-height: 80px;
  a {
    text-decoration: underline;
    font-weight: bold;
  }
  a:hover {
    text-decoration: underline;
  }
`;
