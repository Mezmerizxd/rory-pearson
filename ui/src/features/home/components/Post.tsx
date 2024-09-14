import clsx from "clsx";
import React from "react";
import { formatDate } from "../../../util/format";

export type PostProps = {
  post: PostData;
  className?: string;
};

export const Post = ({ post, className }: PostProps) => {
  return (
    <div
      className={clsx(
        "px-2 grid grid-cols-1",
        "bg-background-dark border border-t-light/20",
        "rounded-md divide-y divide-t-light/20 text-left",
        ...className
      )}
    >
      <div className="px-2 pb-2 pt-4">
        <h1 className="font-bold">{post.title}</h1>
      </div>

      <div className="p-4">
        <p>{post.content}</p>
      </div>

      <div className="flex justify-between items-center p-2 font-bold text-sm">
        <p className="text-xs">{formatDate(post.created_at * 1000)}</p>
        <p>{post.author}</p>
      </div>
    </div>
  );
};
