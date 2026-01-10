/** @type {import('tailwindcss').Config} */
const config = {
  content: [
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        "dark-primary": "#131a1c",
        "dark-secondary": "#1b2224",
        red: "#e74c4c",
        green: "#6bb85d",
        blue: "#0183ff",
        grey: "#dddfe2",
        white: "#fff"
      },
    }
  },
  plugins: {
    "@tailwindcss/postcss": {},
  },
};

export default config;