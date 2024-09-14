import React, { useEffect, useState } from "react";
import { FaArrowLeft, FaArrowRight } from "react-icons/fa";
import { Post } from "../components/Post";
import { Button } from "@components/Elements/Button";
import { Spinner } from "@components/Elements";
import { CreatePostDialog } from "../components/CreatePostDialog";
import { useCreatePost, usePosts } from "../api/posts";

const HomeLayout = React.lazy(() => import("@components/Layout/HomeLayout"));

export const Board = () => {
  const maxPosts = 8;
  const [page, setPage] = useState(1);
  const postsQuery = usePosts({ page: page, pageSize: maxPosts });
  const createPost = useCreatePost();
  const [createPostDialog, setCreatePostDialog] = React.useState(false);

  function increasePage() {
    setPage((prev) => prev + 1);
  }

  function decreasePage() {
    setPage((prev) => prev - 1);
  }

  async function createPostHandler(post) {
    await createPost.mutateAsync(post);
    postsQuery.refetch();
  }

  return (
    <>
      <HomeLayout title="Board">
        <CreatePostDialog
          key={"create-post-dialog"}
          isOpen={createPostDialog}
          onClose={() => setCreatePostDialog(false)}
          onConfirm={(post) => {
            createPostHandler(post);
          }}
        />

        <div className="p-10">
          <div className="flex justify-center items-center pb-4">
            <Button size="md" onClick={() => setCreatePostDialog(true)}>
              Create Post
            </Button>
          </div>

          <div className="flex justify-between items-center space-x-2 my-5">
            <Button
              size="sm"
              onClick={decreasePage}
              disabled={postsQuery?.data?.page === 1}
            >
              <FaArrowLeft />
            </Button>
            <p className="flex justify-center items-center">
              Page: {postsQuery?.data?.page}/
              {Math.ceil(postsQuery?.data?.total_posts / maxPosts)}
            </p>
            <Button
              size="sm"
              onClick={increasePage}
              disabled={
                postsQuery?.data?.page * maxPosts >=
                postsQuery?.data?.total_posts
              }
            >
              <FaArrowRight />
            </Button>
          </div>

          {postsQuery.isLoading ? (
            <div className="flex justify-center items-center">
              <Spinner className="ml-2" size="lg" />
            </div>
          ) : (
            <div className="grid xs:grid-cols-1 sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
              {postsQuery.data?.posts?.map((post) => (
                <Post key={post.id} className="animate-fadeIn" post={post} />
              ))}
            </div>
          )}
        </div>
      </HomeLayout>
    </>
  );
};
