const defaultTheme = require("tailwindcss/defaultTheme");
const withMT = require("@material-tailwind/react/utils/withMT");

module.exports = withMT({
  mode: "jit",
  theme: {
    extend: {
      colors: {
        "accent-light": "#ee1a40",
        "accent-dark": "#d0062a",
        background: "#1e2120",
        "background-light": "#2a2e2c",
        "background-dark": "#131514",
        "t-light": "#ECECEC",
        "t-dark": "#CCCCCC",
      },
      fontFamily: {
        fira: ["Fira Code", ...defaultTheme.fontFamily.sans],
      },
      keyframes: {
        fadeIn: {
          "0%": { opacity: "0", transform: "translateY(20px)" },
          "100%": { opacity: "1", transform: "translateY(0)" },
        },
      },
      animation: {
        fadeIn: "fadeIn 0.6s ease-in-out forwards",
      },
    },
  },
  variants: {
    extend: {},
  },
  plugins: [require("@tailwindcss/typography"), require("flowbite/plugin")],
  content: [
    "./src/**/*.{html,js,ts,jsx,tsx}",
    "./src/index.html",
    "../../node_modules/flowbite-react/lib/esm/**/*.js",
    "../../node_modules/@material-tailwind/react/components/**/*.{js,ts,jsx,tsx}",
    "../../node_modules/@material-tailwind/react/theme/components/**/*.{js,ts,jsx,tsx}",
  ],
});
