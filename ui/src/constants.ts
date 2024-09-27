export const APP_NAME = "Rory Pearson";
export const APP_LOGO_IMAGE_URL = "/images/logo.png";

export const navigation_items = [
  {
    name: "Home",
    to: "/",
  },
  {
    name: "Board",
    to: "/board",
  },
  {
    name: "Background Remover",
    to: "/background-remover",
  },
  {
    name: "Image To Icon",
    to: "/image-to-icon",
  },
  {
    name: "Spotify",
    to: "/spotify",
  },
  {
    name: "Components",
    to: "/components",
  },
].filter(Boolean) as {
  name: string;
  to: string;
  customUrl?: boolean;
}[];
