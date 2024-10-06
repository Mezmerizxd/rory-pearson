import { Dialog, Menu, Transition } from "@headlessui/react";
import * as React from "react";
import { NavLink, Link } from "react-router-dom";
import { clsx } from "clsx";
import { IoIosCloseCircle } from "react-icons/io";
import { RxHamburgerMenu } from "react-icons/rx";

import { APP_NAME, APP_LOGO_IMAGE_URL } from "../../constants";

export type DashboardNavigationItem = {
  name: string;
  to: string;
  icon: React.ReactNode;
};

const SideNavigation = ({
  navigation,
  onSelect,
}: {
  navigation: DashboardNavigationItem[];
  onSelect?: () => void;
}) => {
  return (
    <>
      {navigation.length > 0 ? (
        navigation.map((item, index) => (
          <NavLink
            onClick={() => {
              onSelect && onSelect();
            }}
            end={index === 0}
            key={item.name}
            to={item.to}
            className={clsx(
              "text-gray-300 hover:bg-background-light hover:text-accent-light duration-150",
              "group flex items-center px-2 py-2 text-base font-medium rounded-md space-x-2"
            )}
          >
            <p>{item.icon}</p>
            <p>{item.name}</p>
          </NavLink>
        ))
      ) : (
        <div>
          <p className="text-gray-300 text-sm px-2 py-2">No navigation items</p>
        </div>
      )}
    </>
  );
};

export type UserNavigationItem = {
  name: string;
  to: string;
  onClick?: () => void;
};

const UserNavigation = ({
  userNavigation,
  userElement,
}: {
  userNavigation: UserNavigationItem[];
  userElement: React.ReactNode;
}) => {
  return (
    <Menu as="div" className="ml-3 relative">
      {({ open }) => (
        <>
          <div>
            <Menu.Button className="max-w-xs flex items-center text-sm rounded-full">
              <span className="sr-only">Open user menu</span>
              {userElement}
            </Menu.Button>
          </div>
          <Transition
            show={open}
            as={React.Fragment}
            enter="transition ease-out duration-100"
            enterFrom="transform opacity-0 scale-95"
            enterTo="transform opacity-100 scale-100"
            leave="transition ease-in duration-75"
            leaveFrom="transform opacity-100 scale-100"
            leaveTo="transform opacity-0 scale-95"
          >
            <Menu.Items
              static
              className="origin-top-right absolute right-0 mt-2 w-48 rounded-md shadow-lg py-1 bg-background-dark ring-1 ring-black ring-opacity-5 focus:outline-none"
            >
              {userNavigation.map((item) => (
                <Menu.Item key={item.name}>
                  {({ active }) => (
                    <Link
                      onClick={item.onClick}
                      to={item.to}
                      className={clsx(
                        active ? "bg-background-light" : "",
                        "block px-4 py-2 text-sm text-gray-300 hover:text-white-light duration-150"
                      )}
                    >
                      {item.name}
                    </Link>
                  )}
                </Menu.Item>
              ))}
            </Menu.Items>
          </Transition>
        </>
      )}
    </Menu>
  );
};

type MobileSidebarProps = {
  navigation: DashboardNavigationItem[];
  sidebarOpen: boolean;
  setSidebarOpen: React.Dispatch<React.SetStateAction<boolean>>;
};

