import { Button } from "@components/Elements/Button";
import React from "react";
import { useNavigate } from "react-router-dom";

const HomeLayout = React.lazy(() => import("@components/Layout/HomeLayout"));

export const Landing = () => {
  const navigate = useNavigate();

  return (
    <HomeLayout title="Home">
      <div className="flex items-center justify-center min-h-screen p-10">
        <div className="text-center">
          <h1 className="mb-4 text-4xl font-extrabold leading-none tracking-tight text-t-light md:text-5xl lg:text-6xl">
            Welcome, Create Some Posts On The Board
          </h1>
          <p className="mb-6 text-lg font-normal text-t-dark lg:text-xl sm:px-16 xl:px-48">
            In the Board, You will be able to create Posts about whatever you
            want, and other users will be able to see them.
          </p>
          <Button
            className="inline-flex items-center justify-center px-5 py-3"
            variant="primary"
            size="lg"
            onClick={() => navigate("/board")}
          >
            Start Posting
          </Button>
        </div>
      </div>
    </HomeLayout>
  );
};
