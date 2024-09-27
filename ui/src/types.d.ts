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

type SpotifyPlaylistData = {
  items: {
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
    isPublic: boolean;
    snapshotID: string;
    tracks: {
      href: string;
      total: number;
    };
    uri: string;
  }[];
};
