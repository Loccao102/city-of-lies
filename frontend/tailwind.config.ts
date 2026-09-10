import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  darkMode: "class",
  theme: {
    extend: {
      colors: {
        background: "#0a0c10",
        surface: {
          DEFAULT: "#12161f",
          elevated: "#1a202c",
          border: "#2d3748",
        },
        danger: {
          DEFAULT: "#e53e3e",
          glow: "#fc8181",
        },
        warning: {
          DEFAULT: "#dd6b20",
          glow: "#f6ad55",
        },
        truth: {
          DEFAULT: "#319795",
          glow: "#4fd1c5",
        },
        credibility: {
          DEFAULT: "#3182ce",
          glow: "#63b3ed",
        },
      },
      fontFamily: {
        mono: ["Consolas", "Monaco", "Courier New", "monospace"],
      },
    },
  },
  plugins: [],
};

export default config;
