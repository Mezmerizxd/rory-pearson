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

export const BackgroundRemover = async (file: File): Promise<Blob | null> => {
  const hostname = getHostname();

  const formData = new FormData();
  formData.append("file", file);

  // Returns image data
  const response = await fetch(`${hostname}/api/background-remover`, {
    method: "POST",
    body: formData,
  });

  if (!response.ok) {
    throw new Error("Failed to process image");
  }

  return response.blob();
};

type UseBackgroundRemoverOptions = {
  config?: MutationConfig<typeof BackgroundRemover>;
};

export const useBackgroundRemover = ({
  config,
}: UseBackgroundRemoverOptions = {}) => {
  const { addNotification } = useNotificationStore();

  return useMutation({
    onError: (error, __, context: any) => {
      addNotification({
        type: "error",
        title: "Background Remover",
        message: error.message,
      });
    },
    onSuccess: (image, __, context: any) => {
      addNotification({
        type: "success",
        title: "Background Remover",
        message: "Successfully uploaded image",
      });
    },
    ...config,
    mutationFn: BackgroundRemover,
  });
};
