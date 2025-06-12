/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/**/*.html"],
  theme: {
    extend: {
      colors: {
        "cyber-blue": "#00ffff",
        "cyber-pink": "#ff00ff",
        "cyber-yellow": "#ffff00",
        "cyber-green": "#00ff00",
        "cyber-purple": "#8000ff",
        "cyber-orange": "#ff8000",
      },
      fontFamily: {
        cyber: ["Orbitron", "Courier New", "monospace"],
      },
      animation: {
        glow: "glow 2s ease-in-out infinite alternate",
        scan: "scan 2s linear infinite",
      },
      keyframes: {
        glow: {
          "0%": { boxShadow: "0 0 5px currentColor" },
          "100%": { boxShadow: "0 0 20px currentColor, 0 0 30px currentColor" },
        },
        scan: {
          "0%": { transform: "translateX(-100%)" },
          "100%": { transform: "translateX(100%)" },
        },
      },
    },
  },
  plugins: [],
};