const MobileSidebar = ({
  navigation,
  sidebarOpen,
  setSidebarOpen,
}: MobileSidebarProps) => {
  return (
    <Transition.Root show={sidebarOpen} as={React.Fragment}>
      <Dialog
        as="div"
        static
        className="fixed inset-0 flex z-40 md:hidden"
        open={sidebarOpen}
        onClose={setSidebarOpen}
      >
        <Transition.Child
          as={React.Fragment}
          enter="transition-opacity ease-linear duration-300"
          enterFrom="opacity-0"
          enterTo="opacity-100"
          leave="transition-opacity ease-linear duration-300"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <Dialog.Overlay className="fixed inset-0 bg-background-dark bg-opacity-75" />
        </Transition.Child>
        <Transition.Child
          as={React.Fragment}
          enter="transition ease-in-out duration-300 transform"
          enterFrom="-translate-x-full"
          enterTo="translate-x-0"
          leave="transition ease-in-out duration-300 transform"
          leaveFrom="translate-x-0"
          leaveTo="-translate-x-full"
        >
          <div className="relative flex-1 flex flex-col max-w-xs w-full pt-5 pb-4 bg-background-dark">
            <Transition.Child
              as={React.Fragment}
              enter="ease-in-out duration-300"
              enterFrom="opacity-0"
              enterTo="opacity-100"
              leave="ease-in-out duration-300"
              leaveFrom="opacity-100"
              leaveTo="opacity-0"
            >
              <div className="absolute top-0 right-0 -mr-12 pt-2">
                <button
                  className="ml-1 flex items-center justify-center h-10 w-10 rounded-full focus:outline-none focus:ring-2 focus:ring-inset focus:ring-accent-light"
                  onClick={() => setSidebarOpen(false)}
                >
                  <span className="sr-only">Close sidebar</span>
                  <IoIosCloseCircle className="text-accent-light" size={30} />
                </button>
              </div>
            </Transition.Child>
            <div className="flex-shrink-0 flex items-center px-4">
              <Logo />
            </div>
            <div className="mt-5 flex-1 h-0 overflow-y-auto">
              <nav className="px-2 space-y-1">
                <SideNavigation
                  navigation={navigation}
                  onSelect={() => setSidebarOpen(false)}
                />
              </nav>
            </div>
          </div>
        </Transition.Child>
        <div className="flex-shrink-0 w-14" aria-hidden="true"></div>
      </Dialog>
    </Transition.Root>
  );
};

const Sidebar = ({ navigation }: { navigation: DashboardNavigationItem[] }) => {
  return (
    <div className="hidden md:flex md:flex-shrink-0">
      <div className="flex flex-col w-64">
        <div className="flex flex-col h-0 flex-1">
          <div className="flex items-center h-16 flex-shrink-0 px-4 bg-background-light">
            <Logo />
          </div>
          <div className="flex-1 flex flex-col overflow-y-auto">
            <nav className="flex-1 px-2 py-4 bg-background-dark space-y-1">
              <SideNavigation navigation={navigation} />
            </nav>
          </div>
        </div>
      </div>
    </div>
  );
};

const Logo = () => {
  return (
    <Link className="flex items-center text-white-light" to=".">
      <img
        className="h-16 w-auto mr-2"
        src={APP_LOGO_IMAGE_URL}
        alt={APP_NAME}
      />
      <span className="text-xl text-accent-light font-semibold">
        {APP_NAME}
      </span>
    </Link>
  );
};

type DashLayoutProps = {
  navigation: DashboardNavigationItem[];
  userNavigation?: UserNavigationItem[];
  userElement?: React.ReactNode;
  children: React.ReactNode;
};

export const DashLayout = ({
  navigation,
  userNavigation,
  userElement,
  children,
}: DashLayoutProps) => {
  const [sidebarOpen, setSidebarOpen] = React.useState(false);

  return (
    <div className="h-screen flex overflow-hidden bg-background">
      <MobileSidebar
        navigation={navigation}
        sidebarOpen={sidebarOpen}
        setSidebarOpen={setSidebarOpen}
      />
      <Sidebar navigation={navigation} />
      <div className="flex flex-col w-0 flex-1 overflow-hidden">
        <div className="relative z-10 flex-shrink-0 flex h-16 bg-background-light shadow">
          <button
            className="px-4 text-accent-light focus:outline-none focus:ring-2 focus:ring-inset focus:ring-accent-light md:hidden"
            onClick={() => setSidebarOpen(true)}
          >
            <span className="sr-only">Open sidebar</span>
            <RxHamburgerMenu className="h-6 w-6" aria-hidden="true" />
          </button>

          {!userNavigation && !userElement ? null : (
            <div className="flex-1 flex justify-end items-center px-4 sm:px-6">
              {userNavigation && userElement ? (
                <UserNavigation
                  userNavigation={userNavigation}
                  userElement={userElement}
                />
              ) : null}
            </div>
          )}
        </div>
        <main className="flex-1 relative overflow-y-auto focus:outline-none">
          {children}
        </main>
      </div>
    </div>
  );
};
