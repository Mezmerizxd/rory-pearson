import React, { useState } from "react";
import { Button, Dialog, DialogTitle } from "@components/Elements";
import { Input } from "@components/Elements";
import { TextArea } from "@components/Elements/Input/TextArea";

interface CreatePostDialogProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: (post: CreatePostData) => void;
}

export const CreatePostDialog = ({
  isOpen,
  onClose,
  onConfirm,
}: CreatePostDialogProps) => {
  const [error, setError] = useState<string>(null);
  const [post, setPost] = useState<CreatePostData>({
    title: "",
    content: "",
    author: "",
  });

  function create() {
    setError(null);

    if (post.title == null || post.title === "") {
      setError("Title is required");
      return;
    }

    if (post.content == null || post.content === "") {
      setError("Content is required");
      return;
    }

    if (post.author == null || post.author === "") {
      setError("Author is required");
      return;
    }

    onConfirm(post);
    onClose();
  }

  return (
    <Dialog className="p-5" isOpen={isOpen} onClose={onClose}>
      <DialogTitle as="h2" className="text-lg font-medium">
        Create Post
      </DialogTitle>
      <div className="w-96 m-2 space-y-6 overflow-auto py-4">
        <Input
          label="Title"
          onChange={(e) => {
            setPost({
              ...post,
              title: e.target.value,
            });
          }}
        />
        <TextArea
          label="Content"
          onChange={(e) => {
            setPost({
              ...post,
              content: e.target.value,
            });
          }}
        />
        <Input
          label="Author"
          onChange={(e) => {
            setPost({
              ...post,
              author: e.target.value,
            });
          }}
        />

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
          onClick={() => create()}
        >
          Create
        </Button>
      </div>
    </Dialog>
  );
};
