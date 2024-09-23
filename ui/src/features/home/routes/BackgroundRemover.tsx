import React from "react";
import { FileInput, FileInputProps, Label } from "flowbite-react";
import { Button, Spinner } from "@components/Elements";
import { FaDownload, FaUpload } from "react-icons/fa";
import { FaRotateRight } from "react-icons/fa6";
import { useBackgroundRemover } from "../api/background-remover";

const HomeLayout = React.lazy(() => import("@components/Layout/HomeLayout"));

export const BackgroundRemover = () => {
  const [error, setError] = React.useState<string | null>(null);
  const [file, setFile] = React.useState<File | null>(null);
  const backgroundRemoverMutation = useBackgroundRemover();

  function collectFile(e: React.ChangeEvent<HTMLInputElement>) {
    reset();

    const file = e.target.files && e.target.files[0];
    if (!file || !file.type.includes("image") || e.target.files.length > 1) {
      setError("Please upload a valid image file");
      return;
    }

    setFile(file);
  }

  async function upload() {
    reset();

    if (file === null) {
      setError("Please upload a valid image file");
      return;
    }

    const data = await backgroundRemoverMutation.mutateAsync(file);

    if (backgroundRemoverMutation.isError) {
      setError(
        `An error occurred while processing the image, ${backgroundRemoverMutation.error.message}`
      );
      return;
    }

    if (data === null) {
      setError("An error occurred while processing the image");
      return;
    }

    // set file as new file
    setFile(new File([data], file.name, { type: "image/png" }));
  }

  function reset() {
    backgroundRemoverMutation.reset();
    setFile(null);
    setError(null);
  }

  return (
    <>
      <HomeLayout title="Background Remover">
        <div className="p-10">
          <div className="p-5 flex justify-center align-middle">
            {error && (
              <p className="p-2 bg-orange-500/20 text-orange-500 border border-orange-500/20">
                {error}
              </p>
            )}
          </div>

          {backgroundRemoverMutation.isLoading ? (
            <>
              <div className="p-5 flex justify-center align-middle">
                <Spinner size="lg" />
              </div>
            </>
          ) : (
            <>
              {file === null && (
                <div className="flex w-full items-center justify-center">
                  <Label
                    htmlFor="dropzone-file"
                    className="flex h-64 w-full cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed border-background-dark bg-background-light hover:border-background-dark/50 hover:bg-background-light/50"
                  >
                    <div className="flex flex-col items-center justify-center pb-6 pt-5">
                      <svg
                        className="mb-4 h-8 w-8 text-gray-400"
                        aria-hidden="true"
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 20 16"
                      >
                        <path
                          stroke="currentColor"
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth="2"
                          d="M13 13h3a3 3 0 0 0 0-6h-.025A5.56 5.56 0 0 0 16 6.5 5.5 5.5 0 0 0 5.207 5.021C5.137 5.017 5.071 5 5 5a4 4 0 0 0 0 8h2.167M10 15V6m0 0L8 8m2-2 2 2"
                        />
                      </svg>
                      <p className="mb-2 text-sm text-t-dark">
                        <span className="font-semibold">Click to upload</span>{" "}
                        or drag and drop
                      </p>
                      <p className="text-xs text-t-dark">PNG, JPG, WEBP</p>
                    </div>
                    <FileInput
                      id="dropzone-file"
                      className="hidden"
                      onChange={(e) => e.target.files && collectFile(e)}
                    />
                  </Label>
                </div>
              )}
            </>
          )}

          {file !== null && (
            <div className="flex w-full items-center justify-center">
              <img
                src={URL.createObjectURL(file)}
                alt="Uploaded Image"
                className="max-h-96"
              />
            </div>
          )}

          <div className="p-5 flex justify-center align-middle space-x-4">
            {backgroundRemoverMutation.isSuccess ? (
              <Button
                startIcon={<FaDownload />}
                onClick={() => {
                  const url = URL.createObjectURL(file);
                  const a = document.createElement("a");
                  a.href = url;
                  a.download = file.name;
                  a.click();
                  URL.revokeObjectURL(url);
                }}
              >
                Download
              </Button>
            ) : (
              <Button startIcon={<FaUpload />} onClick={() => upload()}>
                Upload
              </Button>
            )}
            <Button startIcon={<FaRotateRight />} onClick={() => reset()}>
              Reset
            </Button>
          </div>
        </div>
      </HomeLayout>
    </>
  );
};
