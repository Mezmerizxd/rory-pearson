import React from "react";
import { Dialog, Transition } from "@headlessui/react";
import {
  APP_LOGO_IMAGE_URL,
  APP_NAME,
  navigation_items,
} from "../../constants";
import { useNavigate } from "react-router-dom";
import { Button } from "@components/Elements/Button";

const Navigation = ({ onSelect }: { onSelect?: () => void }) => {
  const navigate = useNavigate();

  return (
    <>
      {navigation_items.map((item, index) => (
        <li key={index} className="list-none">
          <a
            onClick={() => {
              navigate(item.to);
            }}
            className="block py-2 px-3 md:p-0 text-t-light hover:text-t-dark bg-background-light transition rounded md:bg-transparent md:text-t-light cursor-pointer"
          >
            {item.name}
          </a>
        </li>
      ))}
    </>
  );
};

type MobileSidebarProps = {
  sidebarOpen: boolean;
  setSidebarOpen: React.Dispatch<React.SetStateAction<boolean>>;
};

const MobileHeader = ({ sidebarOpen, setSidebarOpen }: MobileSidebarProps) => {
  return (
    <Transition.Root show={sidebarOpen} as={React.Fragment}>
      <Dialog
        as="div"
        static
        className="fixed inset-28 z-40"
        open={sidebarOpen}
        onClose={setSidebarOpen}
      >
        <Transition.Child
          as={React.Fragment}
          enter="transition-opacity ease-linear duration-150"
          enterFrom="opacity-0"
          enterTo="opacity-100"
          leave="transition-opacity ease-linear duration-150"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <Dialog.Overlay className="fixed inset-0 bg-black bg-opacity-30" />
        </Transition.Child>
        <Transition.Child
          as={React.Fragment}
          enter="transition ease-in-out duration-150 transform"
          enterFrom="-translate-y-full"
          enterTo="translate-y-0"
          leave="transition ease-in-out duration-150 transform"
          leaveFrom="translate-y-0"
          leaveTo="-translate-y-full"
        >
          <div className="relative flex-1 flex flex-col w-full h-fit max-h-96 mx-auto overflow-y-auto bg-background-dark border border-white/10 rounded-md">
            <div className="flex flex-col justify-between h-full">
              <div className="flex flex-col space-y-2 p-4">
                <Navigation onSelect={() => setSidebarOpen(false)} />
              </div>
            </div>
          </div>
        </Transition.Child>
      </Dialog>
    </Transition.Root>
  );
};

const HomeHeader = () => {
  const [sidebarOpen, setSidebarOpen] = React.useState(false);

  return (
    <div className="fixed w-full z-50">
      <nav className="bg-background-dark">
        <div className="max-w-screen-xl flex flex-wrap items-center justify-between mx-auto p-4">
          <a
            href="https://rory-pearson.com/"
            className="flex items-center space-x-3 rtl:space-x-reverse"
          >
            <img
              src={APP_LOGO_IMAGE_URL}
              className="h-10"
              alt="Flowbite Logo"
            />
            <span className="text-t-light self-center text-2xl font-semibold whitespace-nowrap">
              {APP_NAME}
            </span>
          </a>
          <div className="flex md:order-2 space-x-3 md:space-x-0 rtl:space-x-reverse">
            {/* <Button size="sm" variant="primary">
              Get started
            </Button> */}

            <button
              type="button"
              className="inline-flex items-center p-2 w-10 h-10 justify-center text-sm text-t-dark-500 transition rounded-lg md:hidden hover:bg-background-light focus:outline-none focus:ring-2 focus:ring-gray-200"
              onClick={() => setSidebarOpen(!sidebarOpen)}
            >
              <span className="sr-only">Open main menu</span>
              <svg
                className="w-5 h-5 text-t-light"
                aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 17 14"
              >
                <path
                  stroke="currentColor"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M1 1h15M1 7h15M1 13h15"
                />
              </svg>
            </button>
          </div>
          <div
            className="items-center justify-between hidden w-full md:flex md:w-auto md:order-1"
            id="navbar-cta"
          >
            <ul className="flex flex-col font-medium p-4 md:p-0 mt-4 border border-gray-100 rounded-lg md:space-x-8 rtl:space-x-reverse md:flex-row md:mt-0 md:border-0">
              <Navigation />
            </ul>
          </div>
        </div>
      </nav>

      <MobileHeader sidebarOpen={sidebarOpen} setSidebarOpen={setSidebarOpen} />
    </div>
  );
};

export default HomeHeader;
