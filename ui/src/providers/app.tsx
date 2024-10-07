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
import { Notifications } from "@components/Notifications";
import { useNotificationStore } from "../stores/notifications";
import { GetHost } from "../util/host";

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
  const notifications = useNotificationStore();

  React.useEffect(() => {
    setTimeout(async () => {
      let errorMsg = null;

      try {
        const response = await fetch(`${GetHost()}/api/ping`);

        const data = await response.json();
        if (!data) {
          errorMsg = "Failed to connect to the server";
        }
        if (data.error) {
          errorMsg = data.error;
        }
      } catch (error) {
        errorMsg = `Failed to send request to the server, ${error}`;
      }

      // if there is an error, show it
      if (errorMsg) {
        notifications.addNotification({
          type: "error",
          title: "Server Error",
          message: errorMsg || "Failed to connect to the server",
        });
      }
    }, 5000);
  }, []);

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
              <Notifications />
              {process.env.NODE_ENV !== "test" && <ReactQueryDevtools />}
              <Router>{children}</Router>
            </QueryClientProvider>
          </HelmetProvider>
        </ErrorBoundary>
      </ThemeProvider>
    </React.Suspense>
  );
};
