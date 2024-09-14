import { Dialog as UIDialog, Transition } from "@headlessui/react";
import clsx from "clsx";
import * as React from "react";

type DialogProps = {
  className?: string;
  isOpen: boolean;
  onClose: () => void;
  children: React.ReactNode;
  initialFocus?: React.MutableRefObject<null>;
};

export const DialogTitle = UIDialog.Title;

export const DialogDescription = UIDialog.Description;

export const Dialog = ({
  className,
  isOpen,
  onClose,
  children,
  initialFocus,
}: DialogProps) => {
  return (
    <>
      <Transition appear show={isOpen} as={React.Fragment}>
        <UIDialog as="div" className="relative z-10" onClose={onClose}>
          <Transition.Child
            as={React.Fragment}
            enter="ease-out duration-100"
            enterFrom="opacity-0"
            enterTo="opacity-100"
            leave="ease-in duration-100"
            leaveFrom="opacity-100"
            leaveTo="opacity-0"
          >
            <div className="fixed inset-0" />
          </Transition.Child>

          <div className="fixed inset-0 overflow-y-auto bg-black/30">
            <div className="flex min-h-full items-center justify-center text-center">
              <Transition.Child
                as={React.Fragment}
                enter="ease-out duration-100"
                enterFrom="opacity-0 scale-0"
                enterTo="opacity-100 scale-100"
                leave="ease-in duration-100"
                leaveFrom="opacity-100 scale-100"
                leaveTo="opacity-0 scale-0"
              >
                <UIDialog.Panel
                  className={clsx(
                    "w-min transform overflow-hidden rounded-md bg-background-dark text-t-light transition-all",
                    "border border-background-light",
                    className
                  )}
                >
                  {children}
                </UIDialog.Panel>
              </Transition.Child>
            </div>
          </div>
        </UIDialog>
      </Transition>
    </>
  );
};
