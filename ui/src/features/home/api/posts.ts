import { useMutation, useQuery } from "react-query";

import {
  ExtractFnReturnType,
  QueryConfig,
  MutationConfig,
} from "../../../util/react-query";
import { useNotificationStore } from "../../../stores/notifications";

function getHostname() {
  if (window.location.hostname === "localhost") {
    return "http://localhost:3000";
  }

  return window.location.origin;
}

export const getPosts = async (
  page: number,
  pageSize: number
): Promise<GetPostData> => {
  const hostname = getHostname();

  try {
    const response = await fetch(
      `${hostname}/board/get?page=${page}&pageSize=${pageSize}`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      }
    );

    if (!response.ok) {
      throw new Error(`Error fetching posts: ${response.statusText}`);
    }

    const data = await response.json();

    return data ?? { posts: [], total_posts: 0, page, pageSize };
  } catch (error) {
    console.error(error);
    return { posts: [], total_posts: 0, page, page_size: 0 }; // Return default data on error
  }
};

type QueryFnType = typeof getPosts;

type UseGetPostsOptions = {
  config?: QueryConfig<QueryFnType>;
  page: number;
  pageSize: number;
};

export const usePosts = ({ config, page, pageSize }: UseGetPostsOptions) => {
  return useQuery<ExtractFnReturnType<QueryFnType>>({
    ...config,
    queryKey: ["posts", page, pageSize],
    queryFn: () => getPosts(page, pageSize),
  });
};

export const CreatePost = async (data: CreatePostData) => {
  const hostname = getHostname();
  const response = await fetch(`${hostname}/board/create`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });
};

type UseCreatePostOptions = {
  config?: MutationConfig<typeof CreatePost>;
};

export const useCreatePost = ({ config }: UseCreatePostOptions = {}) => {
  const { addNotification } = useNotificationStore();
  const postsQuery = usePosts({ page: 1, pageSize: 8 });

  return useMutation({
    onError: (error, __, context: any) => {
      addNotification({
        type: "error",
        title: "Create Post",
        message: error.message,
      });
    },
    onSuccess: (newPost, __, context: any) => {
      addNotification({
        type: "success",
        title: "Create Post",
        message: "Successfully created post",
      });

      postsQuery.refetch();
    },
    ...config,
    mutationFn: CreatePost,
  });
};
