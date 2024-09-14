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
