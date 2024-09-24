import React from "react";
import { Transition } from "@headlessui/react";
import { MdClose } from "react-icons/md";
import { FaCheckCircle, FaExclamationCircle } from "react-icons/fa";
import { FaCircleInfo } from "react-icons/fa6";
import { IoIosCloseCircle } from "react-icons/io";

import { Fragment, useEffect, useState } from "react";

const icons = {
  info: <FaCircleInfo className="h-6 w-6 text-blue-500" aria-hidden="true" />,
  success: (
    <FaCheckCircle className="h-6 w-6 text-green-500" aria-hidden="true" />
  ),
  warning: (
    <FaExclamationCircle
      className="h-6 w-6 text-yellow-500"
      aria-hidden="true"
    />
  ),
  error: (
    <IoIosCloseCircle className="h-6 w-6 text-red-500" aria-hidden="true" />
  ),
};

export type NotificationProps = {
  notification: {
    id: string;
    type: keyof typeof icons;
    title: string;
    message?: string;
  };
  onDismiss: (id: string) => void;
};

export const Notification = ({
  notification: { id, type, title, message },
  onDismiss,
}: NotificationProps) => {
  const [remainingTime, setRemainingTime] = useState(4000);

  useEffect(() => {
    const timer = setTimeout(() => {
      onDismiss(id);
    }, remainingTime);

    const interval = setInterval(() => {
      setRemainingTime((prevTime) => prevTime - 10); // Update every 100ms for a smoother progress bar
    }, 10);

    return () => {
      clearInterval(interval);
      clearTimeout(timer);
    };
  }, [id, onDismiss, remainingTime, 4000]);

  const widthPercentage = (remainingTime / 4000) * 100;

  return (
    <div className="w-full flex flex-col items-center space-y-4 sm:items-end">
      <Transition
        show={true}
        as={Fragment}
        enter="transform ease-out duration-300 transition"
        enterFrom="translate-y-2 opacity-0 sm:translate-y-0 sm:translate-x-2"
        enterTo="translate-y-0 opacity-100 sm:translate-x-0"
        leave="transition ease-in duration-100"
        leaveFrom="opacity-100"
        leaveTo="opacity-0"
      >
        <div className="max-w-sm w-full bg-background-dark shadow-lg rounded-lg pointer-events-auto ring-1 ring-black ring-opacity-5 overflow-hidden">
          <div className="p-4" role="alert" aria-label={title}>
            <div className="flex items-start">
              <div className="flex-shrink-0">{icons[type]}</div>
              <div className="ml-3 w-0 flex-1 pt-0.5">
                <p className="text-sm font-medium text-t-dark">{title}</p>
                <p className="mt-1 text-sm text-t-dark">{message}</p>
              </div>
              <div className="ml-4 flex-shrink-0 flex">
                <button
                  className="bg-background-dark rounded-md inline-flex text-accent-light hover:text-accent-dark duration-150 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-accent-light"
                  onClick={() => {
                    onDismiss(id);
                  }}
                >
                  <span className="sr-only">Close</span>
                  <MdClose className="h-5 w-5" aria-hidden="true" />
                </button>
              </div>
            </div>
          </div>
          <div
            className="bg-accent-light"
            style={{ width: `${widthPercentage}%`, height: "2px" }}
          />
        </div>
      </Transition>
    </div>
  );
};
