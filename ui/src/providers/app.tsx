import * as React from "react";
import { ErrorBoundary } from "react-error-boundary";
import { HelmetProvider } from "react-helmet-async";
import { QueryClientProvider } from "react-query";
import { ReactQueryDevtools } from "react-query/devtools";
import { BrowserRouter as Router } from "react-router-dom";
import { ThemeProvider } from "@material-tailwind/react";
import { queryClient } from "../util/react-query";
import { Spinner } from "@components/Elements";
import { Button } from "@components/Elements/Button";

const ErrorFallback = () => {
  return (
    <div
      className="text-t-light w-screen h-screen flex flex-col justify-center items-center"
      role="alert"
    >
      <h2 className="text-lg font-semibold">Something went wrong</h2>
      <Button
        variant="danger"
        size="sm"
        className="mt-4"
        onClick={() => window.location.assign(window.location.origin)}
      >
        Refresh
      </Button>
    </div>
  );
};

type AppProviderProps = {
  children: React.ReactNode;
};

export const AppProvider = ({ children }: AppProviderProps) => {
  return (
    <React.Suspense
      fallback={
        <div className="flex items-center justify-center w-screen h-screen bg-black">
          <Spinner size="xl" />
        </div>
      }
    >
      <ThemeProvider>
        <ErrorBoundary FallbackComponent={ErrorFallback}>
          <HelmetProvider>
            <QueryClientProvider client={queryClient}>
              {process.env.NODE_ENV !== "test" && <ReactQueryDevtools />}
              <Router>{children}</Router>
            </QueryClientProvider>
          </HelmetProvider>
        </ErrorBoundary>
      </ThemeProvider>
    </React.Suspense>
  );
};
