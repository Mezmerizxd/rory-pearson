type RootState = ReturnType<typeof store.getState>;
type AppDispatch = typeof store.dispatch;

declare module "*.jpg" {
  const value: string;
  export default value;
}
declare module "*.png" {
  const value: string;
  export default value;
}

type PostData = {
  id: number;
  title: string;
  content: string;
  author: string;
  created_at: number;
};

type CreatePostData = {
  title: string;
  content: string;
  author: string;
};

type GetPostData = {
  posts: PostData[];
  total_posts: number;
  page: number;
  page_size: number;
};

type SpotifyUserData = {
  display_name: string;
  external_urls: {
    spotify: string;
  };
  followers: {
    total: number;
    href: string;
  };
  href: string;
  id: string;
  images: {
    height: number;
    width: number;
    url: string;
  }[];
  uri: string;
  country: string;
  email: string;
  product: string;
  birthdate: string;
};

type SpotifyPlaylistItemData = {
  collaborative: boolean;
  description: string;
  external_urls: {
    spotify: string;
  };
  href: string;
  id: string;
  images: {
    height: number;
    width: number;
    url: string;
  }[];
  name: string;
  owner: SpotifyUserData;
  public: boolean;
  snapshotID: string;
  tracks: {
    href: string;
    total: number;
  };
  uri: string;
};

type SpotifyPlaylistData = {
  href: string;
  limit: number;
  offset: number;
  total: number;
  next: string | null;
  previous: string | null;
  items: SpotifyPlaylistItemData[];
};

type SpotifyTrackData = {
  artists: {
    name: string;
    id: string;
    uri: string;
    href: string;
    external_urls: {
      spotify: string;
    };
  }[];
  available_markets: string[];
  disc_number: number;
  duration_ms: number;
  explicit: boolean;
  external_urls: {
    spotify: string;
  };
  href: string;
  id: string;
  name: string;
  preview_url: string | null;
  track_number: number;
  uri: string;
  album: {
    name: string;
    id: string;
    uri: string;
    href: string;
    external_urls: {
      spotify: string;
    };
    images: {
      height: number;
      width: number;
      url: string;
    }[];
    release_date: string;
  };
};

// type SpotifyPlaylistItemData = {
//   href: string;
//   limit: number;
//   offset: number;
//   total: number;
//   next: string | null;
//   previous: string | null;
//   items: {
//     added_at: string;
//     added_by: SpotifyUserData;
//     is_local: boolean;
//     track: SpotifyTrackData;
//     episode: any | null;
//   }[];
// };

type SpotifyNowPlayingData = {
  timestamp: number;
  context: {
    external_urls: {
      spotify: string;
    };
    href: string;
    type: string;
    uri: string;
  };
  progress_ms: number;
  is_playing: boolean;
  item: SpotifyTrackData;
};

type YoutubeVideoData = {
  id: string;
  title: string;
  channel: string;
  views: number;
  url: string;
};
