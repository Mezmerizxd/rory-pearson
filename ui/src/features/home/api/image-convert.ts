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

export const GetDownloadUrl = (id: string): string => {
  const hostname = getHostname();

  return `${hostname}/api/image-convert/download/${id}`;
};

export const ImageConvert = async (
  file: File
): Promise<{
  message: string;
  download_id: string;
}> => {
  const hostname = getHostname();

  const formData = new FormData();
  formData.append("file", file);

  // Returns image data
  const response = await fetch(`${hostname}/api/image-convert/upload`, {
    method: "POST",
    body: formData,
  });

  if (!response.ok) {
    throw new Error("Failed to process image");
  }

  return response.json();
};

type UseImageConvertOptions = {
  config?: MutationConfig<typeof ImageConvert>;
};

export const useImageConvert = ({ config }: UseImageConvertOptions = {}) => {
  const { addNotification } = useNotificationStore();

  return useMutation({
    onError: (error, __, context: any) => {
      addNotification({
        type: "error",
        title: "Image Convert",
        message: error.message,
      });
    },
    onSuccess: (image, __, context: any) => {
      addNotification({
        type: "success",
        title: "Image Convert",
        message: "Successfully uploaded image",
      });
    },
    ...config,
    mutationFn: ImageConvert,
  });
};
