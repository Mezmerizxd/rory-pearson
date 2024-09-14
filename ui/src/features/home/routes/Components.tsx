import React from "react";
import {
  Dialog,
  Spinner,
  Button,
  DialogTitle,
  DialogDescription,
} from "@components/Elements";
import { FaAdjust } from "react-icons/fa";

const HomeLayout = React.lazy(() => import("@components/Layout/HomeLayout"));

export const Components = () => {
  const [dialog1, setDialog1] = React.useState(false);

  return (
    <HomeLayout title="Home">
      <div className="p-5 space-y-4">
        <h1 className="text-center p-5">Buttons</h1>
        <div className="grid grid-cols-4 gap-4">
          <Button variant="primary">Primary</Button>
          <Button variant="secondary">Secondary</Button>
          <Button variant="danger">Danger</Button>
          <Button variant="success">Success</Button>
          <Button variant="warning">Warning</Button>
        </div>

        <div className="flex space-x-10">
          <div className="space-y-4">
            <Button variant="primary" size="xs">
              Extra Small
            </Button>
            <Button variant="primary" size="sm">
              Small
            </Button>
            <Button variant="secondary" size="md">
              Medium
            </Button>
            <Button variant="danger" size="lg">
              Large
            </Button>
            <Button variant="success" size="xl">
              Extra Large
            </Button>
          </div>

          <div className="space-y-4">
            <Button variant="primary" size="xs" isLoading>
              Extra Small
            </Button>
            <Button variant="primary" size="sm" isLoading>
              Small
            </Button>
            <Button variant="secondary" size="md" isLoading>
              Medium
            </Button>
            <Button variant="danger" size="lg" isLoading>
              Large
            </Button>
            <Button variant="success" size="xl" isLoading>
              Extra Large
            </Button>
          </div>

          <div className="space-y-4">
            <Button variant="primary" size="xs" startIcon={<FaAdjust />}>
              Extra Small
            </Button>
            <Button variant="primary" size="sm" endIcon={<FaAdjust />}>
              Small
            </Button>
            <Button variant="secondary" size="md" startIcon={<FaAdjust />}>
              Medium
            </Button>
            <Button variant="danger" size="lg" endIcon={<FaAdjust />}>
              Large
            </Button>
            <Button variant="success" size="xl" startIcon={<FaAdjust />}>
              Extra Large
            </Button>
          </div>
        </div>
      </div>

      <div className="mt-10">
        <h1 className="text-center">Spinners</h1>
        <div className="grid grid-cols-3 space-x-3 p-3 mt-5">
          <Spinner className="m-5" size="sm" />
          <Spinner className="m-5" size="md" />
          <Spinner className="m-5" size="lg" />
        </div>
      </div>

      <div className="p-5 space-y-4">
        <h1 className="text-center p-5">Dialogs</h1>
        <Button onClick={() => setDialog1(true)}>Dialog 1</Button>
        {dialog1 && (
          <Dialog1 isOpen={dialog1} onClose={() => setDialog1(false)} />
        )}
      </div>
    </HomeLayout>
  );
};

const Dialog1 = ({
  isOpen,
  onClose,
}: {
  isOpen: boolean;
  onClose: () => void;
}) => {
  return (
    <Dialog className="p-4 w-96" isOpen={isOpen} onClose={onClose}>
      <DialogTitle as="h2" className="text-lg font-medium leading-6">
        Dialog Title
      </DialogTitle>
      <DialogDescription className="mt-2">
        Lorem ipsum dolor sit amet consectetur adipisicing elit. Quisquam,
        voluptatum. lorem ipsum dolor sit amet consectetur adipisicing elit.
        Quisquam, voluptatum.
      </DialogDescription>
      <div className="mt-4">
        <Button
          variant="primary"
          size="sm"
          className="inline-flex justify-center"
          onClick={onClose}
        >
          Close
        </Button>
      </div>
    </Dialog>
  );
};
